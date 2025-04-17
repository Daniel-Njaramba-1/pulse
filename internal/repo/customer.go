package repo

import (
	"time"
)

type Customer struct {
	Id					int			`db:"id" json:"id"`
	Username			string		`db:"username" json:"username"`
	Email				string		`db:"email" json:"email"`
	PasswordHash		string		`db:"password_hash" json:"-"`
    Password        	string      `db:"-" json:"password"`
	IsActive			bool		`db:"is_active" json:"is_active"`
	CreatedAt			time.Time	`db:"created_at" json:"created_at"`
	UpdatedAt			time.Time	`db:"updated_at" json:"updated_at"`
}

type CustomerProfile struct {
	Id					int			`db:"id" json:"id"`
	CustomerId			int			`db:"customer_id" json:"customer_id"`
	FirstName			string		`db:"firstname" json:"firstname"`
	LastName			string		`db:"lastname" json:"lastname"`
	Phone				string		`db:"phone" json:"phone"`
	Address				string		`db:"address" json:"address"`
	CreatedAt			time.Time	`db:"created_at" json:"created_at"`
	UpdatedAt			time.Time	`db:"updated_at" json:"updated_at"`
}

type Cart struct {
	Id				int			`db:"id" json:"id"`
	CustomerId		int			`db:"customer_id" json:"customer_id"`
	IsActive		bool		`db:"is_active" json:"is_active"`
	CreatedAt		time.Time	`db:"created_at" json:"created_at"`
	UpdatedAt		time.Time	`db:"updated_at" json:"updated_at"`
}

type CartItem struct {
	Id				int			`db:"id" json:"id"`
	CartId			int			`db:"cart_id" json:"cart_id"`
	ProductId		int			`db:"product_id" json:"product_id"`
	Quantity		int			`db:"quantity" json:"quantity"`
	IsProcessed		bool		`db:"is_processed" json:"is_processed"`
	CreatedAt		time.Time	`db:"created_at" json:"created_at"`
	UpdatedAt		time.Time	`db:"updated_at" json:"updated_at"`
}

type Wishlist struct {
	Id         int       `db:"id" json:"id"`
	CustomerId int       `db:"customer_id" json:"customer_id"`
	IsActive   bool      `db:"is_active" json:"is_active"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

type WishlistItem struct {
	Id         int       `db:"id" json:"id"`
	WishlistId int       `db:"wishlist_id" json:"wishlist_id"`
	ProductID  int       `db:"product_id" json:"product_id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

type Review struct {
	Id				int			`db:"id" json:"id"`
	CustomerId		int			`db:"customer_id" json:"customer_id"`
	ProductId		int			`db:"product_id" json:"product_id"`
	Rating			float32		`db:"rating" json:"rating"`
	ReviewText		string		`db:"review_text" json:"review_text"`
	CreatedAt		time.Time	`db:"created_at" json:"created_at"`
	UpdatedAt		time.Time	`db:"updated_at" json:"updated_at"`
}


