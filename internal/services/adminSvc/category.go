package adminSvc

import (
	"context"
	"errors"

	"github.com/Daniel-Njaramba-1/pulse/internal/repo"
	"github.com/Daniel-Njaramba-1/pulse/internal/util/generics"
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

	_, err := generics.CreateModel(ctx, s.db, category)
	if err != nil {
		return  nil, err
	}

	return category, nil
}

func (s *CategoryService) GetCategoryByID(ctx context.Context, id int) (*repo.Category, error) {
	var category repo.Category
	err := generics.SelectModelById(ctx, s.db, id, &category)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (s *CategoryService) GetAllCategories(ctx context.Context) ([]*repo.Category, error) {
	var categories []*repo.Category
	err := generics.SelectAllModels(ctx, s.db, &categories)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *CategoryService) UpdateCategory(ctx context.Context, category *repo.Category) error {
	if category.Id == 0 || category.Name == "" || category.Description == "" {
		return errors.New("missing required fields")
	}

	err := generics.UpdateModelDetails(ctx, s.db, category)
	if err != nil {
		return err
	}
	return nil
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id int) error {
	var category repo.Category
	category.Id = id
	err := generics.DeleteModel(ctx, s.db, &category)
	if err != nil {
		return err
	}
	return nil
}

func (s *CategoryService) DeactivateCategory(ctx context.Context, id int) error {
	var category repo.Category
	category.Id = id
	err := generics.DeactivateModel(ctx, s.db, &category)
	if err != nil {
		return err
	}
	return nil
}

func (s *CategoryService) ReactivateCategory(ctx context.Context, id int) error {
	var category repo.Category
	category.Id = id
	err := generics.ReactivateModel(ctx, s.db, &category)
	if err != nil {
		return err
	}
	return nil
}