package customerSvc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Daniel-Njaramba-1/pulse/internal/repo"
	"github.com/Daniel-Njaramba-1/pulse/internal/util/logging"
	"github.com/jmoiron/sqlx"
)

type OrderService struct {
	db *sqlx.DB
}

func NewOrderService(db *sqlx.DB) *OrderService {
	return &OrderService{ db: db}
}

func (s *OrderService) GenerateOrder(ctx context.Context, userId int) error {
	logging.LogInfo("Order generation started")
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		logging.LogError(fmt.Sprintf("Failed to begin transaction: %v", err))
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
	err = tx.QueryRowxContext(ctx, getCartQuery, userId).Scan(&cartId)
	if err != nil {
		logging.LogError(fmt.Sprintf("Failed to get cart for user %d: %v", userId, err))
		return err
	}
	logging.LogInfo(fmt.Sprintf("Cart ID %d found for user %d", cartId, userId))

	// check for pending order, if pending, return error, if not create order
	var existingOrderId int
	getExistingOrderQuery := `
		SELECT id
		FROM orders
		WHERE customer_id = $1 AND status = 'pending'
		LIMIT 1
	`
	err = tx.QueryRowxContext(ctx, getExistingOrderQuery, userId).Scan(&existingOrderId)

	if err == nil {
		logging.LogInfo(fmt.Sprintf("User %d already has a pending order (ID: %d)", userId, existingOrderId))
		return errors.New("there is an existing pending order already")
	} else if !errors.Is(err, sql.ErrNoRows) {
		logging.LogError(fmt.Sprintf("Error checking for existing order for user %d: %v", userId, err))
		return fmt.Errorf("error checking for an existing order: %w", err)
	}
	logging.LogInfo(fmt.Sprintf("No pending order found for user %d", userId))

	// at this point we know there are no pending orders

	var cartItems []repo.CartItem
	getCartItemsQuery := `
		SELECT *
		FROM cart_items
		WHERE cart_id = $1 AND is_processed = FALSE
	`
	err = tx.SelectContext(ctx, &cartItems, getCartItemsQuery, cartId)
	if err != nil {
		logging.LogError(fmt.Sprintf("Failed to get cart items for cart %d: %v", cartId, err))
		return err 
	}

	if len(cartItems) == 0 {
		logging.LogInfo(fmt.Sprintf("No items in cart %d for user %d", cartId, userId))
		return errors.New("no items in cart")
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
			logging.LogError(fmt.Sprintf("Failed to get price for product %d: %v", item.ProductId, err))
			return err 
		}

		itemTotal := currentPrice * float64(item.Quantity)
		totalPrice += itemTotal

		orderItems = append(orderItems, repo.OrderItem{
			ProductId: item.ProductId,
			Price: currentPrice,
			Quantity: item.Quantity,
		})
		logging.LogInfo(fmt.Sprintf("Added product %d (qty %d, price %.2f) to order", item.ProductId, item.Quantity, currentPrice))
	}

	// Create the order with expiration time for price validity
	expirationTime := time.Now().Add(30 * time.Minute)
	order := &repo.Order{
		CustomerId:     userId,
		TotalPrice:     totalPrice,
		Status:         repo.OrderStatusPending,
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
		logging.LogError(fmt.Sprintf("Failed to create order for user %d: %v", userId, err))
		return fmt.Errorf("failed to create order: %w", err)
	}
	logging.LogInfo(fmt.Sprintf("Order %d created for user %d", order.Id, userId))
	
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
			logging.LogError(fmt.Sprintf("Failed to create order item for order %d, product %d: %v", order.Id, item.ProductId, err))
			return fmt.Errorf("failed to create order item: %w", err)
		}
		logging.LogInfo(fmt.Sprintf("Order item for product %d added to order %d", item.ProductId, order.Id))
	}
	
	// Mark cart items as processed
	updateCartItemsQuery := `
		UPDATE cart_items
		SET is_processed = TRUE
		WHERE cart_id = $1
	`
	_, err = tx.ExecContext(ctx, updateCartItemsQuery, cartId)
	if err != nil {
		logging.LogError(fmt.Sprintf("Failed to update cart items for cart %d: %v", cartId, err))
		return fmt.Errorf("failed to update cart items: %w", err)
	}
	logging.LogInfo(fmt.Sprintf("Cart items for cart %d marked as processed", cartId))
	
	// Commit the transaction
	if err = tx.Commit(); err != nil {
		logging.LogError(fmt.Sprintf("Failed to commit transaction for order %d: %v", order.Id, err))
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	logging.LogInfo(fmt.Sprintf("Order generation completed successfully for user %d", userId))
	
	return nil
}

func (s *OrderService) GetOrderWithItems (ctx context.Context, userId int) (*repo.OrderWithItems, error) {
	var order repo.Order
	getOrderQuery := `
		SELECT id, customer_id, total_price, status, price_valid_until, created_at
		FROM orders
		WHERE customer_id = $1 AND status = $2
		LIMIT 1
	`
	err := s.db.GetContext(ctx, &order, getOrderQuery, userId, repo.OrderStatusPending)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("no pending order found")
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	var items []repo.OrderItemDetail
	getItemsQuery := `
		SELECT oi.id, oi.order_id, oi.product_id, oi.price, oi.quantity, oi.created_at, oi.updated_at,
				p.name as product_name, p.image_path as product_image_path
		FROM order_items oi
		JOIN products p ON oi.product_id = p.id
		WHERE order_id = $1
		ORDER BY oi.created_at 
	`
	err = s.db.SelectContext(ctx, &items, getItemsQuery, order.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}

	orderWithItems := &repo.OrderWithItems{
		Id: order.Id,
		CustomerId: order.CustomerId,
		TotalPrice: order.TotalPrice,
		Status: order.Status,
		PriceValidUntil: order.PriceValidUntil,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
		Items: items,
	}

	return orderWithItems, nil
}
