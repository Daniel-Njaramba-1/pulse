package adminHdl

import (
	"errors"
	"net/http"
	"os"
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
	
	basePriceStr := c.FormValue("base_price")
	basePrice, err := strconv.Atoi(basePriceStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid base price"})
	}
	
	initialStockStr := c.FormValue("initial_stock")
	initialStock, err := strconv.Atoi(initialStockStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid stock quantity"})
	}
	
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

		filename := imageHdl.GenerateFilename(product.Name)
		uploadDir, err := imageHdl.EnsureUploadDirectoryExists()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		imagePath := filepath.Join(uploadDir, filename)
		err = imageHdl.SaveImage(src, imagePath)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error":err.Error()})
		}
		product.ImagePath = &imagePath
	}

	createdProduct, err := h.productService.CreateProduct(c.Request().Context(), &product, basePrice, initialStock)
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

func (h *ProductHandler) UpdateProductDetails(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
	}

	// Use a struct that matches the expected JSON payload from the Svelte store
	type UpdateRequest struct {
		BrandId     int    `json:"brand_id"`
		CategoryId  int    `json:"category_id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		IsActive    bool   `json:"is_active"`
	}

	var req UpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// Map the request to the product model
	var product repo.Product
	product.Id = id
	product.BrandId = req.BrandId
	product.CategoryId = req.CategoryId
	product.Name = req.Name
	product.Description = req.Description
	product.IsActive = req.IsActive

	updatedProduct, err := h.productService.UpdateProductDetails(c.Request().Context(), &product)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update product details", "details": err.Error()})
	}

	return c.JSON(http.StatusOK, updatedProduct)
}

func (h *ProductHandler) UpdateProductImage(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid image file"})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to open image file"})
	}
	defer src.Close()

	productName, err := h.productService.GetProductName(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve product name", "details": err.Error()})
	}
	filename := imageHdl.GenerateFilename(productName)

	uploadPath, err := imageHdl.EnsureUploadDirectoryExists()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	newImagePath := filename
	newFullImagePath := filepath.Join(uploadPath, newImagePath)

	// Find and replace old image
	oldImagePath, err := h.productService.GetProductImagePath(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve old image path", "details": err.Error()})
	}

	err = imageHdl.SaveImage(src, newFullImagePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Delete old image if it exists
	if oldImagePath != "" {
		if _, err := os.Stat(oldImagePath); err == nil {
			// File exists, try to delete it
			err = imageHdl.DeleteImage(oldImagePath)
			if err != nil {
				logging.LogInfo("Warning: Failed to delete old image: %v", err)
				// Continue execution even if deletion fails
			}
		} 
	}

	err = h.productService.UpdateProductImage(c.Request().Context(), id, newImagePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update product image path", "details": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Product image updated successfully"})
}

func (h *ProductHandler) UpdateProductPrice(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
	}

	// Parse request body
	type PriceRequest struct {
		BasePrice float64 `json:"base_price"`
	}

	var req PriceRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// Call service method
	err = h.productService.ChangeBasePrice(c.Request().Context(), id, req.BasePrice)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update product price", "details": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Product price updated successfully"})
}

func (h *ProductHandler) UpdateProductStock(c echo.Context) error {
	logging.LogInfo("UpdateProductStock called")
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	logging.LogInfo("Product ID: %s", idStr)
	if err != nil {
		logging.LogInfo("Error: Invalid product ID")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
	}

	// Parse request body
	type StockRequest struct {
		Quantity int `json:"quantity"`
	}
	

	var req StockRequest
	if err := c.Bind(&req); err != nil {
		logging.LogInfo("Error: Invalid request payload")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}
	logging.LogInfo("Quantity: %d", req.Quantity)

	// Call service method with stock object
	stock := &repo.Stock{
		ProductId: id,
		Quantity:  req.Quantity,
	}
	
	err = h.productService.RestockProduct(c.Request().Context(), stock)
	if err != nil {
		logging.LogInfo("Error: Failed to add product stock")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add product stock", "details": err.Error()})
	}

	logging.LogInfo("Success: Product stock added successfully")	

	return c.JSON(http.StatusOK, map[string]string{"message": "Product stock added successfully"})
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