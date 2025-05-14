package customerSvc

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Daniel-Njaramba-1/pulse/internal/repo"
	"github.com/Daniel-Njaramba-1/pulse/internal/util/logging"
	"github.com/jmoiron/sqlx"
)

type CartService struct {
	db *sqlx.DB
}

func NewCartService(db *sqlx.DB) *CartService {
	return &CartService{db: db}
}

// GetCartWithItems retrieves an active cart with all unprocessed items
func (s *CartService) GetCartWithItems(ctx context.Context, userID int) (*repo.CartWithItems, error) {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // Will be ignored if transaction commits successfully

	// Get the cart first
	cart := repo.Cart{}
	getCartQuery := `
		SELECT id, customer_id, is_active 
		FROM carts
		WHERE customer_id = $1 AND is_active = true
		LIMIT 1
	`
	
	err = tx.QueryRowxContext(ctx, getCartQuery, userID).StructScan(&cart)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("cart not found")
		}
		return nil, err
	}
	
	// Now get the unprocessed items with product details
	items := []repo.CartItemDetail{}
	query := `
		SELECT 
			ci.id, ci.cart_id, ci.product_id, ci.quantity, ci.is_processed, ci.created_at, ci.updated_at,
			p.name as product_name, p.image_path as product_image_path,
			pm.adjusted_price as product_adjusted_price,
			s.quantity as product_stock_quantity
		FROM cart_items ci
		JOIN products p ON ci.product_id = p.id
		LEFT JOIN product_metrics pm ON p.id = pm.product_id
		LEFT JOIN stocks s ON p.id = s.product_id
		WHERE ci.cart_id = $1 AND ci.is_processed = false 
		ORDER BY ci.created_at DESC
	`
	
	err = tx.SelectContext(ctx, &items, query, cart.Id)
	if err != nil {
		return nil, err
	}
	
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	
	// Calculate summary data
	totalItems := 0
	totalPrice := 0.0
	for _, item := range items {
		totalItems += item.Quantity
		totalPrice += float64(item.Quantity) * float64(*item.ProductAdjustedPrice)
	}
	
	// Combine cart with items
	cartWithItems := &repo.CartWithItems{
		Id:         cart.Id,
		CustomerId: cart.CustomerId,
		IsActive:   cart.IsActive,
		Items:      items,
		TotalItems: totalItems,
		TotalPrice: totalPrice,
	}
	
	return cartWithItems, nil
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

func (s *CartService) AddItemToCart(ctx context.Context, userId int, productId int, quantity int) error {
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
		if err == sql.ErrNoRows {
			logging.LogInfo("CartService - Add to Cart - Cart not found, creating new cart")
			
			// Create a new cart for the user
			createCartQuery := `
				INSERT INTO carts (customer_id, is_active, created_at)
				VALUES ($1, TRUE, NOW())
				RETURNING id
			`
			err = tx.QueryRowContext(ctx, createCartQuery, userId).Scan(&cartId)
			if err != nil {
				logging.LogError("CartService - Add to Cart - Failed to create cart: %v", err)
				return err
			}
		} else {
			logging.LogError("CartService - Add to Cart - Error finding cart: %v", err)
			return err
		}
	}

	// check if cart item exists and if processed == FALSE
	var existingItem struct {
		ID         int		`db:"id"`
		Quantity   int		`db:"quantity"`
		IsProcessed bool 	`db:"is_processed"`
	}
	
	checkQuery := `
		SELECT id, quantity, is_processed 
		FROM cart_items
		WHERE cart_id = $1 AND product_id = $2 AND is_processed = FALSE
		LIMIT 1
	`
	
	err = tx.QueryRowxContext(ctx, checkQuery, cartId, productId).StructScan(&existingItem)
	
	if err == nil {
		// Item exists, update quantity
		updateQuery := `
			UPDATE cart_items
			SET quantity = quantity + $1, updated_at = NOW()
			WHERE id = $2
		`
		_, err = tx.ExecContext(ctx, updateQuery, quantity, existingItem.ID)
		if err != nil {
			logging.LogError("CartService - Add to Cart - Error updating existing cart item: %v", err)
			return err
		}
		logging.LogInfo("CartService - Updated existing cart item ID: %d, new quantity: %d", existingItem.ID, existingItem.Quantity + quantity)
	} else if err == sql.ErrNoRows {
		// Item doesn't exist, insert new one
		insertQuery := `
			INSERT INTO cart_items (cart_id, product_id, quantity, is_processed)
			VALUES ($1, $2, $3, FALSE)
		`
		_, err = tx.ExecContext(ctx, insertQuery, cartId, productId, quantity)
		if err != nil {
			logging.LogError("CartService - Add to Cart - Error inserting new cart item: %v", err)
			return err
		}
		logging.LogInfo("CartService - Added new item to cart ID: %d, product ID: %d, quantity: %d", cartId, productId, quantity)
	} else {
		// Some other database error
		logging.LogError("CartService - Add to Cart - Database error checking for existing item: %v", err)
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		logging.LogError("CartService - Add to Cart - Failed to commit transaction: %v", err)
		return err
	}

	return nil
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

	// check if cart item exists and if processed == FALSE
	var item repo.CartItem
	getItemQuery := `
		SELECT * FROM cart_items
		WHERE id = $1 AND cart_id = $2 AND is_processed = FALSE 
	`
	err = tx.GetContext(ctx, &item, getItemQuery, itemId, cartId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("item not found in cart")
		}
		return err
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
