package repo

import "time"

type PricingFeatures struct {
	Id                 		int       	`db:"id" json:"id"`
	ProductId          		int       	`db:"product_id" json:"product_id"`
	DaysSinceLastSale  		int       	`db:"days_since_last_sale" json:"days_since_last_sale"`
	SalesVelocity      		float64   	`db:"sales_velocity" json:"sales_velocity"`
	TotalSalesCount    		int       	`db:"total_sales_count" json:"total_sales_count"`
	TotalSalesValue    		float64   	`db:"total_sales_value" json:"total_sales_value"`
	CategoryRank       		int       	`db:"category_rank" json:"category_rank"`
	CategoryPercentile 		float64   	`db:"category_percentile" json:"category_percentile"`
	ReviewScore        		float64   	`db:"review_score" json:"review_score"`
	WishlistToSalesRatio 	float64 	`db:"wishlist_to_sales_ratio" json:"wishlist_to_sales_ratio"`
	DaysInStock        		int       	`db:"days_in_stock" json:"days_in_stock"`
	SeasonalFactor     		float64   	`db:"seasonal_factor" json:"seasonal_factor"`
	LastModelRun       		time.Time 	`db:"last_model_run" json:"last_model_run"`
	CreatedAt          		time.Time 	`db:"created_at" json:"created_at"`
	UpdatedAt          		time.Time 	`db:"updated_at" json:"updated_at"`
}

type PriceAdjustment struct {
	Id				int			`db:"id" json:"id"`
	ProductId		int			`db:"product_id" json:"product_id"`
	OldPrice		float64		`db:"old_price" json:"old_price"`
	NewPrice		float64		`db:"new_price" json:"new_price"`
	ModelVersion	string		`db:"model_version" json:"model_version"`
	ConfidenceScore	float32		`db:"confidence_score" json:"confidence_score"`
	CreatedAt		time.Time	`db:"created_at" json:"created_at"`
	UpdatedAt		time.Time	`db:"updated_at" json:"updated_at"`
}

type PriceModelCoefficients struct {
	Id                    	int       	`db:"id" json:"id"`
	ModelVersion          	string    	`db:"model_version" json:"model_version"`
	TrainingDate          	time.Time 	`db:"training_date" json:"training_date"`
	SampleSize            	int       	`db:"sample_size" json:"sample_size"`
	RSquared              	float64   	`db:"r_squared" json:"r_squared"`
	Intercept             	float64   	`db:"intercept" json:"intercept"`
	SalesCountCoef        	float64   	`db:"sales_count_coef" json:"sales_count_coef"`
	SalesValueCoef        	float64   	`db:"sales_value_coef" json:"sales_value_coef"`
	SalesVelocityCoef     	float64   	`db:"sales_velocity_coef" json:"sales_velocity_coef"`
	DaysSinceSaleCoef     	float64   	`db:"days_since_sale_coef" json:"days_since_sale_coef"`
	CategoryRankCoef      	float64   	`db:"category_rank_coef" json:"category_rank_coef"`
	CategoryPercentileCoef 	float64  	`db:"category_percentile_coef" json:"category_percentile_coef"`
	ReviewScoreCoef       	float64   	`db:"review_score_coef" json:"review_score_coef"`
	WishlistRatioCoef     	float64   	`db:"wishlist_ratio_coef" json:"wishlist_ratio_coef"`
	DaysInStockCoef       	float64   	`db:"days_in_stock_coef" json:"days_in_stock_coef"`
	SeasonalFactorCoef    	float64   	`db:"seasonal_factor_coef" json:"seasonal_factor_coef"`
	CreatedAt            	time.Time 	`db:"created_at" json:"created_at"`
	UpdatedAt             	time.Time 	`db:"updated_at" json:"updated_at"`
}