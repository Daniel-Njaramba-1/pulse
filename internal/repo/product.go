
package repo

import (
	"time"
)

type Product struct {
	Id				int			`db:"id" json:"id"`
	CategoryId	    int			`db:"category_id" json:"category_id"`
	BrandId			int			`db:"brand_id" json:"brand_id"`
	Name			string		`db:"name" json:"name"`
	Description		string		`db:"description" json:"description"`
	ImagePath		*string		`db:"image_path" json:"image_path"` //pointer to handle null values
	IsActive		bool		`db:"is_active" json:"is_active"` 
	CreatedAt		time.Time	`db:"created_at" json:"created_at"`
	UpdatedAt		time.Time	`db:"updated_at" json:"updated_at"`
}

type ProductMetric struct {
	Id				int			`db:"id" json:"id"`
	ProductId		int			`db:"product_id" json:"product_id"`
	AverageRating	float64		`db:"average_rating" json:"average_rating"`
	ReviewCount		int			`db:"review_count" json:"review_count"`
	WishlistCount	int			`db:"wishlist_count" json:"wishlist_count"`
	BasePrice		float64		`db:"base_price" json:"base_price"`
	AdjustedPrice	float64		`db:"adjusted_price" json:"adjusted_price"`
	LastSale		time.Time	`db:"last_sale" json:"last_sale"`
	LastPriceUpdate	time.Time	`db:"last_price_update" json:"last_price_update"`
	CreatedAt		time.Time	`db:"created_at" json:"created_at"`
	UpdatedAt		time.Time	`db:"updated_at" json:"updated_at"`
}

type Stock struct {
	Id					int			`db:"id" json:"id"`
	ProductId			int			`db:"product_id" json:"product_id"`
	Quantity			int			`db:"quantity" json:"quantity"`
	StockThreshold		int			`db:"stock_threshold" json:"stock_threshold"`
	FirstStockedDate	time.Time	`db:"first_stocked_date" json:"first_stocked_date"` 
	LastOutOfStockDate	time.Time	`db:"last_out_of_stock_date" json:"last_out_of_stock_date"`
	CreatedAt			time.Time	`db:"created_at" json:"created_at"`
	UpdatedAt			time.Time	`db:"updated_at" json:"updated_at"`
}

type StockingEvent string 

const (
	StockingEventInStock		StockingEvent = "in stock"
	StockingEventOutOfStock 	StockingEvent = "out of stock"
	StockingEventRestock 		StockingEvent = "restock"
)

type StockHistory struct {
	Id            	int       		`db:"id" json:"id"`
	ProductId     	int       		`db:"product_id" json:"product_id"`
	EventType     	StockingEvent   `db:"event_type" json:"event_type"`
	QuantityChange 	int      		`db:"quantity_change" json:"quantity_change"`
	QuantityAfter 	int       		`db:"quantity_after" json:"quantity_after"`
	EventDate     	time.Time 		`db:"event_date" json:"event_date"`
	CreatedAt     	time.Time 		`db:"created_at" json:"created_at"`
	UpdatedAt     	time.Time 		`db:"updated_at" json:"updated_at"`
}

type Sale struct {
	Id				int			`db:"id" json:"id"`
	OrderItemId		int			`db:"order_item_id" json:"order_item_id"`
	ProductID		int			`db:"product_id" json:"product_id"`
	SalePrice		float64		`db:"sale_price" json:"sale_price"`
	Quantity		int			`db:"quantity" json:"quantity"`
	CreatedAt		time.Time	`db:"created_at" json:"created_at"`
	UpdatedAt		time.Time	`db:"updated_at" json:"updated_at"`
}

// ProductDetail represents a detailed view of a product.
type ProductDetail struct {
    // Basic product information
    Id          int       `db:"id" json:"id"`
    Name        string    `db:"name" json:"name"`
    Description string    `db:"description" json:"description"`
    ImagePath   *string    `db:"image_path" json:"image_path"`
    IsActive    bool      `db:"is_active" json:"is_active"`
    // CreatedAt   time.Time `db:"created_at" json:"created_at"`
    // UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`

    // Associated brand information
    BrandId   int    `db:"brand_id" json:"brand_id"`
    BrandName string `db:"brand_name" json:"brand_name"`

    // Associated category information
    CategoryId   int    `db:"category_id" json:"category_id"`
    CategoryName string `db:"category_name" json:"category_name"`

    // Product metrics
    AverageRating  *float32 `db:"average_rating" json:"average_rating"`
    ReviewCount    *int     `db:"review_count" json:"review_count"`
    WishlistCount  *int     `db:"wishlist_count" json:"wishlist_count"`
    BasePrice      *float32 `db:"base_price" json:"base_price"`
    AdjustedPrice  *float32 `db:"adjusted_price" json:"adjusted_price"`

    // Stock information
    StockQuantity    *int `db:"stock_quantity" json:"stock_quantity"`
    StockThreshold   *int `db:"stock_threshold" json:"stock_threshold"`
}
