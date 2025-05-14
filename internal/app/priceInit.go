package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Daniel-Njaramba-1/pulse/internal/pricing"
	"github.com/Daniel-Njaramba-1/pulse/internal/repo"
	"github.com/robfig/cron/v3"
)

// GetPricingModel returns the pricing model service
func (a *App) GetPricingModel() *pricing.ModelService {
	return a.pricingModel
}

// initializePricingModel sets up the initial pricing model coefficients if needed
func InitializePricingModel(ctx context.Context, modelService *pricing.ModelService) error {
	coeffs, err := modelService.GetLatestModelCoefficients(ctx)
	if err != nil {
		log.Printf("No existing model coefficients found, initializing with seed data")

		seedCoeffs := repo.PriceModelCoefficients{
			ModelVersion:          "v1.0",
			TrainingDate:          time.Now(),
			SampleSize:            0,
			RSquared:              0.0,
			Intercept:             100.0,
			SalesCountCoef:        0.5,
			SalesValueCoef:        0.01,
			SalesVelocityCoef:     2.0,
			DaysSinceSaleCoef:     -0.1,
			CategoryRankCoef:      -0.5,
			CategoryPercentileCoef: 0.3,
			ReviewScoreCoef:       5.0,
			WishlistRatioCoef:     2.0,
			DaysInStockCoef:       -0.05,
			SeasonalFactorCoef:    10.0,
		}

		if err := modelService.SaveModelCoefficients(ctx, seedCoeffs); err != nil {
			return fmt.Errorf("failed to save seed coefficients: %w", err)
		}

		log.Printf("Successfully initiated pricing model with seed data")
		return nil
	}

	log.Printf("Found existing model coefficients (version: %s)", coeffs.ModelVersion)
	return nil
}

// startPriceAdjustmentScheduler runs price adjustments on a daily schedule
func StartPriceAdjustmentScheduler(modelService *pricing.ModelService) {
	ticker := time.NewTicker(24 * time.Hour) // Adjust prices once per day
	defer ticker.Stop()
	
	log.Printf("Starting price adjustment scheduler")
	for range ticker.C {
		ctx := context.Background()
		log.Printf("Running scheduled price adjustments")
		
		// Run price adjustments for all active products
		if err := modelService.AdjustAllPrices(ctx); err != nil {
			log.Printf("Error during scheduled price adjustment: %v", err)
		} else {
			log.Printf("Scheduled price adjustments completed successfully")
		}
	}
}

// startModelTrainingScheduler runs model training on a monthly schedule
func StartModelTrainingScheduler(modelService *pricing.ModelService) {
    c := cron.New()
    defer c.Stop()
    
    _, err := c.AddFunc("0 0 1 * *", func() { // Run at 00:00 on the 1st of every month
        ctx := context.Background()
        log.Printf("Running scheduled monthly model training")
        
        if err := modelService.TrainNewModel(ctx); err != nil {
            log.Printf("Error during scheduled monthly model training: %v", err)
        } else {
            log.Printf("Scheduled monthly model training completed successfully")
        }
    })
    
    if err != nil {
        log.Printf("Error setting up monthly model training schedule: %v", err)
        return
    }
    
    c.Start()
    log.Printf("Started monthly model training scheduler")
}