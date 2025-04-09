package customerSvc

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Daniel-Njaramba-1/pulse/internal/repo"
	"github.com/jmoiron/sqlx"
)

type CartService struct {
	db *sqlx.DB
}

func NewCartService(db *sqlx.DB) *CartService {
	return &CartService{db: db}
}

func (s *CartService) GetCartByUserID(ctx context.Context, userId int) (*repo.Cart, error) {
	var cart repo.Cart
	query := `
		SELECT *
		FROM carts
		WHERE customer_id = $1 AND is_active = TRUE
		LIMIT 1
	`
	err := s.db.GetContext(ctx, &cart, query, userId)
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (s *CartService) AddItemToCart(ctx context.Context, userId int, item *repo.CartItem) error {
	// begin transaction
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err 
	}
	defer tx.Rollback()

	// get cart
	var cartId int
	getCartQuery := `
		SELECT id
		FROM carts
		WHERE customer_id = $1 AND is_active = TRUE
		LIMIT 1
	`
	err = tx.QueryRowContext(ctx, getCartQuery, userId).Scan(&cartId)
	if err != nil {
		return err
	}

	// check if cart item exists and if processed == FALSE
	var existingItem repo.CartItem
	checkQuery := `
		SELECT * FROM cart_items
		WHERE cart_id = $1 AND product_id = $2 AND is_processed = FALSE
		LIMIT 1
	`
	err = tx.GetContext(ctx, &existingItem, checkQuery, item.CartId, item.ProductId)
	if err == nil {
		updateQuery := `
			UPDATE cart_items
			SET quantity = quantity + $1
			WHERE cart_id = $2 AND product_id = $3
		`
		_, err = tx.ExecContext(ctx, updateQuery, item.Quantity, item.CartId, item.ProductId)
		if err != nil {
			return err
		}
	} else if err == sql.ErrNoRows {
		insertQuery := `
			INSERT INTO cart_items (cart_id, product_id, quantity, is_processed)
			VALUES ($1, $2, $3, false)
		`
		_, err = tx.ExecContext(ctx, insertQuery, item.CartId, item.ProductId, item.Quantity)
		if err != nil {
			return err
		}
	} else {
		return err
	}

	return tx.Commit()
}

func (s *CartService) RemoveItemFromCart(ctx context.Context, userId int, itemId int) error {
	// begin transaction
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err 
	}
	defer tx.Rollback()

	// get cart
	var cartId int
	getCartQuery := `
		SELECT id
		FROM carts
		WHERE customer_id = $1 AND is_active = TRUE
		LIMIT 1
	`
	err = tx.QueryRowContext(ctx, getCartQuery, userId).Scan(&cartId)
	if err != nil {
		return err
	}

	// check if cart item exists and if processed == TRUE
	var item repo.CartItem
	getItemQuery := `
		SELECT * FROM cart_items
		WHERE id = $1 AND cart_id = $2 
	`
	err = tx.GetContext(ctx, &item, getItemQuery, itemId, item.CartId)
	if err == nil {
		return err
	}

	if item.IsProcessed {
		return errors.New ("cannot remove an item that has been processed")
	}

	deleteQuery := `
		DELETE FROM cart_items 
		WHERE id = $1
	`
	_, err = tx.ExecContext(ctx, deleteQuery, itemId)
	if err != nil {
		return err
	}

	return tx.Commit()
}


// UpdateCartItemQuantity updates the quantity of an item in the cart
func (s *CartService) UpdateCartItemQuantity(ctx context.Context, userID int, itemID int, newQuantity int) error {
	if newQuantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}
	
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	// Get the customer's cart ID
	var cartID int
	getCartQuery := `
		SELECT id FROM carts
		WHERE customer_id = $1 AND is_active = true
		LIMIT 1
	`
	
	err = tx.QueryRowxContext(ctx, getCartQuery, userID).Scan(&cartID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("no active cart found for user")
		}
		return err
	}
	
	// Get the item first to check if it's processed and belongs to user's cart
	var item repo.CartItem
	getItemQuery := `
		SELECT * FROM cart_items 
		WHERE id = $1 AND cart_id = $2
	`
	
	err = tx.GetContext(ctx, &item, getItemQuery, itemID, cartID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("cart item not found")
		}
		return err
	}
	
	// Check if item has been processed
	if item.IsProcessed {
		return errors.New("cannot update processed cart item")
	}
	
	// Update the quantity
	updateQuery := `UPDATE cart_items SET quantity = $1 WHERE id = $3`
	_, err = tx.ExecContext(ctx, updateQuery, newQuantity, itemID)
	if err != nil {
		return err
	}
	
	return tx.Commit()
}

