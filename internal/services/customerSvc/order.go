package customerSvc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Daniel-Njaramba-1/pulse/internal/repo"
	"github.com/jmoiron/sqlx"
)

type OrderService struct {
	db *sqlx.DB
}

func NewOrderService(db *sqlx.DB) *OrderService {
	return &OrderService{ db: db}
}

func (s *OrderService) GenerateOrder(ctx context.Context, userId int) (*repo.Order, error) {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var cartId int
	getCartQuery := `
		SELECT id
		FROM carts
		WHERE customer_id = $1 AND is_active = TRUE
		LIMIT 1
	`
	err = tx.QueryRowxContext(ctx, getCartQuery, userId).Scan(&cartId)
	if err != nil {
		return nil, err
	}

	var cartItems []repo.CartItem
	getCartItemsQuery := `
		SELECT *
		FROM cart_items
		WHERE cart_id = $1 AND is_processed = FALSE
	`
	err = tx.SelectContext(ctx, &cartItems, getCartItemsQuery, cartId)
	if err != nil {
		return nil, err 
	}

	if len (cartItems) == 0 {
		return nil, errors.New("no items in cart")
	}

	// calc price
	var totalPrice float64
	var orderItems []repo.OrderItem

	for _, item := range cartItems {
		var currentPrice float64
		// locking price rows to prevent race condition using SELECT FOR UPDATE
		getPriceQuery := `
			SELECT adjusted_price
			FROM product_metrics
			WHERE product_id = $1
			FOR UPDATE
		`

		err  = tx.QueryRowxContext(ctx, getPriceQuery, item.ProductId).Scan(&currentPrice)
		if err != nil {
			return nil, err 
		}

		itemTotal := currentPrice * float64(item.Quantity)
		totalPrice += itemTotal

		orderItems = append(orderItems, repo.OrderItem{
			ProductId: item.ProductId,
			Price: currentPrice,
			Quantity: item.Quantity,
		})
	}

	// Create the order with expiration time for price validity
	expirationTime := time.Now().Add(30 * time.Minute)
	order := &repo.Order{
		CustomerId:     userId,
		TotalPrice:     totalPrice,
		Status:         "pending",
		PriceValidUntil: expirationTime,
	}
	
	// Insert the order
	insertOrderQuery := `
		INSERT INTO orders (customer_id, total_price, status, price_valid_until)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	err = tx.QueryRowxContext(
		ctx, 
		insertOrderQuery, 
		order.CustomerId, 
		order.TotalPrice, 
		order.Status, 
		order.PriceValidUntil,
	).Scan(&order.Id)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	
	// Insert order items
	insertOrderItemQuery := `
		INSERT INTO order_items (order_id, product_id, price, quantity)
		VALUES ($1, $2, $3, $4)
	`
	
	for _, item := range orderItems {
		_, err = tx.ExecContext(
			ctx,
			insertOrderItemQuery,
			order.Id,
			item.ProductId,
			item.Price,
			item.Quantity,
		)
		
		if err != nil {
			return nil, fmt.Errorf("failed to create order item: %w", err)
		}
	}
	
	// Mark cart items as processed
	updateCartItemsQuery := `
		UPDATE cart_items
		SET is_processed = TRUE
		WHERE cart_id = $1
	`
	_, err = tx.ExecContext(ctx, updateCartItemsQuery, cartId)
	if err != nil {
		return nil, fmt.Errorf("failed to update cart items: %w", err)
	}
	
	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	return order, nil
}

