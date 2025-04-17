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
func (s *ProductService) GetProductByID(ctx context.Context, id int, name string) (*repo.Product, error) {
    var product repo.Product

	

    return &product, nil
}

// GetAllProducts retrieves all products
func (s *ProductService) GetAllProducts(ctx context.Context) ([]*repo.Product, error) {
    var products []*repo.Product

    return products, nil
}