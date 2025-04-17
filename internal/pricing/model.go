package pricing

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Daniel-Njaramba-1/pulse/internal/repo"
	"github.com/jmoiron/sqlx"
)

type ModelService struct {
	db *sqlx.DB
}

func NewModelService(db *sqlx.DB) *ModelService {
	return &ModelService{db: db}
}

func (s *ModelService) TrainNewModel(ctx context.Context) error {
	// Get Training Data
	features, err := s.getAllPricingFeatures(ctx)
	if err != nil {
		return err
	}

	// Get historical prices
	prices, err := s.getHistoricalPrices(ctx)
	if err != nil {
		return err
	}

	// Create and train model
	model := NewPricingModel(generateModelVersion())
	err = model.Train(features, prices)
	if err != nil {
		return err
	}

	// Save model coefficients
	coeffs := model.GetCoefficients()
	return s.SaveModelCoefficients(ctx, coeffs)
}

func (s *ModelService) SaveModelCoefficients(ctx context.Context, coeffs repo.PriceModelCoefficients) error {
	query := `
        INSERT INTO price_model_coefficients (
            model_version, training_date, sample_size, r_squared,
            intercept, sales_count_coef, sales_value_coef, sales_velocity_coef,
            days_since_sale_coef, category_rank_coef, category_percentile_coef,
            review_score_coef, wishlist_ratio_coef, days_in_stock_coef, seasonal_factor_coef
        ) VALUES (
            :model_version, :training_date, :sample_size, :r_squared,
            :intercept, :sales_count_coef, :sales_value_coef, :sales_velocity_coef,
            :days_since_sale_coef, :category_rank_coef, :category_percentile_coef,
            :review_score_coef, :wishlist_ratio_coef, :days_in_stock_coef, :seasonal_factor_coef
        )
    `
	_, err := s.db.NamedExecContext(ctx, query, coeffs)
	return err
}

func (s *ModelService) AdjustPrice(ctx context.Context, productId int) (float64, float32, error) {
	// Get latest model coefficients
	coeffs, err := s.GetLatestModelCoefficients(ctx)
	if err != nil {
		return 0, 0, err
	}

	// Get current product metrics
	features, err := s.buildPricingFeatures(ctx, productId)
	if err != nil {
		return 0, 0, err
	}

	// Calculate new price using regression model
	newPrice := coeffs.Intercept +
		coeffs.SalesCountCoef*float64(features.TotalSalesCount) +
		coeffs.SalesValueCoef*features.TotalSalesValue +
		coeffs.SalesVelocityCoef*features.SalesVelocity +
		coeffs.DaysSinceSaleCoef*float64(features.DaysSinceLastSale) +
		coeffs.CategoryRankCoef*float64(features.CategoryRank) +
		coeffs.CategoryPercentileCoef*features.CategoryPercentile +
		coeffs.ReviewScoreCoef*features.ReviewScore +
		coeffs.WishlistRatioCoef*features.WishlistToSalesRatio +
		coeffs.DaysInStockCoef*float64(features.DaysInStock) +
		coeffs.SeasonalFactorCoef*features.SeasonalFactor

	// Get current price
	currentPrice, err := s.getCurrentPrice(ctx, productId)
	if err != nil {
		return 0, 0, err
	}

	// Calculate confidence score (simplified)
	confidenceScore := float32(0.85) // You might implement a proper confidence calculation

	// Save the price adjustment record
	err = s.savePriceAdjustment(ctx, productId, currentPrice, newPrice, coeffs.ModelVersion, confidenceScore)
	if err != nil {
		return 0, 0, err
	}

	// Update product metrics with new price
	err = s.updateProductPrice(ctx, productId, newPrice)
	if err != nil {
		return 0, 0, err
	}

	return newPrice, confidenceScore, nil
}

func (s *ModelService) GetTimeLastModelWasRun(ctx context.Context, productId int) (time.Time, error) {
	query := `
        SELECT last_model_run FROM pricing_features WHERE product_id = $1
    `
	var lastRun time.Time
	err := s.db.GetContext(ctx, &lastRun, query, productId)
	return lastRun, err
}

