// filepath: c:\Users\ADMIN\Desktop\pulse\internal\api\adminHdl\brands.go
package adminHdl

import (
	"net/http"
	"strconv"

	"github.com/Daniel-Njaramba-1/pulse/internal/repo"
	"github.com/Daniel-Njaramba-1/pulse/internal/services/adminSvc"
	"github.com/labstack/echo/v4"
)

type BrandHandler struct {
	brandService *adminSvc.BrandService
}

func NewBrandHandler(brandService *adminSvc.BrandService) *BrandHandler {
	return &BrandHandler{brandService: brandService}
}

// CreateBrand handles the creation of a new brand
func (h *BrandHandler) CreateBrand(c echo.Context) error {
	var brand repo.Brand
	if err := c.Bind(&brand); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
	}

	createdBrand, err := h.brandService.CreateBrand(c.Request().Context(), &brand)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, createdBrand)
}

// GetBrandByID handles retrieving a brand by its ID
func (h *BrandHandler) GetBrandByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid brand ID"})
	}

	brand, err := h.brandService.GetBrandByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, brand)
}

// GetAllBrands handles retrieving all brands
func (h *BrandHandler) GetAllBrands(c echo.Context) error {
	brands, err := h.brandService.GetAllBrands(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, brands)
}

// UpdateBrand handles updating an existing brand
func (h *BrandHandler) UpdateBrand(c echo.Context) error {
	var brand repo.Brand
	if err := c.Bind(&brand); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
	}

	updatedBrand, err := h.brandService.UpdateBrand(c.Request().Context(), &brand)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, updatedBrand)
}

// DeactivateBrand handles deactivating a brand
func (h *BrandHandler) DeactivateBrand(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid brand ID"})
	}

	if err := h.brandService.DeactivateBrand(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "brand deactivated successfully"})
}

// ReactivateBrand handles reactivating a brand
func (h *BrandHandler) ReactivateBrand(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid brand ID"})
	}

	if err := h.brandService.ReactivateBrand(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "brand reactivated successfully"})
}

// DeleteBrand handles deleting a brand
func (h *BrandHandler) DeleteBrand(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid brand ID"})
	}

	if err := h.brandService.DeleteBrand(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "brand deleted successfully"})
}