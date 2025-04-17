package app

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Daniel-Njaramba-1/pulse/internal/db"
	"github.com/Daniel-Njaramba-1/pulse/internal/pricing"
	"github.com/Daniel-Njaramba-1/pulse/internal/repo"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type App struct {
	db *sqlx.DB
	echo *echo.Echo
	pricingModel *pricing.ModelService
}

func NewApp() (*App, error) {
	db, err := db.ConnDB()
	if err != nil {
		return nil, err
	}

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "http://localhost:5174"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	
	pricingModel := pricing.NewModelService(db)
	err = initializePricingModel(pricingModel)
	if err != nil {
		log.Printf("Error initializing pricing model: %v", err)
	}

	adminServices := NewAdminServices(db)
	customerServices := NewCustomerServices(db)

	adminHandlers := NewAdminHdl(adminServices)
	customerHandlers := NewCustomerHdl(customerServices)

	AdminRoutes(e, adminHandlers)
	CustomerRoutes(e, customerHandlers)

	// product images are in filepath.Join (os.Getwd() "internal", "assets", "products")
	rootDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}
	e.Static("/assets/products", filepath.Join(rootDir, "internal", "assets", "products"))
	
	// start price adjustment
	go startPriceAdjustmentScheduler(pricingModel)
	
	return &App{
		db: db,
		echo: e,
		pricingModel: pricingModel,
	}, nil
}

func (a *App) Start() error {
	return a.echo.Start(":8080")
}

func (a *App) Close () {
	if a.db != nil {
		a.db.Close()
	}
}

func initializePricingModel(modelService *pricing.ModelService) error {
	ctx := context.Background()

	coeffs, err := modelService.GetLatestModelCoefficients(ctx)
	if err != nil {
		log.Printf("No existing model coefficients found, initializing with seed data")

		seedCoeffs := repo.PriceModelCoefficients{
			ModelVersion: "v1.0",
			TrainingDate: time.Now(),
			SampleSize: 0,
			RSquared: 0.0,
			Intercept: 100.0,
			SalesCountCoef: 0.5,
			SalesValueCoef: 0.01,
			SalesVelocityCoef: 2.0,
			DaysSinceSaleCoef: -0.1,
			CategoryRankCoef: -0.5,
			CategoryPercentileCoef: 0.3,
			ReviewScoreCoef: 5.0,
			WishlistRatioCoef: 2.0,
			DaysInStockCoef: -0.05,
			SeasonalFactorCoef: 10.0,
		}

		err = modelService.SaveModelCoefficients(ctx, seedCoeffs)
		if err != nil {
			return err
		}

		log.Printf("Succesfully initiated pricing model with seed data")
		return nil
	}

	log.Printf("Found existing model coefficients (version: %s)", coeffs.ModelVersion)
	return nil
}

// Start a background goroutine to periodically adjust prices
func startPriceAdjustmentScheduler(modelService *pricing.ModelService) {
	ticker := time.NewTicker(24 * time.Hour) // Adjust prices once per day
	go func() {
		log.Printf("Starting price adjustment scheduler")
		for range ticker.C {
			ctx := context.Background()
			log.Printf("Running scheduled price adjustments")
			
			// Run price adjustments for all active products
			err := modelService.AdjustAllPrices(ctx)
			if err != nil {
				log.Printf("Error during scheduled price adjustment: %v", err)
			} else {
				log.Printf("Scheduled price adjustments completed successfully")
			}
		}
	}()
}


