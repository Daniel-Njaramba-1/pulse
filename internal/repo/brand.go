package repo

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type Brand struct {
	Id				int			`db:"id" json:"id"`
	Name			string		`db:"name" json:"name"`
	Description		string		`db:"description" json:"description"`
	IsActive		bool		`db:"is_active" json:"is_active"`
	CreatedAt		time.Time	`db:"created_at" json:"created_at"`
	UpdatedAt		time.Time	`db:"updated_at" json:"updated_at"`
}

func (brand *Brand) FeedGetId() *int {
	return &brand.Id
}   

func (brand *Brand) FeedCreateQuery() string {
	return `
		INSERT INTO brands (name, description, is_active)
		VALUES (:name, :description, :is_active)
		RETURNING id
	`   
}

func (brand *Brand) FeedGetByIdQuery() string {
	return `
		SELECT *
		FROM brands
		WHERE id = $1
	`
}

func (brand *Brand) FeedGetAllQuery() string {
	return `
		SELECT *
		FROM brands
		ORDER BY id ASC
	`
}

func (brand *Brand) FeedUpdateDetailsQuery() string {
	return `    
		UPDATE brands
		SET name = :name, 
            description = :description
		WHERE id = :id
	`
}

func (brand *Brand) FeedDeactivateQuery() string {
	return `
		UPDATE brands
		SET is_active = false
		WHERE id = :id
	`
}

func (brand *Brand) FeedReactivateQuery() string {
	return `
		UPDATE brands
		SET is_active = true
		WHERE id = :id
	`
}

func (brand *Brand) FeedDeleteQuery() string {
	return `
		DELETE FROM brands
		WHERE id = :id
	`
}

// SearchBrandsByName searches brands with a similar name.
func SearchBrandsByName(ctx context.Context, db *sqlx.DB, name string, brands *[]Brand) error {
    query := `
        SELECT *
        FROM brands
        WHERE name ILIKE $1
    `
    return db.Select(brands, query, "%"+name+"%")
}

// GetAllProductsByBrand retrieves all products for a brand.
func GetAllProductsByBrand(ctx context.Context, db *sqlx.DB, brandID int, products *[]Product) error {
    query := `	
        SELECT *
        FROM products 
        WHERE brand_id = $1
    `
    return db.Select(products, query, brandID)
}