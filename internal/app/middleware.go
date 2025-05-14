package app

import (
	"net/http"
	"strings"

	"github.com/Daniel-Njaramba-1/pulse/internal/services/adminSvc"
	"github.com/Daniel-Njaramba-1/pulse/internal/services/customerSvc"
	"github.com/Daniel-Njaramba-1/pulse/internal/util/logging"
	"github.com/labstack/echo/v4"
)

func AdminAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func (c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				logging.LogError("Authorization header is missing")
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				logging.LogError("Invalid authorization format")
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid authorization format"})
			}
			tokenString := parts[1]
			claims, err := adminSvc.VerifyAdminToken(tokenString)
			if err != nil {
				logging.LogError("Invalid token: %v", err)
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			}
			c.Set("username", claims.Username)
			return next(c)
		}
	}
}

func CustomerAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func (c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid authorization format"})
			}
			tokenString := parts[1]
			claims, err := customerSvc.VerifyCustomerToken(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			}
			c.Set("username", claims.Username)
			c.Set("userId", claims.Id)
			return next(c)
		}
	}
}