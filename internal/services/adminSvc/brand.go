package adminSvc

import (
    "context"
    "database/sql"
    "errors"

    "github.com/Daniel-Njaramba-1/pulse/internal/repo"
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
	brand.IsActive = true

    query := `
        INSERT INTO brands (name, description, is_active)
        VALUES ($1, $2, $3)
        RETURNING id
    `
    err := s.db.QueryRowxContext(ctx, query, brand.Name, brand.Description, brand.IsActive).Scan(&brand.Id)
    if err != nil {
        return nil, err
    }

    return brand, nil
}

// GetBrandByID retrieves a brand by its ID.
func (s *BrandService) GetBrandByID(ctx context.Context, id int) (*repo.Brand, error) {
    var brand repo.Brand
    query := `
        SELECT id, name, description
        FROM brands
        WHERE id = $1
    `
    err := s.db.GetContext(ctx, &brand, query, id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.New("brand not found")
        }
        return nil, err
    }
    return &brand, nil
}

// GetAllBrands retrieves all brands.
func (s *BrandService) GetAllBrands(ctx context.Context) ([]*repo.Brand, error) {
    var brands []*repo.Brand
    query := `
        SELECT id, name, description
        FROM brands
    `
    err := s.db.SelectContext(ctx, &brands, query)
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

    query := `
        UPDATE brands
        SET name = $1, description = $2
        WHERE id = $3
    `
    _, err := s.db.ExecContext(ctx, query, brand.Name, brand.Description, brand.Id)
    if err != nil {
        return nil, err
    }

    return s.GetBrandByID(ctx, brand.Id)
}

// DeactivateBrand deactivates a brand.
func (s *BrandService) DeactivateBrand(ctx context.Context, id int) error {
    query := `
        UPDATE brands
        SET is_active = false
        WHERE id = $1
    `
    _, err := s.db.ExecContext(ctx, query, id)
    if err != nil {
        return err
    }
    return nil
}

// ReactivateBrand reactivates a brand.
func (s *BrandService) ReactivateBrand(ctx context.Context, id int) error {
    query := `
        UPDATE brands
        SET is_active = true
        WHERE id = $1
    `
    _, err := s.db.ExecContext(ctx, query, id)
    if err != nil {
        return err
    }
    return nil
}

// DeleteBrand deletes a brand by its ID.
func (s *BrandService) DeleteBrand(ctx context.Context, id int) error {
    query := `
        DELETE FROM brands
        WHERE id = $1
    `
    _, err := s.db.ExecContext(ctx, query, id)
    if err != nil {
        return err
    }
    return nil
}