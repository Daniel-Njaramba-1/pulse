package adminHdl

import (
	"errors"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/Daniel-Njaramba-1/pulse/internal/repo"
	"github.com/Daniel-Njaramba-1/pulse/internal/services/adminSvc"
	"github.com/Daniel-Njaramba-1/pulse/internal/util/imageHdl"
	"github.com/Daniel-Njaramba-1/pulse/internal/util/logging"
	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	productService *adminSvc.ProductService
}

func NewProductHandler(productService *adminSvc.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) CreateProduct(c echo.Context) error {
	var product repo.Product
    
    // Get form values directly
    categoryIdStr := c.FormValue("category_id")
    brandIdStr := c.FormValue("brand_id")
    
    // Convert to int
    categoryId, err := strconv.Atoi(categoryIdStr)
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid category ID"})
    }
    
    brandId, err := strconv.Atoi(brandIdStr)
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid brand ID"})
    }
    
    // Assign values manually
    product.CategoryId = categoryId
    product.BrandId = brandId
    product.Name = c.FormValue("name")
    product.Description = c.FormValue("description")
    product.IsActive = true  // Set default value for IsActive
	
	file, err := c.FormFile("image")
	if err != nil { 
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid image file"})
	}

	if file != nil {
		src, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to open image file"})
		}
		defer src.Close()

		filename := imageHdl.GenerateFilename(file.Filename)
		uploadDir, err := imageHdl.EnsureUploadDirectoryExists()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		imagePath := filepath.Join(uploadDir, filename)
		err = imageHdl.SaveImage(src, imagePath)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error":err.Error()})
		}
		product.ImagePath = imagePath
	}

	createdProduct, err := h.productService.CreateProduct(c.Request().Context(), &product)
	if err != nil {
		logging.LogInfo("Service function called, error: could not create product")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create product", "details": err.Error()})
	}

	return c.JSON(http.StatusCreated, createdProduct)
}

func (h *ProductHandler) GetProductByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
	}

	product, err := h.productService.GetProductByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, errors.New("product not found")) { // Replace with your actual not found error
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Product not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve product", "details": err.Error()})
	}

	return c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) GetAllProducts(c echo.Context) error {
	products, err := h.productService.GetAllProducts(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve products", "details": err.Error()})
	}
	return c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) UpdateProduct(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
	}

	var product repo.Product
	err = c.Bind(&product)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}
	product.Id = id

	file, err := c.FormFile("image")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid image file"})
	}

	if file != nil {
		src, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to open image file"})
		}
		defer src.Close()

		filename := imageHdl.GenerateFilename(file.Filename)
		uploadDir, err := imageHdl.EnsureUploadDirectoryExists()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		newImagePath := filepath.Join(uploadDir, filename)

		err = imageHdl.SaveImage(src, newImagePath)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		product.ImagePath = newImagePath
	}

	updatedProduct, err := h.productService.UpdateProduct(c.Request().Context(), &product)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update product", "details": err.Error()})
	}

	return c.JSON(http.StatusOK, updatedProduct)
}

func (h *ProductHandler) DeactivateProduct(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
	}

	err = h.productService.DeactivateProduct(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to deactivate product", "details": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Product deactivated successfully"})
}

func (h *ProductHandler) ReactivateProduct(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
	}

	err = h.productService.ReactivateProduct(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to reactivate product", "details": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Product reactivated successfully"})
}

func (h *ProductHandler) DeleteProduct(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
	}

	err = h.productService.DeleteProduct(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete product", "details": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Product deleted successfully"})
}