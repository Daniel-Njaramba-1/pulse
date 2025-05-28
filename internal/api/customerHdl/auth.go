package customerHdl

import (
	"net/http"

	"github.com/Daniel-Njaramba-1/pulse/internal/repo"
	"github.com/Daniel-Njaramba-1/pulse/internal/services/customerSvc"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authentication *customerSvc.Authentication
}

func NewAuthHandler(authentication *customerSvc.Authentication) *AuthHandler {
	return &AuthHandler{authentication: authentication}
}

// Register handles customer registration
func (h *AuthHandler) Register(c echo.Context) error {
    var customer repo.Customer
    if err := c.Bind(&customer); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
    } 

    token, user, err := h.authentication.RegisterCustomer(c.Request().Context(), &customer)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "token": token,
		"user": map[string]interface{}{
			"id":       user.Id,
			"username": user.Username,
			"email":    user.Email,
		},
    })
}

// Login handles customer login 
func (h *AuthHandler) Login(c echo.Context) error {
    var credentials struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    if err := c.Bind(&credentials); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
    }

    token, user, err := h.authentication.LoginCustomer(c.Request().Context(), credentials.Username, credentials.Password)
    if err != nil {
        return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "token": token,
		"user": map[string]interface{}{
			"id":       user.Id,
			"username": user.Username,
			"email":    user.Email,
		},
    })
}