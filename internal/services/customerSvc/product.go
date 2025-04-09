package customerSvc

import (
	"context"
	"errors"

	"github.com/Daniel-Njaramba-1/pulse/internal/repo"
	"github.com/Daniel-Njaramba-1/pulse/internal/util/generics"
	"github.com/jmoiron/sqlx"
)

type ProductService struct {
	db *sqlx.DB
}

func NewProductService(db *sqlx.DB) *ProductService {
	return &ProductService{db: db}
}

// GetProductByID retrieves a product by its ID
func (s *ProductService) GetProductByID(ctx context.Context, id int) (*repo.Product, error) {
    var product repo.Product
    err := generics.SelectModelById(ctx, s.db, id, &product)
    if err != nil {
        return nil, errors.New("product not found")
    }
    return &product, nil
}

// GetAllProducts retrieves all products
func (s *ProductService) GetAllProducts(ctx context.Context) ([]*repo.Product, error) {
    var products []*repo.Product
    err := generics.SelectAllModels(ctx, s.db, &products)
    if err != nil {
        return nil, err
    }
    return products, nil
}