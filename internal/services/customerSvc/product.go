package customerSvc

import (
	"context"

	"github.com/Daniel-Njaramba-1/pulse/internal/repo"
	"github.com/jmoiron/sqlx"
)

type ProductService struct {
	db *sqlx.DB
}

func NewProductService(db *sqlx.DB) *ProductService {
	return &ProductService{db: db}
}

// GetProductByID retrieves a product by its ID
func (s *ProductService) GetProductByID(ctx context.Context, productId int) (*repo.ProductDetail, error) {
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
    err := s.db.GetContext(ctx, &productDetail, query, productId)
    if err != nil {
        return nil, err
    }

    return &productDetail, nil
}

// GetProductByID retrieves a product by its ID
func (s *ProductService) GetProductByName(ctx context.Context, productName string) (*repo.ProductDetail, error) {
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
        WHERE p.name = $1
    `
    err := s.db.GetContext(ctx, &productDetail, query, productName)
    if err != nil {
        return nil, err
    }

    return &productDetail, nil
}

// GetAllProducts retrieves all products
func (s *ProductService) GetAllProducts(ctx context.Context) ([]*repo.ProductDetail, error) {
	var productDetail []*repo.ProductDetail
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
    err := s.db.SelectContext(ctx, &productDetail, query)
    if err != nil {
        return nil, err
    }
    return productDetail, nil
}