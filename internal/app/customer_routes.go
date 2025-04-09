package app

import "github.com/labstack/echo/v4"

func CustomerRoutes(e *echo.Echo, customerHandlers *CustomerHdl) {
    customer := e.Group("/api/customer")

    // Customer routes
    customer.POST("/register", func(c echo.Context) error {
        return customerHandlers.AuthHandler.Register(c)
    })
    customer.POST("/login", func(c echo.Context) error {
        return customerHandlers.AuthHandler.Login(c)
    })
}