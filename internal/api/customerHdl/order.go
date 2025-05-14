package customerHdl

import (
	"net/http"

	"github.com/Daniel-Njaramba-1/pulse/internal/services/customerSvc"
	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	orderService *customerSvc.OrderService 
}

func NewOrderHandler(orderService *customerSvc.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) GenerateOrder(c echo.Context) error {
	// Get the user ID from the context (assuming it's set during authentication)
	userId := c.Get("userId").(int)

	// Call the service to create the order
	err := h.orderService.GenerateOrder(c.Request().Context(), userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Order generated successfully"})
}

func (h *OrderHandler) GetOrderWithItems (c echo.Context) error {
	userId := c.Get("userId").(int)

	orderWithItems, err := h.orderService.GetOrderWithItems(c.Request().Context(), userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, orderWithItems)
}

