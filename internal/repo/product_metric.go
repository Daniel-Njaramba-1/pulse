package repo

import "time"

type ProductMetric struct {
	Id				int			`db:"id" json:"id"`
	ProductId		int			`db:"product_id" json:"product_id"`
	AverageRating	float32		`db:"average_rating" json:"average_rating"`
	ReviewCount		int			`db:"review_count" json:"review_count"`
	WishlistCount	int			`db:"wishlist_count" json:"wishlist_count"`
	BasePrice		float32		`db:"base_price" json:"base_price"`
	AdjustedPrice	float32		`db:"adjusted_price" json:"adjusted_price"`
	LastPrice 		time.Time	`db:"last_price" json:"last_price"`
	CreatedAt		time.Time	`db:"created_at" json:"created_at"`
	UpdatedAt		time.Time	`db:"updated_at" json:"updated_at"`
}

func (pm *ProductMetric) FeedGetId() *int {
	return &pm.Id
}

func (pm *ProductMetric) FeedCreateQuery() string {
	return `
		INSERT INTO product_metrics (categoryId, brandId, name, description, imagePath, is_active)
		VALUES (:categoryId, :brandId, :name, :description, :imagePath, :is_active)
		RETURNING id
	`
}

func (pm *ProductMetric) FeedGetByIdQuery() string {
	return `
		SELECT * FROM products
		WHERE id = :id
	`
}

func (pm *ProductMetric) FeedGetAllQuery() string {
	return `
		SELECT * FROM products
		FROM products
		ORDER BY id ASC
	`
}

func (pm *ProductMetric) FeedUpdateDetailsQuery() string {
	return `
		UPDATE products
		SET categoryId = :categoryId, 
			brandId = :brandId, 
			name = :name, 
			description = :description, 
			imagePath = :imagePath, 
		WHERE id = :id
	`
}

func (pm *ProductMetric) FeedDeleteQuery() string {
	return `
		DELETE FROM products
		WHERE id = :id
	`
}

func (pm *ProductMetric) FeedDeactivateQuery() string {
	return `
		UPDATE products
		SET is_active = false
		WHERE id = :id
	`
}

func (pm *ProductMetric) FeedReactivateQuery() string {
	return `
		UPDATE products
		SET is_active = true
		WHERE id = :id
	`
}
