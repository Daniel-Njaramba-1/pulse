package adminSvc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Daniel-Njaramba-1/pulse/internal/repo"
	"github.com/jmoiron/sqlx"
)

type ProductService struct {
    db              *sqlx.DB
    categoryService *CategoryService
    brandService    *BrandService
}

func NewProductService(db *sqlx.DB, categoryService *CategoryService, brandService *BrandService) *ProductService {
    return &ProductService{
        db:              db,
        categoryService: categoryService,
        brandService:    brandService,
    }
}

// CreateProduct creates a new product along with its metrics and initial stock using a single query with CTEs
func (s *ProductService) CreateProduct(ctx context.Context, product *repo.Product, basePrice int, initialStock int) (*repo.Product, error) {
    tx, err := s.db.BeginTxx(ctx, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback()

    // Validate category
    _, err = s.categoryService.GetCategoryByID(ctx, product.CategoryId)
    if err != nil {
        return nil, errors.New("category not found")
    }

    // Validate brand
    _, err = s.brandService.GetBrandByID(ctx, product.BrandId)
    if err != nil {
        return nil, errors.New("brand not found")
    }

    // Validate required fields
    if product.Name == "" || product.Description == "" || product.ImagePath == nil {
        return nil, errors.New("missing required fields")
    }

    // Use a single query with CTEs to insert product, metrics, and stock
    query := `
        WITH new_product AS (
            INSERT INTO products (name, description, image_path, is_active, category_id, brand_id)
            VALUES ($1, $2, $3, $4, $5, $6)
            RETURNING id
        ),
        new_metrics AS (
            INSERT INTO product_metrics (product_id, average_rating, review_count, wishlist_count, base_price, adjusted_price)
            SELECT id, 0, 0, 0, $7, $7 FROM new_product
        ),
        new_stock AS (
            INSERT INTO stocks (product_id, quantity, stock_threshold)
            SELECT id, $8, 0 FROM new_product
        )
        SELECT id FROM new_product
    `
    err = tx.QueryRowContext(ctx, query, product.Name, product.Description, product.ImagePath, true, product.CategoryId, product.BrandId, basePrice, initialStock).Scan(&product.Id)
    if err != nil {
        return nil, fmt.Errorf("failed to execute CTE query: %w", err)
    }

    // Commit transaction
    if err = tx.Commit(); err != nil {
        return nil, fmt.Errorf("failed to commit transaction: %w", err)
    }

    return product, nil
}

// GetProductByID retrieves a product by its ID along with brand and category names and metrics/stock
func (s *ProductService) GetProductByID(ctx context.Context, id int) (*repo.ProductDetail, error) {
    var productDetail repo.ProductDetail
    query := `
        SELECT
            p.id, p.name, p.description, p.image_path, p.is_active,
            b.id AS brand_id,  b.name AS brand_name,
            c.id AS category_id, c.name AS category_name,
            pm.average_rating, pm.review_count, pm.wishlist_count, pm.base_price, pm.adjusted_price,
            st.quantity AS stock_quantity, st.stock_threshold
        FROM products p
        INNER JOIN brands b ON p.brand_id = b.id
        INNER JOIN categories c ON p.category_id = c.id
        LEFT JOIN product_metrics pm ON p.id = pm.product_id
        LEFT JOIN stocks st ON p.id = st.product_id
        WHERE p.id = $1
    `
    err := s.db.GetContext(ctx, &productDetail, query, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, errors.New("product not found")
        }
        return nil, err
    }
    return &productDetail, nil
}

// GetAllProducts retrieves all products
func (s *ProductService) GetAllProducts(ctx context.Context) ([]*repo.ProductDetail, error) {
    var products []*repo.ProductDetail
    query := `
        SELECT
            p.id, p.name, p.description, p.image_path, p.is_active,
            b.id AS brand_id,  b.name AS brand_name,
            c.id AS category_id, c.name AS category_name,
            pm.average_rating, pm.review_count, pm.wishlist_count, pm.base_price, pm.adjusted_price,
            st.quantity AS stock_quantity, st.stock_threshold
        FROM products p
        INNER JOIN brands b ON p.brand_id = b.id
        INNER JOIN categories c ON p.category_id = c.id
        LEFT JOIN product_metrics pm ON p.id = pm.product_id
        LEFT JOIN stocks st ON p.id = st.product_id
        ORDER BY p.id
    `
    err := s.db.SelectContext(ctx, &products, query)
    if err != nil {
        return nil, err
    }
    return products, nil
}

