package repo

import (
	"time"
)

type Product struct {
	Id				int			`db:"id" json:"id"`
	CategoryId	    int			`db:"category_id" json:"category_id"`
	BrandId			int			`db:"brand_id" json:"brand_id"`
	Name			string		`db:"name" json:"name"`
	Description		string		`db:"description" json:"description"`
	ImagePath		string		`db:"image_path" json:"image_path"`
	IsActive		bool		`db:"is_active" json:"is_active"` 
	CreatedAt		time.Time	`db:"created_at" json:"created_at"`
	UpdatedAt		time.Time	`db:"updated_at" json:"updated_at"`
}

func (product *Product) FeedGetId() *int {
	return &product.Id
}

func (product *Product) FeedCreateQuery() string {
	return `
		INSERT INTO products (category_id, brand_id, name, description, image_path, is_active)
		VALUES (:category_id, :brand_id, :name, :description, :image_path, :is_active)
		RETURNING id
	`
}

func (product *Product) FeedGetByIdQuery() string {
	return `
		SELECT * 
		FROM products
		WHERE id = :id
	`
}

func (product *Product) FeedGetAllQuery() string {
	return `
		SELECT * 
		FROM products
		ORDER BY id ASC
	`
}

func (product *Product) FeedUpdateDetailsQuery() string {
	return `
		UPDATE products
		SET category_id = :category_id, 
			brand_id = :brand_id, 
			name = :name, 
			description = :description, 
			image_path = :image_path
		WHERE id = :id
	`
}

func (product *Product) FeedDeleteQuery() string {
	return `
		DELETE FROM products
		WHERE id = :id
	`
}

func (product *Product) FeedDeactivateQuery() string {
	return `
		UPDATE products
		SET is_active = false
		WHERE id = :id
	`
}

func (product *Product) FeedReactivateQuery() string {
	return `
		UPDATE products
		SET is_active = true
		WHERE id = :id
	`
}
