package adminSvc

import (
	"time"

	"github.com/Daniel-Njaramba-1/pulse/internal/util/logging"
	"github.com/jmoiron/sqlx"
)

type DashboardService struct {
	db *sqlx.DB
}

func NewDashboardService(db *sqlx.DB) *DashboardService {
	return &DashboardService{db: db}
}

// Data structures for dashboard responses
type ModelPerformanceData struct {
	ModelVersion   string    `json:"model_version" db:"model_version"`
	TrainingDate   time.Time `json:"training_date" db:"training_date"`
	RSquared       float64   `json:"r_squared" db:"r_squared"`
	MSE            float64   `json:"mse" db:"mse"`
	RMSE           float64   `json:"rmse" db:"rmse"`
	MAE            float64   `json:"mae" db:"mae"`
	SampleSize     int       `json:"sample_size" db:"sample_size"`
}

type SalesAnalytics struct {
	Date         time.Time `json:"date" db:"date"`
	ProductID    int       `json:"product_id" db:"product_id"`
	ProductName  string    `json:"product_name" db:"product_name"`
	TotalSales   int       `json:"total_sales" db:"total_sales"`
	TotalRevenue float64   `json:"total_revenue" db:"total_revenue"`
	CategoryName string    `json:"category_name" db:"category_name"`
	BrandName    string    `json:"brand_name" db:"brand_name"`
}

type InventoryStatus struct {
	ProductID       int       `json:"product_id" db:"product_id"`
	ProductName     string    `json:"product_name" db:"product_name"`
	CurrentStock    int       `json:"current_stock" db:"current_stock"`
	StockThreshold  int       `json:"stock_threshold" db:"stock_threshold"`
	IsLowStock      bool      `json:"is_low_stock" db:"is_low_stock"`
	LastRestock     time.Time `json:"last_restock" db:"last_restock"`
	DaysSinceRestock int      `json:"days_since_restock" db:"days_since_restock"`
}

type PricingAnalytics struct {
	ProductID     int       `json:"product_id" db:"product_id"`
	ProductName   string    `json:"product_name" db:"product_name"`
	BasePrice     float64   `json:"base_price" db:"base_price"`
	AdjustedPrice float64   `json:"adjusted_price" db:"adjusted_price"`
	PriceChange   float64   `json:"price_change" db:"price_change"`
	LastAdjusted  time.Time `json:"last_adjusted" db:"last_adjusted"`
	ModelVersion  string    `json:"model_version" db:"model_version"`
}

type CustomerBehavior struct {
	ProductID              int     `json:"product_id" db:"product_id"`
	ProductName            string  `json:"product_name" db:"product_name"`
	AverageRating          float64 `json:"average_rating" db:"average_rating"`
	ReviewCount            int     `json:"review_count" db:"review_count"`
	WishlistCount          int     `json:"wishlist_count" db:"wishlist_count"`
	WishlistToSalesRatio   float64 `json:"wishlist_to_sales_ratio" db:"wishlist_to_sales_ratio"`
	SalesVelocity          float64 `json:"sales_velocity" db:"sales_velocity"`
}

type OperationalHealth struct {
	Metric string  `json:"metric" db:"metric"`
	Value  float64 `json:"value" db:"value"`
	Status string  `json:"status" db:"status"`
}

// ViewModelPerformance - Track ML model performance over time
func (s *DashboardService) ViewModelPerformance() ([]ModelPerformanceData, error) {
	logging.LogInfo("DashboardService: ViewModelPerformance called")
	query := `
		SELECT 
			model_version,
			training_date,
			COALESCE(r_squared, 0.0) as r_squared,
			COALESCE(mse, 0.0) as mse,
			COALESCE(rmse, 0.0) as rmse,
			COALESCE(mae, 0.0) as mae,
			sample_size
		FROM price_model_coefficients 
		ORDER BY training_date DESC
		LIMIT 10
	`
	var results []ModelPerformanceData
	err := s.db.Select(&results, query)
	if err != nil {
		logging.LogError("DashboardService: ViewModelPerformance error: " + err.Error())
	} else {
		logging.LogInfo("DashboardService: ViewModelPerformance success")
	}
	return results, err
}

