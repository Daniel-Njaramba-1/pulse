package repo

import "time"

type PriceModelCoefficients struct {
	Id         			int       	`db:"id" json:"id"`
	ProductCategoryId  	int       	`db:"product_category_id" json:"product_category_id"`
	FeatureName			string		`db:"feature_name" json:"feature_name"`
	Coefficient			float32		`db:"coefficient" json:"coefficient"`
	LastTrainedAt		time.Time	`db:"last_trained_at" json:"last_trained_at"`
	CreatedAt  			time.Time 	`db:"created_at" json:"created_at"`
	UpdatedAt  			time.Time 	`db:"updated_at" json:"updated_at"`
}