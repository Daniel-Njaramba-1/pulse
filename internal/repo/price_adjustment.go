package repo

import "time"

type PriceAdjustment struct {
	Id				int			`db:"id" json:"id"`
	ProductId		int			`db:"product_id" json:"product_id"`
	OldPrice		float32		`db:"old_price" json:"old_price"`
	NewPrice		float32		`db:"new_price" json:"new_price"`
	ModelVersion	string		`db:"model_version" json:"model_version"`
	ConfidenceScore	float32		`db:"confidence_score" json:"confidence_score"`
	CreatedAt		time.Time	`db:"created_at" json:"created_at"`
	UpdatedAt		time.Time	`db:"updated_at" json:"updated_at"`
}