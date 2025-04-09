package repo

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type Category struct {
	Id			int			`db:"id" json:"id"`
	Name		string		`db:"name" json:"name"`
	Description	string		`db:"description" json:"description"`
	IsActive	bool		`db:"is_active" json:"is_active"`
	CreatedAt	time.Time	`db:"created_at" json:"created_at"`
	UpdatedAt	time.Time	`db:"updated_at" json:"updated_at"`
}

func (category *Category) FeedGetId() *int {
    return &category.Id
}

func (category *Category) FeedCreateQuery() string {
    return `
        INSERT INTO categories (name, description, is_active)
        VALUES (:name, :description, :is_active)
        RETURNING id
    `
}

func (category *Category) FeedGetByIdQuery() string {
    return `
        SELECT *
        FROM categories
        WHERE id = $1
    `
}

func (category *Category) FeedGetAllQuery() string {
    return `
        SELECT *
        FROM categories
        ORDER BY id ASC
    `
}

func (category *Category) FeedUpdateDetailsQuery() string {
    return `
        UPDATE categories
        SET name = :name,
            description = :description
        WHERE id = :id
    `
}

func (category *Category) FeedDeactivateQuery() string {
    return `
        UPDATE categories
        SET is_active = FALSE
        WHERE id = :id
    `
}

func (category *Category) FeedReactivateQuery() string {
    return `
        UPDATE categories
        SET is_active = TRUE
        WHERE id = :id
    `
}

func (category *Category) FeedDeleteQuery() string {
    return `
        DELETE FROM categories
        WHERE id = :id
    `
}

func SearchCategoriesByName(ctx context.Context, db *sqlx.DB, name string, categories *[]Category) error {
    query := `	
        SELECT *
        FROM categories 
        WHERE name ILIKE $1
    `
    return db.Select(categories, query, "%"+name+"%")
}

func GetAllProductsByCategory(ctx context.Context, db *sqlx.DB, category_id int, products *[]Product) error {
    query := `	
        SELECT *
        FROM products 
        WHERE category_id = $1
    `
    return db.Select(products, query, category_id)
}