package customerHdl

import (
	"net/http"

	"github.com/Daniel-Njaramba-1/pulse/internal/services/customerSvc"
	"github.com/labstack/echo/v4"
)

type PaymentHandler struct {
	paymentService *customerSvc.PaymentService
}

func NewPaymentHandler(paymentService *customerSvc.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

func (h *PaymentHandler) ProcessPayment(c echo.Context) error {
	// Get the user ID from the context (assuming it's set during authentication)
	userId := c.Get("userId").(int)

	// Call the service to process the payment
	payment, err := h.paymentService.ProcessPayment(c.Request().Context(), userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, payment)
}