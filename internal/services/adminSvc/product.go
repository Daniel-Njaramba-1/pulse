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

// CreateProduct creates a new product along with its metrics and initial stock
func (s *ProductService) CreateProduct(ctx context.Context, product *repo.Product, initialStock int) (*repo.Product, error) {
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
    if product.Name == "" || product.Description == "" || product.ImagePath == "" {
        return nil, errors.New("missing required fields")
    }

    // Insert product
    product.IsActive = true
    insertProductQuery := `
        INSERT INTO products (name, description, image_path, is_active, category_id, brand_id)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `
    err = tx.QueryRowContext(ctx, insertProductQuery, product.Name, product.Description, product.ImagePath, product.IsActive, product.CategoryId, product.BrandId).Scan(&product.Id)
    if err != nil {
        return nil, fmt.Errorf("failed to insert product: %w", err)
    }

    // Insert product metrics
    insertMetricsQuery := `
        INSERT INTO product_metrics (product_id, average_rating, review_count, wishlist_count, base_price, adjusted_price)
        VALUES ($1, 0, 0, 0, 0, 0)
    `
    _, err = tx.ExecContext(ctx, insertMetricsQuery, product.Id)
    if err != nil {
        return nil, fmt.Errorf("failed to insert product metrics: %w", err)
    }

    // Insert initial stock
    insertStockQuery := `
        INSERT INTO stocks (product_id, quantity, stock_threshold)
        VALUES ($1, $2, 0)
    `
    _, err = tx.ExecContext(ctx, insertStockQuery, product.Id, initialStock)
    if err != nil {
        return nil, fmt.Errorf("failed to insert initial stock: %w", err)
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
func (s *ProductService) GetAllProducts(ctx context.Context) ([]*repo.Product, error) {
    var products []*repo.Product
    query := `
        SELECT id, name, description, image_path, is_active, category_id, brand_id
        FROM products
    `
    err := s.db.SelectContext(ctx, &products, query)
    if err != nil {
        return nil, err
    }
    return products, nil
}

// UpdateProduct updates an existing product
func (s *ProductService) UpdateProduct(ctx context.Context, product *repo.Product) (*repo.Product, error) {
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

    // Update product
    updateProductQuery := `
        UPDATE products
        SET name = $1, description = $2, image_path = $3, category_id = $4, brand_id = $5
        WHERE id = $6
    `
    _, err = tx.ExecContext(ctx, updateProductQuery, product.Name, product.Description, product.ImagePath, product.CategoryId, product.BrandId, product.Id)
    if err != nil {
        return nil, fmt.Errorf("failed to update product: %w", err)
    }

    // Commit transaction
    if err = tx.Commit(); err != nil {
        return nil, fmt.Errorf("failed to commit transaction: %w", err)
    }

    return product, nil
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

func (s *ProductService) SetBasePrice(ctx context.Context, id int, price float64) error {
    return nil
}