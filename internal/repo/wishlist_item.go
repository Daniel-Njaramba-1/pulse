package repo

import "time"

type WishlistItem struct {
	Id         int       `db:"id" json:"id"`
	WishlistId int       `db:"wishlist_id" json:"wishlist_id"`
	ProductID  int       `db:"product_id" json:"product_id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}