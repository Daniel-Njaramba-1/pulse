package adminSvc

import (
	"context"
	"database/sql"

	"github.com/Daniel-Njaramba-1/pulse/internal/repo"
	"github.com/jmoiron/sqlx"
)

type StockService struct {
	db *sqlx.DB
}

func NewStockService(db *sqlx.DB) StockService {
	return StockService{db: db}
}

func (s *StockService) StockUpProduct(ctx context.Context, stock *repo.Stock) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Validate product
	var productId int
	getProductQuery := `
		SELECT id
		FROM products
		WHERE id = $1
	`
	err = tx.QueryRowContext(ctx, getProductQuery, stock.ProductId).Scan(&productId)
	if err != nil {
		if err == sql.ErrNoRows {
			return err // Or a custom error indicating product not found
		}
		return err
	}

	// Check if stock for the product already exists
	var existingStock repo.Stock
	getStockQuery := `
		SELECT id, product_id, quantity, stock_threshold, created_at, updated_at
		FROM stocks
		WHERE product_id = $1
	`
	err = tx.GetContext(ctx, &existingStock, getStockQuery, stock.ProductId)

	if err == nil {
		// Update existing stock
		updateStockQuery := `
			UPDATE stocks
			SET quantity = quantity + $1, updated_at = NOW()
			WHERE product_id = $2
		`
		_, err = tx.ExecContext(ctx, updateStockQuery, stock.Quantity, stock.ProductId)
		if err != nil {
			return err
		}
	} else if err == sql.ErrNoRows {
		// Create new stock level
		insertStockQuery := `
			INSERT INTO stocks(product_id, quantity, stock_threshold)
			VALUES ($1, $2, $3)
			RETURNING id
		`
		err = tx.QueryRowContext(ctx, insertStockQuery, stock.ProductId, stock.Quantity, stock.StockThreshold).Scan(&stock.Id)
		if err != nil {
			return err
		}
	} else {
		return err
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}