package customerSvc

import (
	"context"
	"fmt"
	"time"

	"github.com/Daniel-Njaramba-1/pulse/internal/repo"
	"github.com/jmoiron/sqlx"
)

type PaymentService struct {
	db *sqlx.DB
}

func NewPaymentService (db *sqlx.DB) *PaymentService {
	return &PaymentService{db: db}
}

func (s *PaymentService) ProcessPayment(ctx context.Context, userId int) (string, error) {
	// Start transaction
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	// Get order details and check price validity
	var order repo.Order
	orderQuery := `
		SELECT *
		FROM orders
		WHERE customer_id = $1 AND status = 'pending'
		FOR UPDATE
	`
	err = tx.GetContext(ctx, &order, orderQuery, userId)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve order: %w", err)
	}

	// Check if price is still valid
	isPriceValid := time.Now().Before(order.PriceValidUntil)
	
	// Get order items
	var orderItems []repo.OrderItem
	itemsQuery := `
		SELECT * FROM order_items
		WHERE order_id = $1
	`
	err = tx.SelectContext(ctx, &orderItems, itemsQuery, order.Id)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve order items: %w", err)
	}

	// If price is no longer valid, recalculate prices
	if !isPriceValid {
		var newTotalPrice float64 = 0
		
		for i, item := range orderItems {
			// Get current price for each product
			var currentPrice float64
			priceQuery := `
				SELECT adjusted_price
				FROM product_metrics
				WHERE product_id = $1
			`
			err = tx.QueryRowxContext(ctx, priceQuery, item.ProductId).Scan(&currentPrice)
			if err != nil {
				return "", fmt.Errorf("failed to get current price for product %d: %w", item.ProductId, err)
			}
			
			// Update item price
			updateItemQuery := `
				UPDATE order_items
				SET price = $1
				WHERE id = $2
			`
			_, err = tx.ExecContext(ctx, updateItemQuery, currentPrice, item.Id)
			if err != nil {
				return "", fmt.Errorf("failed to update item price: %w", err)
			}
			
			// Update local item price for sales record later
			orderItems[i].Price = currentPrice
			
			// Add to new total
			newTotalPrice += currentPrice * float64(item.Quantity)
		}
		
		// Update order total price
		updateOrderQuery := `
			UPDATE orders
			SET total_price = $1, price_valid_until = $2
			WHERE id = $3
		`
		newValidUntil := time.Now().Add(30 * time.Minute)
		_, err = tx.ExecContext(ctx, updateOrderQuery, newTotalPrice, newValidUntil, order.Id)
		if err != nil {
			return "", fmt.Errorf("failed to update order price: %w", err)
		}
		
		// Update order in memory
		order.TotalPrice = newTotalPrice
		order.PriceValidUntil = newValidUntil
	}

	// Check if all products are in stock
	for _, item := range orderItems {
		var stockQuantity int
		stockQuery := `
			SELECT quantity FROM stocks 
			WHERE product_id = $1
			FOR UPDATE
		`
		err = tx.QueryRowxContext(ctx, stockQuery, item.ProductId).Scan(&stockQuantity)
		if err != nil {
			return "", fmt.Errorf("failed to check stock for product %d: %w", item.ProductId, err)
		}
		
		if stockQuantity < item.Quantity {
			return "", fmt.Errorf("insufficient stock for product %d: requested %d, available %d", 
				item.ProductId, item.Quantity, stockQuantity)
		}
	}

	// Create payment record
	paymentId := fmt.Sprintf("TXN-%d", time.Now().UnixNano())
	paymentQuery := `
		INSERT INTO payments (order_id, payment_method, amount, status, transaction_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	var paymentDbId int
	err = tx.QueryRowxContext(
		ctx,
		paymentQuery,
		order.Id,
		repo.PaymentMethodCard, // This would come from a parameter in a real implementation
		order.TotalPrice,
		repo.PaymentStatusSuccess,		// always success for now
		paymentId,
	).Scan(&paymentDbId)
	if err != nil {
		return "", fmt.Errorf("failed to create payment record: %w", err)
	}

	// Update order status
	_, err = tx.ExecContext(ctx, `
		UPDATE orders 
		SET status = $2 
		WHERE id = $1
	`, order.Id, repo.OrderStatusCompleted)
	if err != nil {
		return "", fmt.Errorf("failed to update order status: %w", err)
	}

	// Generate sales records and update stock for each item
	for _, item := range orderItems {
		// Create sales record
		_, err = tx.ExecContext(ctx, `
			INSERT INTO sales (order_item_id, product_id, sale_price, quantity)
			VALUES ($1, $2, $3, $4)
		`, item.Id, item.ProductId, item.Price, item.Quantity)
		if err != nil {
			return "", fmt.Errorf("failed to create sales record: %w", err)
		}
		
		// Update product metrics
		_, err = tx.ExecContext(ctx, `
			UPDATE product_metrics
			SET last_sale = NOW()
			WHERE product_id = $1
		`, item.ProductId)
		if err != nil {
			return "", fmt.Errorf("failed to update product metrics: %w", err)
		}
		
		// Update stock quantity
		_, err = tx.ExecContext(ctx, `
			UPDATE stocks
			SET quantity = quantity - $1
			WHERE product_id = $2
		`, item.Quantity, item.ProductId)
		if err != nil {
			return "", fmt.Errorf("failed to update stock: %w", err)
		}
		
		// Record stock history
		// _, err = tx.ExecContext(ctx, `
		// 	INSERT INTO stock_history (
		// 		product_id, event_type, quantity_change, quantity_after
		// 	) SELECT 
		// 		$1, 'sale', -$2, 
		// 		(SELECT quantity FROM stocks WHERE product_id = $1)
		// 	FROM stocks WHERE product_id = $1
		// `, item.ProductId, item.Quantity)
		// if err != nil {
		// 	return "", fmt.Errorf("failed to record stock history: %w", err)
		// }
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return paymentId, nil
}

