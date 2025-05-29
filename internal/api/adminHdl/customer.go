package adminHdl

import (
	"net/http"

	"github.com/Daniel-Njaramba-1/pulse/internal/services/adminSvc"
	"github.com/labstack/echo/v4"
)

type CustomerHandler struct {
	customerService adminSvc.CustomerService
}

func NewCustomerHandler(customerService adminSvc.CustomerService) *CustomerHandler {
	return &CustomerHandler{
		customerService: customerService,
	}
}

func (h *CustomerHandler) GetAllCustomers(c echo.Context) error {
	customers, err := h.customerService.GetAllCustomers(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch customers",
		})
	}

	return c.JSON(http.StatusOK, customers)
}