// GetCartItems retrieves all items in a customer's cart
func (s *CartService) GetCartItems(ctx context.Context, userID int) ([]repo.CartItem, error) {
	// Get the cart ID first
	var cartID int
	getCartQuery := `
		SELECT id FROM carts
		WHERE customer_id = $1 AND is_active = true
		LIMIT 1
	`
	
	err := s.db.QueryRowxContext(ctx, getCartQuery, userID).Scan(&cartID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No cart means no items
			return []repo.CartItem{}, nil
		}
		return nil, err
	}
	
	// Now get the items
	var items []repo.CartItem
	query := `
		SELECT * FROM cart_items 
		WHERE cart_id = $1 AND is_processed = false 
		ORDER BY created_at DESC
	`
	
	err = s.db.SelectContext(ctx, &items, query, cartID)
	if err != nil {
		return nil, err
	}
	
	return items, nil
}

// ClearCart removes all unprocessed items from a customer's cart
func (s *CartService) ClearCart(ctx context.Context, userID int) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	// Get the cart ID first
	var cartID int
	getCartQuery := `
		SELECT id FROM carts
		WHERE customer_id = $1 AND is_active = true
		LIMIT 1
	`
	
	err = tx.QueryRowxContext(ctx, getCartQuery, userID).Scan(&cartID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No cart to clear
			return nil
		}
		return err
	}
	
	// Delete all unprocessed items
	deleteQuery := `DELETE FROM cart_items WHERE cart_id = $1 AND is_processed = false`
	_, err = tx.ExecContext(ctx, deleteQuery, cartID)
	if err != nil {
		return err
	}
	
	return tx.Commit()
}

// MarkCartItemsAsProcessed marks all items in a customer's cart as processed (e.g., after checkout)
func (s *CartService) MarkCartItemsAsProcessed(ctx context.Context, userID int) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	// Get the cart ID first
	var cartID int
	getCartQuery := `
		SELECT id FROM carts
		WHERE customer_id = $1 AND is_active = true
		LIMIT 1
	`
	
	err = tx.QueryRowxContext(ctx, getCartQuery, userID).Scan(&cartID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No cart to process
			return nil
		}
		return err
	}
	
	// Mark all unprocessed items as processed
	updateQuery := `
		UPDATE cart_items 
		SET is_processed = true
		WHERE cart_id = $2 AND is_processed = false
	`
	_, err = tx.ExecContext(ctx, updateQuery, cartID)
	if err != nil {
		return err
	}
	
	return tx.Commit()
}

// GetCartSummary gets count of items and total items in cart
func (s *CartService) GetCartSummary(ctx context.Context, userID int) (int, int, error) {
	// Get the cart ID first
	var cartID int
	getCartQuery := `
		SELECT id FROM carts
		WHERE customer_id = $1 AND is_active = true
		LIMIT 1
	`
	
	err := s.db.QueryRowxContext(ctx, getCartQuery, userID).Scan(&cartID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No cart means empty summary
			return 0, 0, nil
		}
		return 0, 0, err
	}
	
	// Get summary info using a single query
	var itemCount int
	var totalQuantity int
	summaryQuery := `
		SELECT 
			COUNT(id) as item_count,
			COALESCE(SUM(quantity), 0) as total_quantity
		FROM cart_items
		WHERE cart_id = $1 AND is_processed = false
	`
	
	err = s.db.QueryRowxContext(ctx, summaryQuery, cartID).Scan(&itemCount, &totalQuantity)
	if err != nil {
		return 0, 0, err
	}
	
	return itemCount, totalQuantity, nil
}