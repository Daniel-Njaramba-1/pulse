package customerHdl

import (
	"net/http"
	"strconv"

	"github.com/Daniel-Njaramba-1/pulse/internal/services/customerSvc"
	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	productService *customerSvc.ProductService
}

func NewProductHandler(productService *customerSvc.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) GetAllProducts(c echo.Context) error {
	// Call the service to get the products
	products, err := h.productService.GetAllProducts(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) GetProductByID(c echo.Context) error {	
	// Get the product ID from the URL parameters
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid brand ID"})
	}

	// Call the service to get the product by ID
	product, err := h.productService.GetProductByID(c.Request().Context(), productID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) GetProductByName(c echo.Context) error {
	// Get the product name from the URL parameters
	productName := c.Param("name")

	// Call the service to get the product by name
	product, err := h.productService.GetProductByName(c.Request().Context(), productName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, product)
}