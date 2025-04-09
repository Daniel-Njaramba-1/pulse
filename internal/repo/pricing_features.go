package repo

import "time"

type PricingFeatures struct {
	Id         			int       	`db:"id" json:"id"`
	ProductId  			int       	`db:"product_id" json:"product_id"`
	DemandScore 		float32		`db:"demand_score" json:"demand_score"`
	CompetitiveIndex	float32		`db:"competitive_index" json:"competitive_index"`
	SeasonalityFactor	float32		`db:"seasonality_factor" json:"seasonality_factor"`
	InventoryRatio		float32		`db:"inventory_ratio" json:"inventory_ratio"`
	DaysInStock			int			`db:"days_in_stock" json:"days_in_stock"`
	ViewToPurchaseRatio	float32		`db:"view_to_purchase_ratio" json:"view_to_purchase_ratio"`
	CreatedAt  			time.Time 	`db:"created_at" json:"created_at"`
	UpdatedAt  			time.Time 	`db:"updated_at" json:"updated_at"`
}