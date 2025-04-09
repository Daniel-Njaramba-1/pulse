package repo

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type Cart struct {
	Id				int			`db:"id" json:"id"`
	CustomerId		int			`db:"customer_id" json:"customer_id"`
	IsActive		bool		`db:"is_active" json:"is_active"`
	CreatedAt		time.Time	`db:"created_at" json:"created_at"`
	UpdatedAt		time.Time	`db:"updated_at" json:"updated_at"`
}

func (cart *Cart) FeedGetId() *int {
	return &cart.Id
}

func (cart *Cart) FeedCreateQuery() string {
	return `
		INSERT INTO carts (customer_id, is_active)
		VALUES (:customer_id, :is_active)
		RETURNING id
	`
}

func (cart *Cart) FeedGetByIdQuery() string {
	return `
		SELECT *
		FROM carts
		WHERE id = $1
	`
}

func (cart *Cart) FeedGetAllQuery() string {
	return `
		SELECT *
		FROM carts
		ORDER BY id ASC
	`
}

func (cart *Cart) FeedUpdateDetailsQuery() string {
	return `
		UPDATE carts
		SET customer_id = :customer_id
		WHERE id = :id
	`
}

func (cart *Cart) FeedDeactivateQuery() string {
	return `
		UPDATE carts
		SET is_active = FALSE
		WHERE id = :id
	`
}

func (cart *Cart) FeedReactivateQuery() string {
	return `
		UPDATE carts
		SET is_active = TRUE
		WHERE id = :id
	`
}

func (cart *Cart) FeedDeleteQuery() string {
	return `
		DELETE FROM carts
		WHERE id = :id
	`
}

func GetCartByCustomer(ctx context.Context, db *sqlx.DB, customer_id int, cart *Cart) error {
	query := `
		SELECT *
		FROM carts
		WHERE customer_id = $1
	`
	return db.Get(cart, query, customer_id)
}