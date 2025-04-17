package adminSvc

import (
    "context"
    "errors"

    "github.com/Daniel-Njaramba-1/pulse/internal/repo"
    "github.com/jmoiron/sqlx"
)

type CategoryService struct {
    db *sqlx.DB
}

func NewCategoryService(db *sqlx.DB) *CategoryService {
    return &CategoryService{db: db}
}

func (s *CategoryService) CreateCategory(ctx context.Context, category *repo.Category) (*repo.Category, error) {
    if category.Name == "" || category.Description == "" {
        return nil, errors.New("missing required fields")
    }
    category.IsActive = true

    query := `
        INSERT INTO categories (name, description, is_active)
        VALUES ($1, $2, $3)
        RETURNING id
    `
    err := s.db.QueryRowxContext(ctx, query, category.Name, category.Description, category.IsActive).Scan(&category.Id)
    if err != nil {
        return nil, err
    }

    return category, nil
}

func (s *CategoryService) GetCategoryByID(ctx context.Context, id int) (*repo.Category, error) {
    var category repo.Category
    query := `
        SELECT id, name, description
        FROM categories
        WHERE id = $1
    `
    err := s.db.GetContext(ctx, &category, query, id)
    if err != nil {
        return nil, err
    }
    return &category, nil
}

func (s *CategoryService) GetAllCategories(ctx context.Context) ([]*repo.Category, error) {
    var categories []*repo.Category
    query := `
        SELECT id, name, description
        FROM categories
    `
    err := s.db.SelectContext(ctx, &categories, query)
    if err != nil {
        return nil, err
    }
    return categories, nil
}

func (s *CategoryService) UpdateCategory(ctx context.Context, category *repo.Category) error {
    if category.Id == 0 || category.Name == "" || category.Description == "" {
        return errors.New("missing required fields")
    }

    query := `
        UPDATE categories
        SET name = $1, description = $2
        WHERE id = $3
    `
    _, err := s.db.ExecContext(ctx, query, category.Name, category.Description, category.Id)
    if err != nil {
        return err
    }
    return nil
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id int) error {
    query := `
        DELETE FROM categories
        WHERE id = $1
    `
    _, err := s.db.ExecContext(ctx, query, id)
    if err != nil {
        return err
    }
    return nil
}

func (s *CategoryService) DeactivateCategory(ctx context.Context, id int) error {
    query := `
        UPDATE categories
        SET is_active = false
        WHERE id = $1
    `
    _, err := s.db.ExecContext(ctx, query, id)
    if err != nil {
        return err
    }
    return nil
}

func (s *CategoryService) ReactivateCategory(ctx context.Context, id int) error {
    query := `
        UPDATE categories
        SET is_active = true
        WHERE id = $1
    `
    _, err := s.db.ExecContext(ctx, query, id)
    if err != nil {
        return err
    }
    return nil
}