// UpdateProductDetails updates the details of an existing product excluding the image
func (s *ProductService) UpdateProductDetails(ctx context.Context, product *repo.Product) (*repo.Product, error) {
    tx, err := s.db.BeginTxx(ctx, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback()

    // Validate category
    _, err = s.categoryService.GetCategoryByID(ctx, product.CategoryId)
    if err != nil {
        return nil, errors.New("category not found")
    }

    // Validate brand
    _, err = s.brandService.GetBrandByID(ctx, product.BrandId)
    if err != nil {
        return nil, errors.New("brand not found")
    }

    // Update product details
    updateDetailsQuery := `
        UPDATE products
        SET name = $1, description = $2, category_id = $3, brand_id = $4
        WHERE id = $5
    `
    _, err = tx.ExecContext(ctx, updateDetailsQuery, product.Name, product.Description, product.CategoryId, product.BrandId, product.Id)
    if err != nil {
        return nil, fmt.Errorf("failed to update product details: %w", err)
    }

    // Commit transaction
    if err = tx.Commit(); err != nil {
        return nil, fmt.Errorf("failed to commit transaction: %w", err)
    }

    return product, nil
}

func (s *ProductService) GetProductImagePath(ctx context.Context, id int) (string, error) {
    var product struct {
        ImagePath *string `db:"image_path"`
    }
    query := `
        SELECT image_path
        FROM products
        WHERE id = $1
    `
    err := s.db.GetContext(ctx, &product, query, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return "", errors.New("product not found")
        }
        return "", err
    }
    
    // Handle null image paths
    if product.ImagePath == nil {
        return "", nil
    }
    
    return *product.ImagePath, nil
}

// UpdateProductImage updates the image of an existing product
func (s *ProductService) UpdateProductImage(ctx context.Context, productId int, imagePath string) error {
    tx, err := s.db.BeginTxx(ctx, nil)
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback()

    // Update product image
    updateImageQuery := `
        UPDATE products
        SET image_path = $1
        WHERE id = $2
    `
    _, err = tx.ExecContext(ctx, updateImageQuery, imagePath, productId)
    if err != nil {
        return fmt.Errorf("failed to update product image: %w", err)
    }

    // Commit transaction
    if err = tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }

    return nil
}

// DeactivateProduct deactivates a product
func (s *ProductService) DeactivateProduct(ctx context.Context, id int) error {
    query := `
        UPDATE products
        SET is_active = false
        WHERE id = $1
    `
    _, err := s.db.ExecContext(ctx, query, id)
    return err
}

// ReactivateProduct reactivates a product
func (s *ProductService) ReactivateProduct(ctx context.Context, id int) error {
    query := `
        UPDATE products
        SET is_active = true
        WHERE id = $1
    `
    _, err := s.db.ExecContext(ctx, query, id)
    return err
}

// DeleteProduct deletes a product
func (s *ProductService) DeleteProduct(ctx context.Context, id int) error {
    query := `
        DELETE FROM products
        WHERE id = $1
    `
    _, err := s.db.ExecContext(ctx, query, id)
    return err
}

func (s *ProductService) ChangeBasePrice(ctx context.Context, id int, price float64) error {
    tx, err := s.db.BeginTxx(ctx, nil)
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback()

    // Get Product
    _, err = s.GetProductByID(ctx, id)
    if err != nil {
        return fmt.Errorf("failed to get product: %w", err)
    }

    // Update base price
    updatePriceQuery := `
        UPDATE product_metrics
        SET base_price = $1
        WHERE product_id = $2
    `
    _, err = tx.ExecContext(ctx, updatePriceQuery, price, id)
    if err != nil {
        return fmt.Errorf("failed to update base price: %w", err)
    }

    // Commit transaction
    if err = tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }

    return nil
}

func (s *ProductService) RestockProduct(ctx context.Context, stock *repo.Stock) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Validate product
	var productId int
	getProductQuery := `
		SELECT id
		FROM products
		WHERE id = $1
	`
	err = tx.QueryRowContext(ctx, getProductQuery, stock.ProductId).Scan(&productId)
	if err != nil {
		if err == sql.ErrNoRows {
			return err // Or a custom error indicating product not found
		}
		return err
	}

	// Check if stock for the product already exists
	var existingStock repo.Stock
	getStockQuery := `
		SELECT id, product_id, quantity, stock_threshold, created_at, updated_at
		FROM stocks
		WHERE product_id = $1
	`
	err = tx.GetContext(ctx, &existingStock, getStockQuery, stock.ProductId)

	if err == nil {
		// Update existing stock
		updateStockQuery := `
			UPDATE stocks
			SET quantity = quantity + $1 
			WHERE product_id = $2
		`
		_, err = tx.ExecContext(ctx, updateStockQuery, stock.Quantity, stock.ProductId)
		if err != nil {
			return err
		}
	} else if err == sql.ErrNoRows {
		// Create new stock level
		insertStockQuery := `
			INSERT INTO stocks(product_id, quantity, stock_threshold)
			VALUES ($1, $2, $3)
			RETURNING id
		`
		err = tx.QueryRowContext(ctx, insertStockQuery, stock.ProductId, stock.Quantity, stock.StockThreshold).Scan(&stock.Id)
		if err != nil {
			return err
		}
	} else {
		return err
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *ProductService) GetProductName (ctx context.Context, id int) (string, error) {
    var product struct {
        Name string `db:"name"`
    }
    query := `
        SELECT name
        FROM products
        WHERE id = $1
    `
    err := s.db.GetContext(ctx, &product, query, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return "", errors.New("product not found")
        }
        return "", err
    }
    
    return product.Name, nil
}