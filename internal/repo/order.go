package repo

import (
	"time"
)

type OrderStatus string

const (
	OrderStatusPending		OrderStatus = "pending"
	OrderStatusCancelled	OrderStatus = "cancelled"
	OrderStatusFailed 		OrderStatus = "failed"
	OrderStatusCompleted 	OrderStatus = "completed"
)

type Order struct {
	Id				int			`db:"id" json:"id"`
	CustomerId		int			`db:"customer_id" json:"customer_id"`
	TotalPrice		float64		`db:"total_price" json:"total_price"`
	Status			OrderStatus	`db:"status" json:"status"`
	PriceValidUntil time.Time	`db:"price_valid_until" json:"price_valid_until"`
	CreatedAt		time.Time	`db:"created_at" json:"created_at"`
	UpdatedAt		time.Time	`db:"updated_at" json:"updated_at"`
}

type OrderItem struct {
	Id				int			`db:"id" json:"id"`
	OrderId			int			`db:"order_id" json:"order_id"`
	ProductId		int			`db:"product_id" json:"product_id"`
	Price			float64		`db:"price" json:"price"`
	Quantity		int			`db:"quantity" json:"quantity"`
	CreatedAt		time.Time	`db:"created_at" json:"created_at"`
	UpdatedAt		time.Time	`db:"updated_at" json:"updated_at"`
}
