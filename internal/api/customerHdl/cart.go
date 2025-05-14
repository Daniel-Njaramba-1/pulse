package customerHdl

import (
	"net/http"

	"github.com/Daniel-Njaramba-1/pulse/internal/services/customerSvc"
	"github.com/Daniel-Njaramba-1/pulse/internal/util/logging"
	"github.com/labstack/echo/v4"
)

type CartHandler struct {
	cartService *customerSvc.CartService
}

func NewCartHandler(cartService *customerSvc.CartService) *CartHandler {
	return &CartHandler{cartService: cartService}
}

func (h *CartHandler) GetCart(c echo.Context) error {
	// Get the user ID from the context (assuming it's set during authentication)
	userId := c.Get("userId").(int)

	// Call the service to get the cart
	cart, err := h.cartService.GetCartByUserID(c.Request().Context(), userId)
	if err != nil {
		// Specific error handling based on error type
		if err.Error() == "cart not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "No active cart found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, cart)
}

func (h *CartHandler) GetCartWithItems(c echo.Context) error {
	// Get the user ID from the context
	userId := c.Get("userId").(int)

	// Call the service to get the cart items
	cartWithItems, err := h.cartService.GetCartWithItems(c.Request().Context(), userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, cartWithItems)
}

func (h *CartHandler) AddToCart(c echo.Context) error {
	// Get the user ID from the context
	userId := c.Get("userId").(int)

	// Define a struct that matches the expected JSON format
	var requestBody struct {
		ProductId int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}

	// Parse the JSON request body
	if err := c.Bind(&requestBody); err != nil {
		logging.LogInfo("CartHandler - Add to Cart - Invalid Request Body: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Call the service to add the item to the cart
	err := h.cartService.AddItemToCart(c.Request().Context(), userId, requestBody.ProductId, requestBody.Quantity)
	if err != nil {
		logging.LogInfo("CartHandler - Add to Cart - Error Adding to Cart")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Item added to cart successfully"})
}

func (h *CartHandler) RemoveFromCart(c echo.Context) error {
	// Get the user ID from the context
	userId := c.Get("userId").(int)

	// Parse the request body to get the item ID
	var request struct {
		ItemId int `json:"itemId"`
	}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Call the service to remove the item from the cart
	err := h.cartService.RemoveItemFromCart(c.Request().Context(), userId, request.ItemId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Item removed from cart successfully"})
}

func (h *CartHandler) ClearCart(c echo.Context) error {
	// Get the user ID from the context
	userId := c.Get("userId").(int)

	// Call the service to clear the cart
	err := h.cartService.ClearCart(c.Request().Context(), userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Cart cleared successfully"})
}