// AnalyseSales - Comprehensive sales analytics
func (s *DashboardService) AnalyseSales(days int) ([]SalesAnalytics, error) {
	logging.LogInfo("DashboardService: AnalyseSales called with days=%d", days)
	query := `
		SELECT 
			DATE(s.created_at) as date,
			s.product_id,
			p.name as product_name,
			COUNT(s.id) as total_sales,
			SUM(s.sale_price * s.quantity) as total_revenue,
			c.name as category_name,
			b.name as brand_name
		FROM sales s
		JOIN products p ON s.product_id = p.id
		JOIN categories c ON p.category_id = c.id
		JOIN brands b ON p.brand_id = b.id
		WHERE s.created_at >= NOW() - ($1 || ' days')::interval
		GROUP BY DATE(s.created_at), s.product_id, p.name, c.name, b.name
		ORDER BY date DESC, total_revenue DESC
	`
	var results []SalesAnalytics
	err := s.db.Select(&results, query, days)
	if err != nil {
		logging.LogError("DashboardService: AnalyseSales error: " + err.Error())
	} else {
		logging.LogInfo("DashboardService: AnalyseSales success")
	}
	return results, err
}

// ManageInventory - Current inventory status and alerts
func (s *DashboardService) ManageInventory() ([]InventoryStatus, error) {
	logging.LogInfo("DashboardService: ManageInventory called")
	query := `
		SELECT 
			st.product_id,
			p.name as product_name,
			st.quantity as current_stock,
			st.stock_threshold,
			CASE WHEN st.quantity <= st.stock_threshold THEN true ELSE false END as is_low_stock,
			COALESCE(
				(SELECT created_at 
				 FROM stock_history sh 
				 WHERE sh.product_id = st.product_id 
				 AND sh.event_type = 'restock' 
				 ORDER BY created_at DESC 
				 LIMIT 1), 
				st.created_at
			) as last_restock,
			COALESCE(
				EXTRACT(DAY FROM NOW() - 
					(SELECT created_at 
					 FROM stock_history sh 
					 WHERE sh.product_id = st.product_id 
					 AND sh.event_type = 'restock' 
					 ORDER BY created_at DESC 
					 LIMIT 1)
				)::INTEGER, 
				0
			) as days_since_restock
		FROM stocks st
		JOIN products p ON st.product_id = p.id
		WHERE p.is_active = true
		ORDER BY is_low_stock DESC, current_stock ASC
	`
	var results []InventoryStatus
	err := s.db.Select(&results, query)
	if err != nil {
		logging.LogError("DashboardService: ManageInventory error: " + err.Error())
	} else {
		logging.LogInfo("DashboardService: ManageInventory success")
	}
	return results, err
}

// AnalyseDynamicPricing - Price adjustment analytics
func (s *DashboardService) AnalyseDynamicPricing(days int) ([]PricingAnalytics, error) {
	logging.LogInfo("DashboardService: AnalyseDynamicPricing called with days=%d", days)
	query := `
		SELECT 
			pa.product_id,
			p.name as product_name,
			pm.base_price,
			pm.adjusted_price,
			(pm.adjusted_price - pm.base_price) as price_change,
			pa.created_at as last_adjusted,
			pa.model_version
		FROM price_adjustments pa
		JOIN products p ON pa.product_id = p.id
		JOIN product_metrics pm ON pa.product_id = pm.product_id
		WHERE pa.created_at >= NOW() - ($1 || ' days')::interval
		ORDER BY pa.created_at DESC
	`
	var results []PricingAnalytics
	err := s.db.Select(&results, query, days)
	if err != nil {
		logging.LogError("DashboardService: AnalyseDynamicPricing error: " + err.Error())
	} else {
		logging.LogInfo("DashboardService: AnalyseDynamicPricing success")
	}
	return results, err
}

// AnalyseCustomerBehavior - Customer interaction patterns
func (s *DashboardService) AnalyseCustomerBehavior() ([]CustomerBehavior, error) {
	logging.LogInfo("DashboardService: AnalyseCustomerBehavior called")
	query := `
		SELECT 
			pm.product_id,
			p.name as product_name,
			pm.average_rating,
			pm.review_count,
			pm.wishlist_count,
			pf.wishlist_to_sales_ratio,
			pf.sales_velocity
		FROM product_metrics pm
		JOIN products p ON pm.product_id = p.id
		LEFT JOIN pricing_features pf ON pm.product_id = pf.product_id
		WHERE p.is_active = true
		ORDER BY pm.average_rating DESC, pm.review_count DESC
	`
	var results []CustomerBehavior
	err := s.db.Select(&results, query)
	if err != nil {
		logging.LogError("DashboardService: AnalyseCustomerBehavior error: " + err.Error())
	} else {
		logging.LogInfo("DashboardService: AnalyseCustomerBehavior success")
	}
	return results, err
}

