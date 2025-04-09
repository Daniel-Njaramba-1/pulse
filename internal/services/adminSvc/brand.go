// filepath: c:\Users\ADMIN\Desktop\pulse\internal\services\adminSvc\brands.go
package adminSvc

import (
	"context"
	"errors"

	"github.com/Daniel-Njaramba-1/pulse/internal/repo"
	"github.com/Daniel-Njaramba-1/pulse/internal/util/generics"
	"github.com/jmoiron/sqlx"
)

type BrandService struct {
	db *sqlx.DB
}

func NewBrandService(db *sqlx.DB) *BrandService {
	return &BrandService{db: db}
}

// CreateBrand creates a new brand.
func (s *BrandService) CreateBrand(ctx context.Context, brand *repo.Brand) (*repo.Brand, error) {
	if brand.Name == "" || brand.Description == "" {
		return nil, errors.New("missing required fields")
	}

	_, err := generics.CreateModel(ctx, s.db, brand)
	if err != nil {
		return nil, err
	}

	return brand, nil
}

// GetBrandByID retrieves a brand by its ID.
func (s *BrandService) GetBrandByID(ctx context.Context, id int) (*repo.Brand, error) {
	var brand repo.Brand
	err := generics.SelectModelById(ctx, s.db, id, &brand)
	if err != nil {
		return nil, err
	}
	return &brand, nil
}

// GetAllBrands retrieves all brands.
func (s *BrandService) GetAllBrands(ctx context.Context) ([]*repo.Brand, error) {
	var brands []*repo.Brand
	err := generics.SelectAllModels(ctx, s.db, &brands)
	if err != nil {
		return nil, err
	}
	return brands, nil
}

// UpdateBrand updates the details of an existing brand.
func (s *BrandService) UpdateBrand(ctx context.Context, brand *repo.Brand) (*repo.Brand, error) {
	if brand.Id == 0 || brand.Name == "" || brand.Description == "" {
		return nil, errors.New("missing required fields")
	}

	err := generics.UpdateModelDetails(ctx, s.db, brand)
	if err != nil {
		return nil, err
	}
	brand, err = s.GetBrandByID(ctx, brand.Id)
	if err != nil {
		return nil, err
	}
	return brand, nil
}

// DeactivateBrand deactivates a brand.
func (s *BrandService) DeactivateBrand(ctx context.Context, id int) error {
    var brand repo.Brand
    brand.Id = id
    err := generics.DeactivateModel(ctx, s.db, &brand)
    if err != nil {
        return err
    }
    return nil
}

// ReactivateBrand reactivates a brand.
func (s *BrandService) ReactivateBrand(ctx context.Context, id int) error {
    var brand repo.Brand
    brand.Id = id
    err := generics.ReactivateModel(ctx, s.db, &brand)
    if err != nil {
        return err
    }
    return nil
}

// DeleteBrand deletes a brand by its ID.
func (s *BrandService) DeleteBrand(ctx context.Context, id int) error {
    var brand repo.Brand
    brand.Id = id
    err := generics.DeleteModel(ctx, s.db, &brand)
    if err != nil {
        return err
    }
    return nil
}