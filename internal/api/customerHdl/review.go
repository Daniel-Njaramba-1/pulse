package customerHdl

import (
	"net/http"
	"strconv"

	"github.com/Daniel-Njaramba-1/pulse/internal/repo"
	"github.com/Daniel-Njaramba-1/pulse/internal/services/customerSvc"
	"github.com/labstack/echo/v4"
)

type ReviewHandler struct {
	reviewService *customerSvc.ReviewService
}

func NewReviewHandler(reviewService *customerSvc.ReviewService) *ReviewHandler {
	return &ReviewHandler{
		reviewService: reviewService,
	}
}

// GetReviewsForProduct handles GET /reviews?product_id= endpoint
func (h *ReviewHandler) GetReviewsForProduct(c echo.Context) error {
	productIdStr := c.Param("id")
	productId, err := strconv.Atoi(productIdStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
	}

	reviews, err := h.reviewService.GetReviewforProduct(c.Request().Context(), productId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, reviews)
}

// VerifyPurchase handles GET /verify-purchase endpoint
func (h *ReviewHandler) VerifyPurchase(c echo.Context) error {
	userId, ok := c.Get("userId").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	productIdStr := c.Param("id")
	productId, err := strconv.Atoi(productIdStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
	}

	isVerified, err := h.reviewService.VerifyPurchase(c.Request().Context(), userId, productId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]bool{"purchased": isVerified})
}

// ReviewProduct handles POST /review endpoint
func (h *ReviewHandler) ReviewProduct(c echo.Context) error {
	var review repo.Review

	userId, ok := c.Get("userId").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	type ReviewRequest struct {
		ProductID  int    `json:"product_id"`
		Rating     int    `json:"rating"`
		ReviewText string `json:"review_text"`
	}

	var req ReviewRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	productIdStr := strconv.Itoa(req.ProductID)
	ratingStr := strconv.Itoa(req.Rating)
	reviewText := req.ReviewText

	productId, err := strconv.Atoi(productIdStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
	}
	rating, err := strconv.Atoi(ratingStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid rating"})
	}

	review.CustomerId = userId
	review.ProductId = productId
	review.Rating = float32(rating)
	review.ReviewText = reviewText

	err = h.reviewService.ReviewProduct(c.Request().Context(), &review)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Review submitted successfully"})
}