// Helper functions
func (s *ModelService) buildPricingFeatures(ctx context.Context, productId int) (repo.PricingFeatures, error) {
	var features repo.PricingFeatures
	features.ProductId = productId

	// Get all the required metrics
	daysSinceSale, err := s.GetDaysSinceLastSale(ctx, productId)
	if err != nil {
		return features, err
	}
	features.DaysSinceLastSale = daysSinceSale

	salesVelocity, err := s.CalculateSalesVelocity(ctx, productId)
	if err != nil {
		return features, err
	}
	features.SalesVelocity = salesVelocity

	salesCount, err := s.CalculateTotalSalesCount(ctx, productId)
	if err != nil {
		return features, err
	}
	features.TotalSalesCount = salesCount

	salesValue, err := s.CalculateTotalSalesValue(ctx, productId)
	if err != nil {
		return features, err
	}
	features.TotalSalesValue = salesValue

	categoryRank, err := s.CalculateCategoryRank(ctx, productId)
	if err != nil {
		return features, err
	}
	features.CategoryRank = categoryRank

	categoryPercentile, err := s.CalculateCategoryPercentile(ctx, productId)
	if err != nil {
		return features, err
	}
	features.CategoryPercentile = categoryPercentile

	reviewScore, err := s.CalculateReviewScore(ctx, productId)
	if err != nil {
		return features, err
	}
	features.ReviewScore = reviewScore

	wishlistRatio, err := s.CalculateWishlistToSalesRatio(ctx, productId)
	if err != nil {
		return features, err
	}
	features.WishlistToSalesRatio = wishlistRatio

	daysInStock, err := s.CalculateDaysInStock(ctx, productId)
	if err != nil {
		return features, err
	}
	features.DaysInStock = daysInStock

	seasonalFactor, err := s.CalculateSeasonalFactor(ctx, productId)
	if err != nil {
		return features, err
	}
	features.SeasonalFactor = seasonalFactor

	// Set last model run time
	features.LastModelRun = time.Now()

	return features, nil
}

func (s *ModelService) GetLatestModelCoefficients(ctx context.Context) (repo.PriceModelCoefficients, error) {
	query := `
        SELECT * FROM price_model_coefficients
        ORDER BY created_at DESC
        LIMIT 1
    `
	var coeffs repo.PriceModelCoefficients
	err := s.db.GetContext(ctx, &coeffs, query)
	return coeffs, err
}

func (s *ModelService) getCurrentPrice(ctx context.Context, productId int) (float64, error) {
	query := `
        SELECT adjusted_price FROM product_metrics WHERE product_id = $1
    `
	var price float64
	err := s.db.GetContext(ctx, &price, query, productId)
	return price, err
}

func (s *ModelService) updateProductPrice(ctx context.Context, productId int, newPrice float64) error {
	query := `
        UPDATE product_metrics
        SET adjusted_price = $1, last_price_update = NOW()
        WHERE product_id = $2
    `
	_, err := s.db.ExecContext(ctx, query, newPrice, productId)
	return err
}

func (s *ModelService) savePriceAdjustment(ctx context.Context, productId int, oldPrice, newPrice float64, modelVersion string, confidenceScore float32) error {
	query := `
        INSERT INTO price_adjustments (
            product_id, old_price, new_price, model_version, confidence_score
        ) VALUES ($1, $2, $3, $4, $5)
    `
	_, err := s.db.ExecContext(ctx, query, productId, oldPrice, newPrice, modelVersion, confidenceScore)
	return err
}

func (s *ModelService) getAllPricingFeatures(ctx context.Context) ([]repo.PricingFeatures, error) {
	query := `SELECT * FROM pricing_features`
	var features []repo.PricingFeatures
	err := s.db.SelectContext(ctx, &features, query)
	return features, err
}

func (s *ModelService) getHistoricalPrices(ctx context.Context) ([]float64, error) {
	query := `
        SELECT adjusted_price FROM product_metrics
        ORDER BY product_id
    `
	var prices []float64
	err := s.db.SelectContext(ctx, &prices, query)
	return prices, err
}

func generateModelVersion() string {
	return fmt.Sprintf("v%s", time.Now().Format("20060102150405"))
}

func (s *ModelService) AdjustAllPrices(ctx context.Context) error {
	// Get all active product IDs
	query := `SELECT id FROM products WHERE is_active = true`
	var productIDs []int
	err := s.db.SelectContext(ctx, &productIDs, query)
	if err != nil {
		return err
	}

	// Process in batches to avoid overwhelming the database
	batchSize := 50
	for i := 0; i < len(productIDs); i += batchSize {
		end := i + batchSize
		if end > len(productIDs) {
			end = len(productIDs)
		}

		batch := productIDs[i:end]
		for _, productID := range batch {
			_, _, err := s.AdjustPrice(ctx, productID)
			if err != nil {
				// Log the error but continue with other products
				log.Printf("Error adjusting price for product %d: %v", productID, err)
			}
		}

		// Sleep briefly between batches to reduce database load
		time.Sleep(500 * time.Millisecond)
	}

	return nil
}

func (s *ModelService) SchedulePriceAdjustments(frequency time.Duration) {
	ticker := time.NewTicker(frequency)
	go func() {
		for range ticker.C {
			s.AdjustAllPrices(context.Background())
		}
	}()
}
