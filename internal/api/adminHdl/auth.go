package adminHdl

import (
	"net/http"

	"github.com/Daniel-Njaramba-1/pulse/internal/repo"
	"github.com/Daniel-Njaramba-1/pulse/internal/services/adminSvc"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
    authentication *adminSvc.Authentication
}

func NewAuthHandler(authentication *adminSvc.Authentication) *AuthHandler {
    return &AuthHandler{authentication: authentication}
}

// Register handles admin registration
func (h *AuthHandler) Register(c echo.Context) error {
    var admin repo.Admin
    if err := c.Bind(&admin); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
    }

    token, user, err := h.authentication.RegisterAdmin(c.Request().Context(), &admin)
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

// Login handles admin login
func (h *AuthHandler) Login(c echo.Context) error {
    var credentials struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    if err := c.Bind(&credentials); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
    }

    token, user, err := h.authentication.LoginAdmin(c.Request().Context(), credentials.Username, credentials.Password)
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