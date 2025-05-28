package customerHdl

import (
	"net/http"
	"strconv"

	"github.com/Daniel-Njaramba-1/pulse/internal/services/customerSvc"
	"github.com/labstack/echo/v4"
)

type WishlistHandler struct {
	wishlistService *customerSvc.WishlistService
}

func NewWishlistHandler(wishlistService *customerSvc.WishlistService) *WishlistHandler {
	return &WishlistHandler{
		wishlistService: wishlistService,
	}
}

func (h *WishlistHandler) AddToWishlist(c echo.Context) error {
	userId, ok := c.Get("userId").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user id"})
	}

	productIdStr := c.Param("id")
	productId, err := strconv.Atoi(productIdStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid product id"})
	}

	err = h.wishlistService.AddToWishlist(c.Request().Context(), userId, productId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "item added to wishlist"})
}

func (h *WishlistHandler) RemoveFromWishlist(c echo.Context) error {
	userId, ok := c.Get("userId").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user id"})
	}

	productIdStr := c.Param("id")
	productId, err := strconv.Atoi(productIdStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid product id"})
	}

	err = h.wishlistService.RemoveFromWishlist(c.Request().Context(), userId, productId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "item removed from wishlist"})
}

func (h *WishlistHandler) GetWishlistItems(c echo.Context) error {
	userId, ok := c.Get("userId").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user id"})
	}

	items, err := h.wishlistService.GetWishlistDetail(c.Request().Context(), userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, items)
}

func (h *WishlistHandler) CheckProductInWishlist(c echo.Context) error {
	userId, ok := c.Get("userId").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user id"})
	}

	productIdStr := c.Param("id")
	productId, err := strconv.Atoi(productIdStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid product id"})
	}

	exists, err := h.wishlistService.CheckProductInWishlist(c.Request().Context(), userId, productId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]bool{"exists": exists})
}