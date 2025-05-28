package app

import "github.com/labstack/echo/v4"

func CustomerRoutes(e *echo.Echo, customerHandlers *CustomerHdl) {
    customer := e.Group("/api/customer")

    customer.POST("/register", func(c echo.Context) error {
        return customerHandlers.AuthHandler.Register(c)
    })
    customer.POST("/login", func(c echo.Context) error {
        return customerHandlers.AuthHandler.Login(c)
    })
    customer.GET("/products", func(c echo.Context) error{
        return customerHandlers.ProductHandler.GetAllProducts(c)
    })
    customer.GET("/product-by-id/:id", func(c echo.Context) error{
        return customerHandlers.ProductHandler.GetProductByID(c)
    })
    customer.GET("/product-by-name/:name", func(c echo.Context) error {
        return customerHandlers.ProductHandler.GetProductByName(c)
    })

    protected := customer.Group("", CustomerAuthMiddleware())
    
    // cart
    protected.GET("/cart-with-items", func(c echo.Context) error{
        return customerHandlers.CartHandler.GetCartWithItems(c)
    })
    protected.POST("/add-to-cart", func(c echo.Context) error {
        return customerHandlers.CartHandler.AddToCart(c)
    })
    protected.DELETE("/remove-from-cart", func(c echo.Context) error {
        return customerHandlers.CartHandler.RemoveFromCart(c)
    })
    protected.DELETE("/clear-cart", func(c echo.Context) error {
        return customerHandlers.CartHandler.ClearCart(c)
    })

    // order
    protected.POST("/order", func(c echo.Context) error {
        return customerHandlers.OrderHandler.GenerateOrder(c)
    })
    protected.GET("/order-with-items", func(c echo.Context) error {
        return customerHandlers.OrderHandler.GetOrderWithItems(c)
    })

    // payment
    protected.POST("/payment", func(c echo.Context) error {
        return customerHandlers.PaymentHandler.ProcessPayment(c)
    })

    // wishlist
    protected.GET("/wishlist", func(c echo.Context) error {
        return customerHandlers.WishlistHandler.GetWishlistItems(c)
    })

    protected.POST("/add-to-wishlist/:id", func(c echo.Context) error {
        return customerHandlers.WishlistHandler.AddToWishlist(c)
    })

    protected.DELETE("/remove-from-wishlist/:id", func(c echo.Context) error {
        return customerHandlers.WishlistHandler.RemoveFromWishlist(c)
    })

    protected.GET("/check-product-in-wishlist/:id", func(c echo.Context) error {
        return customerHandlers.WishlistHandler.CheckProductInWishlist(c)
    })

    // review
    protected.GET("/product-reviews/:id", func(c echo.Context) error {
        return customerHandlers.ReviewHandler.GetReviewsForProduct(c)
    })
    
    protected.POST("/review-product", func(c echo.Context) error {
        return customerHandlers.ReviewHandler.ReviewProduct(c)
    })

    protected.GET("/verify-purchase/:id", func(c echo.Context) error {
        return customerHandlers.ReviewHandler.VerifyPurchase(c)
    })
}