package app

import (
	"log"
	"os"
	"path/filepath"

	"github.com/Daniel-Njaramba-1/pulse/internal/db"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type App struct {
	db *sqlx.DB
	echo *echo.Echo

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
	
	return &App{
		db: db,
		echo: e,
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


