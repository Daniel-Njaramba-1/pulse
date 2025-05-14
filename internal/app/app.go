package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Daniel-Njaramba-1/pulse/internal/db"
	"github.com/Daniel-Njaramba-1/pulse/internal/pricing"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// App represents the application with its dependencies
type App struct {
	db            *sqlx.DB
	echo          *echo.Echo
	pricingModel  *pricing.ModelService
	dbConfig      *db.DBConfig
}

// NewApp initializes and returns a new app instance
func NewApp() (*App, error) {
	ctx := context.Background()
	
	dbConfig, err := db.LoadDBConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load DB config: %w", err)
	}
	
	database, err := db.InitDB(ctx, dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Start client manager in a goroutine
	go db.Manager.Run()
	
	// Start price adjustment listener in a goroutine
	connStr := db.BuildConnStr(dbConfig)
	go db.StartPriceAdjustmentListener(connStr)

	// Initialize Echo framework
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5190", "http://localhost:5195"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderCookie},
		AllowCredentials: true,
	}))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Set up SSE endpoint
	e.GET("/api/price-adjustments", HandleSSE)
	
	// Initialize pricing model
	pricingModel := pricing.NewModelService(database)
	if err := InitializePricingModel(ctx, pricingModel); err != nil {
		log.Printf("Error initializing pricing model: %v", err)
	}

	// Set up service handlers
	adminServices := NewAdminServices(database)
	customerServices := NewCustomerServices(database)

	adminHandlers := NewAdminHdl(adminServices)
	customerHandlers := NewCustomerHdl(customerServices)

	AdminRoutes(e, adminHandlers)
	CustomerRoutes(e, customerHandlers)

	// Set up static file serving for product images using Go's built-in file server
    rootDir, err := os.Getwd()
    if err != nil {
        log.Printf("Error getting current working directory: %v", err)
    } else {
        // Create a file server handler
        productImagesPath := filepath.Join(rootDir, "internal", "assets", "products")
        fs := http.FileServer(http.Dir(productImagesPath))
        
        // Register the handler with Echo
        e.GET("/assets/products/*", echo.WrapHandler(http.StripPrefix("/assets/products/", fs)))
        log.Printf("Serving static files from: %s", productImagesPath)
    }
	
	// Start Price Adjustment
	go db.StartSaleListener(connStr, pricingModel)
	// Start schedulers
	go StartPriceAdjustmentScheduler(pricingModel)
	go StartModelTrainingScheduler(pricingModel)
	
	return &App{
		db:           database,
		echo:         e,
		pricingModel: pricingModel,
		dbConfig:     dbConfig,
	}, nil
}

// Start begins listening for HTTP requests
func (a *App) Start() error {
	return a.echo.Start(":8080")
}

// Close properly shuts down the application
func (a *App) Close() {
	if a.db != nil {
		db.CloseDB(a.db)
	}
}

// GetDB returns the database connection
func (a *App) GetDB() *sqlx.DB {
	return a.db
}



// HandleSSE handles Server-Sent Events connections
func HandleSSE(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "text/event-stream")
    c.Response().Header().Set("Cache-Control", "no-cache")
    c.Response().Header().Set("Connection", "keep-alive")
    c.Response().Header().Set("Access-Control-Allow-Origin", "*")

	c.Response().WriteHeader(http.StatusOK)

	initialData := map[string]string{
        "type": "connection",
        "status": "established",
        "timestamp": fmt.Sprintf("%d", time.Now().Unix()),
    }
    initialJSON, _ := json.Marshal(initialData)
    fmt.Fprintf(c.Response().Writer, "event: connect\ndata: %s\n\n", initialJSON)
    c.Response().Flush()
	log.Printf("SSE client connected from %s", c.Request().RemoteAddr)

    // Each client gets its own channel
    messageChan := make(chan string)
    
    // Register this client
    db.Manager.RegisterChannel(messageChan)
    
    // Ensure client is unregistered when connection closes
    defer func() {
        db.Manager.UnregisterChannel(messageChan)
    }()

	// Send a heartbeat every 30 seconds to keep connection alive
	heartbeat := time.NewTicker(30 * time.Second)
    defer heartbeat.Stop()
    
    // Keep connection open
    for {
        select {
        case msg := <-messageChan:
            // Format as SSE
            if _, err := fmt.Fprintf(c.Response().Writer, "data: %s\n\n", msg); err != nil {
                return err
            }
            c.Response().Flush()
		case <-heartbeat.C:
            // Send heartbeat comment to keep connection alive
            if _, err := fmt.Fprintf(c.Response().Writer, ": heartbeat %v\n\n", time.Now().Unix()); err != nil {
                log.Printf("Error sending SSE heartbeat: %v", err)
                return err
            }
            c.Response().Flush()
        case <-c.Request().Context().Done():
            log.Printf("SSE client connection closed")
			return nil
        }
    }
}