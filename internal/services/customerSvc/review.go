package customerSvc

import (
	"context"
	"fmt"

	"github.com/Daniel-Njaramba-1/pulse/internal/repo"
	"github.com/jmoiron/sqlx"
)

type ReviewService struct {
	db *sqlx.DB
}

func NewReviewService(db *sqlx.DB) *ReviewService {
	return &ReviewService{db: db}
}

func (s *ReviewService) GetReviewforProduct(ctx context.Context, productId int) ([]*repo.ReviewDetail, error) {
	query := `
		SELECT r.id, r.customer_id, r.product_id, r.rating, r.review_text, r.created_at, c.username
		FROM reviews r
		INNER JOIN customers c ON r.customer_id = c.id
		WHERE r.product_id = $1
		ORDER BY r.id DESC
	`
	var reviews []*repo.ReviewDetail
	err := s.db.SelectContext(ctx, &reviews, query, productId)
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

func (s *ReviewService) VerifyPurchase(ctx context.Context, userId int, productId int) (bool, error) {
	var count int
	query := `
		SELECT COUNT(*)
		FROM sales
		WHERE product_id = $1 AND order_item_id IN (
			SELECT id FROM order_items WHERE order_id IN (
				SELECT id FROM orders WHERE customer_id = $2 AND status = 'completed'
			)
		)
	`
	err := s.db.QueryRowContext(ctx, query, productId, userId).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}


func (s *ReviewService) ReviewProduct(ctx context.Context, review *repo.Review) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Verify user has purchased the product
	var count int
	query := `
		SELECT COUNT(*)
		FROM sales
		WHERE product_id = $1 AND order_item_id IN (
			SELECT id FROM order_items WHERE order_id IN (
				SELECT id FROM orders WHERE customer_id = $2 AND status = 'completed'
			)
		)
	`
	err = tx.QueryRowContext(ctx, query, review.ProductId, review.CustomerId).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("user has not purchased this product")
	}

	// 2. Upsert review (insert or update)
	var exists bool
	checkQuery := `
		SELECT EXISTS (
			SELECT 1 FROM reviews WHERE customer_id = $1 AND product_id = $2
		)
	`
	err = tx.QueryRowContext(ctx, checkQuery, review.CustomerId, review.ProductId).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		updateQuery := `
			UPDATE reviews
			SET rating = $1, review_text = $2
			WHERE customer_id = $3 AND product_id = $4
		`
		_, err = tx.ExecContext(ctx, updateQuery,
			review.Rating, review.ReviewText, review.CustomerId, review.ProductId)
		if err != nil {
			return err
		}
	} else {
		insertQuery := `
			INSERT INTO reviews (customer_id, product_id, rating, review_text)
			VALUES ($1, $2, $3, $4)
		`
		_, err = tx.ExecContext(ctx, insertQuery,
			review.CustomerId, review.ProductId, review.Rating, review.ReviewText)
		if err != nil {
			return err
		}
	}

	// 3. Update product_metrics (average_rating and review_count)
	var avgRating float64
	var totalReviews int
	aggQuery := `
		SELECT COALESCE(AVG(rating), 0), COUNT(*)
		FROM reviews
		WHERE product_id = $1
	`
	err = tx.QueryRowContext(ctx, aggQuery, review.ProductId).Scan(&avgRating, &totalReviews)
	if err != nil {
		return err
	}

	updateMetrics := `
		UPDATE product_metrics
		SET average_rating = $1, review_count = $2, updated_at = NOW()
		WHERE product_id = $3
	`
	_, err = tx.ExecContext(ctx, updateMetrics, avgRating, totalReviews, review.ProductId)
	if err != nil {
		return err
	}

	return tx.Commit()
}
