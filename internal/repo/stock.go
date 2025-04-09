package repo

import "time"

type Stock struct {
	Id				int			`db:"id" json:"id"`
	ProductId		int			`db:"product_id" json:"product_id"`
	Quantity		int			`db:"quantity" json:"quantity"`
	StockThreshold	int			`db:"stock_threshold" json:"stock_threshold"`
	CreatedAt		time.Time	`db:"created_at" json:"created_at"`
	UpdatedAt		time.Time	`db:"updated_at" json:"updated_at"`
}