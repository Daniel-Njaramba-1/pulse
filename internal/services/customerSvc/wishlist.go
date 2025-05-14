package customerSvc

import (
	"context"

	"github.com/Daniel-Njaramba-1/pulse/internal/repo"
	"github.com/jmoiron/sqlx"
)

type WishlistService struct {
	db *sqlx.DB
}

func NewWishlistService(db *sqlx.DB) *WishlistService {
	return &WishlistService{db: db}
}

func (s *WishlistService) AddToWishlist(ctx context.Context, customerId int, productId int) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var wishlistId int
	getWishlistQuery := `
		SELECT id
		FROM wishlists
		WHERE customer_id = $1 AND is_active = TRUE
		LIMIT 1
	`
	err = tx.QueryRowxContext(ctx, getWishlistQuery, customerId).Scan(&wishlistId)
	if err != nil {
		// Create wishlist if not found
		createQuery := `
			INSERT INTO wishlists (customer_id, is_active)
			VALUES ($1, TRUE)
			RETURNING id
		`
		err = tx.QueryRowxContext(ctx, createQuery, customerId).Scan(&wishlistId)
		if err != nil {
			return err
		}
	}

	// Check if product already exists
	var exists bool
	checkQuery := `
		SELECT EXISTS (
			SELECT 1 FROM wishlist_items
			WHERE wishlist_id = $1 AND product_id = $2
		)
	`
	err = tx.QueryRowxContext(ctx, checkQuery, wishlistId, productId).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	// Insert item
	insertQuery := `
		INSERT INTO wishlist_items (wishlist_id, product_id)
		VALUES ($1, $2)
	`
	_, err = tx.ExecContext(ctx, insertQuery, wishlistId, productId)
	if err != nil {
		return err
	}

	// Update product_metrics.wishlist_count
	_, err = tx.ExecContext(ctx, `
		UPDATE product_metrics
		SET wishlist_count = wishlist_count + 1
		WHERE product_id = $1
	`, productId)
	if err != nil {
		return err
	}

	return tx.Commit()
}


func (s *WishlistService) RemoveFromWishlist(ctx context.Context, customerId int, productId int) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var wishlistId int
	getWishlistQuery := `
		SELECT id FROM wishlists
		WHERE customer_id = $1 AND is_active = TRUE
		LIMIT 1
	`
	err = tx.QueryRowxContext(ctx, getWishlistQuery, customerId).Scan(&wishlistId)
	if err != nil {
		return err
	}

	// Delete item
	deleteQuery := `
		DELETE FROM wishlist_items
		WHERE wishlist_id = $1 AND product_id = $2
	`
	res, err := tx.ExecContext(ctx, deleteQuery, wishlistId, productId)
	if err != nil {
		return err
	}

	affected, _ := res.RowsAffected()
	if affected > 0 {
		// Decrement product_metrics.wishlist_count
		_, err = tx.ExecContext(ctx, `
			UPDATE product_metrics
			SET wishlist_count = GREATEST(0, wishlist_count - 1)
			WHERE product_id = $1
		`, productId)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}


func (s *WishlistService) GetWishlistItems(ctx context.Context, customerId int) ([]repo.WishlistItem, error) {
	var items []repo.WishlistItem

	query := `
		SELECT wi.*
		FROM wishlist_items wi
		JOIN wishlists w ON wi.wishlist_id = w.id
		WHERE w.customer_id = $1 AND w.is_active = TRUE
		ORDER BY wi.created_at DESC
	`

	err := s.db.SelectContext(ctx, &items, query, customerId)
	if err != nil {
		return nil, err
	}

	return items, nil
}
