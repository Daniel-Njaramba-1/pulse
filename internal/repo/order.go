package repo

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type Order struct {
	Id				int			`db:"id" json:"id"`
	CustomerId		int			`db:"customer_id" json:"customer_id"`
	TotalPrice		float64		`db:"total_price" json:"total_price"`
	Status			string		`db:"status" json:"status"`
	PriceValidUntil time.Time	`db:"price_valid_until" json:"price_valid_until"`
	CreatedAt		time.Time	`db:"created_at" json:"created_at"`
	UpdatedAt		time.Time	`db:"updated_at" json:"updated_at"`
}

func (order *Order) FeedGetId() *int {
	return &order.Id
}

func (order *Order) FeedCreateQuery() string {
	return `
		INSERT INTO orders (customer_id, total_price, status, created_at, updated_at)
		VALUES (:customer_id, :total_price, :status, :created_at, :updated_at)
		RETURNING id
	`
}

func (order *Order) FeedGetByIdQuery() string {
	return `
		SELECT *
		FROM orders
		WHERE id = $1
	`
}

func (order *Order) FeedGetAllQuery() string {
	return `
		SELECT *
		FROM orders
	`
}

func (order *Order) FeedUpdateDetailsQuery() string {
	return `
		UPDATE orders
		SET customer_id = :customer_id, 
			total_price = :total_price,
			status = :status
		WHERE id = :id
	`
}

func (order *Order) FeedDeleteQuery() string {
	return `
		DELETE FROM orders
		WHERE id = $1
	`
}

func GetOrdersByCustomerIdQuery(ctx context.Context, db *sqlx.DB, customer_id int, orders *[]Order) error {
	query := `	
		SELECT *
		FROM orders 
		WHERE customer_id = $1
	`
	return db.Get(orders, query, customer_id)
}

