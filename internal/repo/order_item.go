package repo

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type OrderItem struct {
	Id				int			`db:"id" json:"id"`
	OrderId			int			`db:"order_id" json:"order_id"`
	ProductId		int			`db:"product_id" json:"product_id"`
	Price			float64		`db:"price" json:"price"`
	Quantity		int			`db:"quantity" json:"quantity"`
	CreatedAt		time.Time	`db:"created_at" json:"created_at"`
	UpdatedAt		time.Time	`db:"updated_at" json:"updated_at"`
}

func (order_item *OrderItem) FeedGetId() *int {
	return &order_item.Id
}

func (order_item *OrderItem) FeedCreateQuery() string {
	return `
		INSERT INTO order_items (product_id, order_id, price, quantity, created_at, updated_at)
		VALUES (:product_id, :order_id, :price, :quantity, :created_at, :updated_at)
		RETURNING id
	`
}

func (order_item *OrderItem) FeedGetByIdQuery() string {
	return `
		SELECT id, product_id, order_id, price, quantity, created_at, updated_at
		FROM order_items
		WHERE id = $1
	`
}

func (order_item *OrderItem) FeedGetAllQuery() string {
	return `
		SELECT id, product_id, order_id, price, quantity, created_at, updated_at
		FROM order_items
	`
}

func (order_item *OrderItem) FeedUpdateDetailsQuery() string {
	return `
		UPDATE order_items
		SET product_id = :product_id, 
			order_id = :order_id, 
			price = :price,
			quantity = :quantity
		WHERE id = :id
		RETURNING id, product_id, order_id, price, quantity, created_at, updated_at
	`
}

func (order_item *OrderItem) FeedDeleteQuery() string {
	return `
		DELETE FROM order_items
		WHERE id = $1
	`
}

func GetOrderItemsByOrderIdQuery(ctx context.Context, db *sqlx.DB, order_id int, order_items *[]OrderItem) error {
	query := `	
		SELECT *
		FROM order_items 
		WHERE order_id = $1
	`
	return db.Get(order_items, query, order_id)
}