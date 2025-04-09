package repo

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type CartItem struct {
	Id				int			`db:"id" json:"id"`
	CartId			int			`db:"cart_id" json:"cart_id"`
	ProductId		int			`db:"product_id" json:"product_id"`
	Quantity		int			`db:"quantity" json:"quantity"`
	IsProcessed		bool		`db:"is_processed" json:"is_processed"`
	CreatedAt		time.Time	`db:"created_at" json:"created_at"`
	UpdatedAt		time.Time	`db:"updated_at" json:"updated_at"`
}

func (cart_item *CartItem) FeedGetId() *int {
	return &cart_item.Id
}

func (cart_item *CartItem) FeedCreateQuery() string {
	return `
		INSERT INTO cart_items (cart_id, product_id, quantity, is_processed)
		VALUES (:cart_id, :product_id, :quantity, :is_processed)
		RETURNING id
	`
}

func (cart_item *CartItem) FeedGetAllQuery() string {
	return `
		SELECT *
		FROM cart_items
		ORDER BY id ASC
	`
}

func (cart_item *CartItem) FeedUpdateDetailsQuery() string {
	return `
		UPDATE cart_items
		SET cart_id = :cart_id
			product_id = :product_id
			quantity = :quantity
		WHERE id = :id
	`
}

func (cart_item *CartItem) FeedDeleteQuery() string {
	return `
		DELETE FROM cart_items
		WHERE id = :id
	`
}

func GetCartItemsByCart(ctx context.Context, db *sqlx.DB, cart_id int, cart_items *[]CartItem) error {
	query := `
		SELECT id, cart_id, product_id, quantity, is_active, created_at, updated_at
		FROM cart_items
		WHERE cart_id = $1
		ORDER BY id ASC
	`
	return db.Select(cart_items, query, cart_id)
}