// ViewOperationalHealth - System health metrics
func (s *DashboardService) ViewOperationalHealth() ([]OperationalHealth, error) {
	logging.LogInfo("DashboardService: ViewOperationalHealth called")
	query := `
		SELECT * FROM (
			-- Order completion rate
			SELECT 
				'Order Completion Rate' as metric,
				ROUND(
					(COUNT(CASE WHEN status = 'completed' THEN 1 END) * 100.0 / COUNT(*)), 2
				) as value,
				CASE 
					WHEN (COUNT(CASE WHEN status = 'completed' THEN 1 END) * 100.0 / COUNT(*)) >= 80 
					THEN 'Good' 
					ELSE 'Needs Attention' 
				END as status
			FROM orders 
			WHERE created_at >= NOW() - INTERVAL '7 days'
			
			UNION ALL
			
			-- Payment success rate
			SELECT 
				'Payment Success Rate' as metric,
				ROUND(
					(COUNT(CASE WHEN status = 'completed' THEN 1 END) * 100.0 / COUNT(*)), 2
				) as value,
				CASE 
					WHEN (COUNT(CASE WHEN status = 'completed' THEN 1 END) * 100.0 / COUNT(*)) >= 90 
					THEN 'Good' 
					ELSE 'Needs Attention' 
				END as status
			FROM payments 
			WHERE created_at >= NOW() - INTERVAL '7 days'
			
			UNION ALL
			
			-- Low stock products percentage
			SELECT 
				'Low Stock Products %' as metric,
				ROUND(
					(COUNT(CASE WHEN quantity <= stock_threshold THEN 1 END) * 100.0 / COUNT(*)), 2
				) as value,
				CASE 
					WHEN (COUNT(CASE WHEN quantity <= stock_threshold THEN 1 END) * 100.0 / COUNT(*)) <= 20 
					THEN 'Good' 
					ELSE 'Needs Attention' 
				END as status
			FROM stocks s
			JOIN products p ON s.product_id = p.id
			WHERE p.is_active = true
			
			UNION ALL
			
			-- Active customers (last 30 days)
			SELECT 
				'Active Customers (30d)' as metric,
				COUNT(DISTINCT customer_id)::FLOAT as value,
				'Info' as status
			FROM orders 
			WHERE created_at >= NOW() - INTERVAL '30 days'
		) metrics
	`
	var results []OperationalHealth
	err := s.db.Select(&results, query)
	if err != nil {
		logging.LogError("DashboardService: ViewOperationalHealth error: " + err.Error())
	} else {
		logging.LogInfo("DashboardService: ViewOperationalHealth success")
	}
	return results, err
}

// GetTopSellingProducts - Best performing products
func (s *DashboardService) GetTopSellingProducts(limit int) ([]SalesAnalytics, error) {
	logging.LogInfo("DashboardService: GetTopSellingProducts called with limit=%d", limit)
	query := `
		SELECT 
			s.product_id,
			p.name as product_name,
			COUNT(s.id) as total_sales,
			SUM(s.sale_price * s.quantity) as total_revenue,
			c.name as category_name,
			b.name as brand_name,
			MAX(s.created_at) as date
		FROM sales s
		JOIN products p ON s.product_id = p.id
		JOIN categories c ON p.category_id = c.id
		JOIN brands b ON p.brand_id = b.id
		WHERE s.created_at >= NOW() - INTERVAL '30 days'
		GROUP BY s.product_id, p.name, c.name, b.name
		ORDER BY total_revenue DESC
		LIMIT $1
	`
	var results []SalesAnalytics
	err := s.db.Select(&results, query, limit)
	if err != nil {
		logging.LogError("DashboardService: GetTopSellingProducts error: " + err.Error())
	} else {
		logging.LogInfo("DashboardService: GetTopSellingProducts success")
	}
	return results, err
}

// GetRevenueByCategory - Revenue breakdown by category
func (s *DashboardService) GetRevenueByCategory(days int) ([]struct {
	CategoryName string  `json:"category_name" db:"category_name"`
	TotalRevenue float64 `json:"total_revenue" db:"total_revenue"`
	SalesCount   int     `json:"sales_count" db:"sales_count"`
}, error) {
	logging.LogInfo("DashboardService: GetRevenueByCategory called with days=%d", days)
	query := `
		SELECT 
			c.name as category_name,
			SUM(s.sale_price * s.quantity) as total_revenue,
			COUNT(s.id) as sales_count
		FROM sales s
		JOIN products p ON s.product_id = p.id
		JOIN categories c ON p.category_id = c.id
		WHERE s.created_at >= NOW() - ($1 || ' days')::interval
		GROUP BY c.name
		ORDER BY total_revenue DESC
	`
	var results []struct {
		CategoryName string  `json:"category_name" db:"category_name"`
		TotalRevenue float64 `json:"total_revenue" db:"total_revenue"`
		SalesCount   int     `json:"sales_count" db:"sales_count"`
	}

	err := s.db.Select(&results, query, days)
	if err != nil {
		logging.LogError("DashboardService: GetRevenueByCategory error: " + err.Error())
	} else {
		logging.LogInfo("DashboardService: GetRevenueByCategory success")
	}
	return results, err
}