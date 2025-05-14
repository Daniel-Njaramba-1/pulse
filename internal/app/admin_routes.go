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

    protected := admin.Group("")
    protected.Use(AdminAuthMiddleware())

    // Brand routes
    protected.GET("/brands", func(c echo.Context) error {
        return adminHandlers.BrandHandler.GetAllBrands(c)
    })
    protected.GET("/brands/:id", func(c echo.Context) error {
        return adminHandlers.BrandHandler.GetBrandByID(c)
    })
    protected.POST("/brands", func(c echo.Context) error {
        return adminHandlers.BrandHandler.CreateBrand(c)
    })
    protected.PUT("/brands/:id", func(c echo.Context) error {
        return adminHandlers.BrandHandler.UpdateBrand(c)
    })
    protected.DELETE("/brands/:id", func(c echo.Context) error {
        return adminHandlers.BrandHandler.DeleteBrand(c)
    })
    protected.PUT("/brands/:id/deactivate", func(c echo.Context) error {
        return adminHandlers.BrandHandler.DeactivateBrand(c)
    })
    protected.PUT("/brands/:id/reactivate", func(c echo.Context) error {
        return adminHandlers.BrandHandler.ReactivateBrand(c)
    })

    // Category routes
    protected.GET("/categories", func(c echo.Context) error {
        return adminHandlers.CategoryHandler.GetAllCategories(c)
    })
    protected.GET("/categories/:id", func(c echo.Context) error {
        return adminHandlers.CategoryHandler.GetCategoryByID(c)
    })
    protected.POST("/categories", func(c echo.Context) error {
        return adminHandlers.CategoryHandler.CreateCategory(c)
    })
    protected.PUT("/categories/:id", func(c echo.Context) error {
        return adminHandlers.CategoryHandler.UpdateCategory(c)
    })
    protected.DELETE("/categories/:id", func(c echo.Context) error {
        return adminHandlers.CategoryHandler.DeleteCategory(c)
    })
    protected.PUT("/categories/:id/deactivate", func(c echo.Context) error {
        return adminHandlers.CategoryHandler.DeactivateCategory(c)
    })
    protected.PUT("/categories/:id/reactivate", func(c echo.Context) error {
        return adminHandlers.CategoryHandler.ReactivateCategory(c)
    })

    // Product routes
    protected.GET("/products", func(c echo.Context) error {
        return adminHandlers.ProductHandler.GetAllProducts(c)
    })
    protected.GET("/products/:id", func(c echo.Context) error {
        return adminHandlers.ProductHandler.GetProductByID(c)
    })
    protected.POST("/products", func(c echo.Context) error {
        return adminHandlers.ProductHandler.CreateProduct(c)
    })
    protected.PUT("/products/:id/details", func(c echo.Context) error {
        return adminHandlers.ProductHandler.UpdateProductDetails(c)
    })
    protected.PUT("/products/:id/image", func(c echo.Context) error {
        return adminHandlers.ProductHandler.UpdateProductImage(c)
    })
    protected.DELETE("/products/:id", func(c echo.Context) error {
        return adminHandlers.ProductHandler.DeleteProduct(c)
    })
    protected.PUT("/products/:id/deactivate", func(c echo.Context) error {
        return adminHandlers.ProductHandler.DeactivateProduct(c)
    })
    protected.PUT("/products/:id/reactivate", func(c echo.Context) error {
        return adminHandlers.ProductHandler.ReactivateProduct(c)
    })
    protected.PUT("/products/:id/reprice", func(c echo.Context) error {
        return adminHandlers.ProductHandler.UpdateProductPrice(c)
    })
    protected.PUT("/products/:id/restock", func(c echo.Context) error {
        return adminHandlers.ProductHandler.UpdateProductStock(c)
    })
}

