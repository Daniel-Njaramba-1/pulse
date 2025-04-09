package repo

import "time"

type Review struct {
	Id				int			`db:"id" json:"id"`
	CustomerId		int			`db:"customer_id" json:"customer_id"`
	ProductId		int			`db:"product_id" json:"product_id"`
	Rating			float32		`db:"rating" json:"rating"`
	ReviewText		string		`db:"review_text" json:"review_text"`
	CreatedAt		time.Time	`db:"created_at" json:"created_at"`
	UpdatedAt		time.Time	`db:"updated_at" json:"updated_at"`
}