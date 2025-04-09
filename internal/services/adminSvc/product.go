package adminSvc

import (
	"context"
	"errors"
	"fmt"

	"github.com/Daniel-Njaramba-1/pulse/internal/repo"
	"github.com/Daniel-Njaramba-1/pulse/internal/util/generics"
	"github.com/Daniel-Njaramba-1/pulse/internal/util/logging"
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

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(ctx context.Context, product *repo.Product) (*repo.Product, error) {
    // Validate category
    _, err := s.categoryService.GetCategoryByID(ctx, product.CategoryId)
    if err != nil {
        logging.LogError("Category validation failed: category not found")
        return nil, errors.New("category not found")
    }

    // Validate brand
    _, err = s.brandService.GetBrandByID(ctx, product.BrandId)
    if err != nil {
        logging.LogError("Brand validation failed: brand not found")
        return nil, errors.New("brand not found")
    }

    // Validate required fields
    if product.Name == "" || product.Description == "" || product.ImagePath == "" {
        logging.LogError("Product creation failed: missing required fields")
        return nil, errors.New("missing required fields")
    }

    // Create product
    _, err = generics.CreateModel(ctx, s.db, product)
    if err != nil {
        logging.LogError(fmt.Sprintf("Failed to create product in database: %v", err))
        return nil, err
    }

    logging.LogInfo("Product successfully created")
    return product, nil
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

// UpdateProduct updates an existing product
func (s *ProductService) UpdateProduct(ctx context.Context, product *repo.Product) (*repo.Product, error) {
    // Validate product ID
    if product.Id == 0 {
        return nil, errors.New("product ID is required")
    }

    // Validate category
    _, err := s.categoryService.GetCategoryByID(ctx, product.CategoryId)
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

    // Update product
    err = generics.UpdateModelDetails(ctx, s.db, product)
    if err != nil {
        return nil, err
    }

    // Retrieve updated product
    updatedProduct, err := s.GetProductByID(ctx, product.Id)
    if err != nil {
        return nil, err
    }
    return updatedProduct, nil
}

// DeactivateProduct deactivates a product
func (s *ProductService) DeactivateProduct(ctx context.Context, id int) error {
    var product repo.Product
    product.Id = id
    err := generics.DeactivateModel(ctx, s.db, &product)
    if err != nil {
        return err
    }
    return nil
}

// ReactivateProduct reactivates a product
func (s *ProductService) ReactivateProduct(ctx context.Context, id int) error {
    var product repo.Product
    product.Id = id
    err := generics.ReactivateModel(ctx, s.db, &product)
    if err != nil {
        return err
    }
    return nil
}

// DeleteProduct deletes a product
func (s *ProductService) DeleteProduct(ctx context.Context, id int) error {
    var product repo.Product
    product.Id = id
    err := generics.DeleteModel(ctx, s.db, &product)
    if err != nil {
        return err
    }
    return nil
}