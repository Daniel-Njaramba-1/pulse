package customerSvc

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type ReviewService struct {
	db *sqlx.DB
}

func NewReviewService(db *sqlx.DB) ReviewService {
	return ReviewService{db: db}
}

func (s *ReviewService) ReviewProduct (ctx context.Context, productId int, userId int) (error) {
	// check if user has pruchased product

	// generate review

	return nil
}