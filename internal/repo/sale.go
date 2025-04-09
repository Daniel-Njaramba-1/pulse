package repo

import "time"

type Sale struct {
	Id				int			`db:"id" json:"id"`
	OrderItemId		int			`db:"order_item_id" json:"order_item_id"`
	ProductID		int			`db:"product_id" json:"product_id"`
	SalePrice		float32		`db:"sale_price" json:"sale_price"`
	Quantity		int			`db:"quantity" json:"quantity"`
	CreatedAt		time.Time	`db:"created_at" json:"created_at"`
	UpdatedAt		time.Time	`db:"updated_at" json:"updated_at"`
}