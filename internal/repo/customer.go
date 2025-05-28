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
	IsEmailVerified 	bool        `db:"is_email_verified"`
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

type CartItemDetail struct {
	Id				int			`db:"id" json:"id"`
	CartId			int			`db:"cart_id" json:"cart_id"`
	ProductId		int			`db:"product_id" json:"product_id"`
	Quantity		int			`db:"quantity" json:"quantity"`
	IsProcessed		bool		`db:"is_processed" json:"is_processed"`
	CreatedAt		time.Time	`db:"created_at" json:"created_at"`
	UpdatedAt		time.Time	`db:"updated_at" json:"updated_at"`

	ProductName				string		`db:"product_name" json:"product_name"`
	ProductImagePath		*string		`db:"product_image_path" json:"product_image_path"`
	ProductAdjustedPrice  	*float32 	`db:"product_adjusted_price" json:"product_adjusted_price"`

	ProductStockQuantity    *int `db:"product_stock_quantity" json:"product_stock_quantity"`
}


// CartWithItems represents a cart with its items
type CartWithItems struct {
	Id         int        `json:"id" db:"id"`
	CustomerId int        `json:"customer_id" db:"customer_id"`
	IsActive   bool       `json:"is_active" db:"is_active"`
	
	Items      []CartItemDetail `json:"items"`
	
	// Summary data
	TotalItems     int     `json:"total_items"`
	TotalPrice     float64 `json:"total_price"`
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

type WishlistItemDetail struct {
	// Basic product information
    Id          int       	`db:"id" json:"id"`
    WishlistId  int			`db:"wishlist_id" json:"wishlist_id"`
	ProductId	int			`db:"product_id" json:"product_id"`

	ProductName				string		`db:"product_name" json:"product_name"`
	ProductImagePath		*string		`db:"product_image_path" json:"product_image_path"`
	ProductAdjustedPrice  	*float32 	`db:"product_adjusted_price" json:"product_adjusted_price"`

    ProductStockQuantity    *int `db:"product_stock_quantity" json:"product_stock_quantity"`
}

type WishlistDetail struct {
	Id         int        `json:"id" db:"id"`
	CustomerId int        `json:"customer_id" db:"customer_id"`
	IsActive   bool       `json:"is_active" db:"is_active"`
	
	Items      []WishlistItemDetail `json:"items"`
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

type ReviewDetail struct {
	Id				int			`db:"id" json:"id"`
	CustomerId		int			`db:"customer_id" json:"customer_id"`
	CustomerName 	string		`db:"username" json:"customer_name"`
	ProductId		int			`db:"product_id" json:"product_id"`
	Rating			float32		`db:"rating" json:"rating"`
	ReviewText		string		`db:"review_text" json:"review_text"`
	CreatedAt		time.Time	`db:"created_at" json:"created_at"`
	UpdatedAt		time.Time	`db:"updated_at" json:"updated_at"`
}


