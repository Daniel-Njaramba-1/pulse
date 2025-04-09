package app

import "github.com/labstack/echo/v4"

func AdminRoutes(e *echo.Echo, adminHandlers *AdminHdl) {
    admin := e.Group("/api/admin")

    // Authentication routes
    admin.POST("/register", func(c echo.Context) error {
        return adminHandlers.AuthHandler.Register(c)
    })
    admin.POST("/login", func(c echo.Context) error {
        return adminHandlers.AuthHandler.Login(c)
    })

    // Brand routes
    admin.GET("/brands", func(c echo.Context) error {
        return adminHandlers.BrandHandler.GetAllBrands(c)
    })
    admin.GET("/brands/:id", func(c echo.Context) error {
        return adminHandlers.BrandHandler.GetBrandByID(c)
    })
    admin.POST("/brands", func(c echo.Context) error {
        return adminHandlers.BrandHandler.CreateBrand(c)
    })
    admin.PUT("/brands/:id", func(c echo.Context) error {
        return adminHandlers.BrandHandler.UpdateBrand(c)
    })
    admin.DELETE("/brands/:id", func(c echo.Context) error {
        return adminHandlers.BrandHandler.DeleteBrand(c)
    })
    admin.PUT("/brands/:id/deactivate", func(c echo.Context) error {
        return adminHandlers.BrandHandler.DeactivateBrand(c)
    })
    admin.PUT("/brands/:id/reactivate", func(c echo.Context) error {
        return adminHandlers.BrandHandler.ReactivateBrand(c)
    })

    // Category routes
    admin.GET("/categories", func(c echo.Context) error {
        return adminHandlers.CategoryHandler.GetAllCategories(c)
    })
    admin.GET("/categories/:id", func(c echo.Context) error {
        return adminHandlers.CategoryHandler.GetCategoryByID(c)
    })
    admin.POST("/categories", func(c echo.Context) error {
        return adminHandlers.CategoryHandler.CreateCategory(c)
    })
    admin.PUT("/categories/:id", func(c echo.Context) error {
        return adminHandlers.CategoryHandler.UpdateCategory(c)
    })
    admin.DELETE("/categories/:id", func(c echo.Context) error {
        return adminHandlers.CategoryHandler.DeleteCategory(c)
    })
    admin.PUT("/categories/:id/deactivate", func(c echo.Context) error {
        return adminHandlers.CategoryHandler.DeactivateCategory(c)
    })
    admin.PUT("/categories/:id/reactivate", func(c echo.Context) error {
        return adminHandlers.CategoryHandler.ReactivateCategory(c)
    })

    // Product routes
    admin.GET("/products", func(c echo.Context) error {
        return adminHandlers.ProductHandler.GetAllProducts(c)
    })
    admin.GET("/products/:id", func(c echo.Context) error {
        return adminHandlers.ProductHandler.GetProductByID(c)
    })
    admin.POST("/products", func(c echo.Context) error {
        return adminHandlers.ProductHandler.CreateProduct(c)
    })
    admin.PUT("/products/:id", func(c echo.Context) error {
        return adminHandlers.ProductHandler.UpdateProduct(c)
    })
    admin.DELETE("/products/:id", func(c echo.Context) error {
        return adminHandlers.ProductHandler.DeleteProduct(c)
    })
    admin.PUT("/products/:id/deactivate", func(c echo.Context) error {
        return adminHandlers.ProductHandler.DeactivateProduct(c)
    })
    admin.PUT("/products/:id/reactivate", func(c echo.Context) error {
        return adminHandlers.ProductHandler.ReactivateProduct(c)
    })
}

