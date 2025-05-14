package customerHdl

import "github.com/Daniel-Njaramba-1/pulse/internal/services/customerSvc"

type ReviewHandler struct {
	reviewService *customerSvc.ReviewService
}

func NewReviewHandler(reviewService *customerSvc.ReviewService) *ReviewHandler {
	return &ReviewHandler{
		reviewService: reviewService,
	}
}