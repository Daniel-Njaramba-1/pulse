package pricing

import "context"

// Coalesce returns the first non-null value in list of arguments,
// Number of days since last sale of a product
// If product has never been sold, it returns 100 as fallback
func (s *ModelService) GetDaysSinceLastSale(ctx context.Context, productId int) (int, error) {
	query := `
        SELECT COALESCE(EXTRACT(DAY FROM NOW() - MAX(created_at)), 100) as days_since_sale
        FROM sales 
        WHERE product_id = $1
    `
	var daysSinceSale int
	err := s.db.GetContext(ctx, &daysSinceSale, query, productId)
	return daysSinceSale, err
}

// sales velocity  - average sales for last 30 days, return 0 if null
// filters sales from last last 30 days
func (s *ModelService) CalculateSalesVelocity(ctx context.Context, productId int) (float64, error) {
	query := `
        SELECT COALESCE(
          COUNT(*) / NULLIF(EXTRACT(DAY FROM NOW() - MIN(created_at)), 0), 
          0
        ) as sales_velocity
        FROM sales 
        WHERE product_id = $1 
          AND created_at > NOW() - INTERVAL '30 days'
    `
	var velocity float64
	err := s.db.GetContext(ctx, &velocity, query, productId)
	return velocity, err
}

func (s *ModelService) CalculateTotalSalesCount(ctx context.Context, productId int) (int, error) {
	query := `
        SELECT COUNT(*) FROM sales WHERE product_id = $1
    `
	var count int
	err := s.db.GetContext(ctx, &count, query, productId)
	return count, err
}

func (s *ModelService) CalculateTotalSalesValue(ctx context.Context, productId int) (float64, error) {
	query := `
        SELECT COALESCE(SUM(sale_price * quantity), 0) FROM sales WHERE product_id = $1
    `
	var value float64
	err := s.db.GetContext(ctx, &value, query, productId)
	return value, err
}

// return int rank for product, 0 if not found / inactive
// Coalesce makes sure NULL ratings dont break ranking logic
func (s *ModelService) CalculateCategoryRank(ctx context.Context, productId int) (int, error) {
	query := `
        WITH product_category AS (
            SELECT category_id FROM products WHERE id = $1
        ),
        ranked_products AS (
            SELECT p.id, 
                  RANK() OVER (
                    PARTITION BY p.category_id 
                    ORDER BY COALESCE(pm.average_rating, 0) DESC, 
                            COALESCE(pm.review_count, 0) DESC
                  ) as rank
            FROM products p
            JOIN product_metrics pm ON p.id = pm.product_id
            WHERE p.category_id = (SELECT category_id FROM product_category)
              AND p.is_active = true
        )
        SELECT rank FROM ranked_products WHERE id = $1
    `
	var rank int
	err := s.db.GetContext(ctx, &rank, query, productId)
	return rank, err
}

func (s *ModelService) CalculateCategoryPercentile(ctx context.Context, productId int) (float64, error) {
	// First get the category rank
	rank, err := s.CalculateCategoryRank(ctx, productId)
	if err != nil {
		return 0, err
	}

	// Then get the total products in the category
	query := `
        WITH product_category AS (
            SELECT category_id FROM products WHERE id = $1
        )
        SELECT COUNT(*) FROM products
        WHERE category_id = (SELECT category_id FROM product_category)
          AND is_active = true
    `
	var total int
	err = s.db.GetContext(ctx, &total, query, productId)
	if err != nil {
		return 0, err
	}

	// Calculate percentile (higher is better)
	if total == 0 {
		return 0, nil
	}
	return (1 - float64(rank)/float64(total)) * 100, nil
}

func (s *ModelService) CalculateReviewScore(ctx context.Context, productId int) (float64, error) {
	query := `
        SELECT COALESCE(average_rating, 0) FROM product_metrics WHERE product_id = $1
    `
	var score float64
	err := s.db.GetContext(ctx, &score, query, productId)
	return score, err
}

func (s *ModelService) CalculateWishlistToSalesRatio(ctx context.Context, productId int) (float64, error) {
	query := `
        WITH wishlist_count AS (
            SELECT COUNT(*) as count FROM wishlist_items WHERE product_id = $1
        ),
        sales_count AS (
            SELECT COUNT(*) as count FROM sales WHERE product_id = $1
        )
        SELECT 
            CASE 
                WHEN (SELECT count FROM sales_count) = 0 THEN 0
                ELSE (SELECT count FROM wishlist_count)::float / (SELECT count FROM sales_count)::float
            END as ratio
    `
	var ratio float64
	err := s.db.GetContext(ctx, &ratio, query, productId)
	return ratio, err
}

func (s *ModelService) CalculateDaysInStock(ctx context.Context, productId int) (int, error) {
	query := `
        WITH first_stock AS (
            SELECT first_stocked_date FROM stocks WHERE product_id = $1
        )
        SELECT COALESCE(EXTRACT(DAY FROM NOW() - first_stocked_date)::integer, 0) 
        FROM first_stock
    `
	var days int
	err := s.db.GetContext(ctx, &days, query, productId)
	return days, err
}

func (s *ModelService) CalculateSeasonalFactor(ctx context.Context, productId int) (float64, error) {
	// A simple seasonal factor based on month (can be made more sophisticated)
	query := `
        WITH product_category AS (
            SELECT category_id FROM products WHERE id = $1
        ),
        month_factor AS (
            -- Different categories might have different seasonal patterns
            SELECT 
                CASE 
                    WHEN EXTRACT(MONTH FROM CURRENT_DATE) IN (11, 12) THEN 1.2 -- Holiday season
                    WHEN EXTRACT(MONTH FROM CURRENT_DATE) IN (1, 2) THEN 0.8  -- Post-holiday slump
                    ELSE 1.0  -- Normal season
                END as factor
        )
        SELECT factor FROM month_factor
    `
	var factor float64
	err := s.db.GetContext(ctx, &factor, query, productId)
	return factor, err
}