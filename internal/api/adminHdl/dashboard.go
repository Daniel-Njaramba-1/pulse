package adminHdl

import (
	"net/http"
	"strconv"

	"github.com/Daniel-Njaramba-1/pulse/internal/services/adminSvc"
	"github.com/Daniel-Njaramba-1/pulse/internal/util/logging"
	"github.com/labstack/echo/v4"
)

type DashboardHandler struct {
	dashboardService *adminSvc.DashboardService
}

func NewDashboardHandler(dashboardService *adminSvc.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService}
}

// Response structures for API
type DashboardResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type ComprehensiveAnalytics struct {
	ModelPerformance   interface{} `json:"model_performance"`
	SalesAnalytics     interface{} `json:"sales_analytics"`
	InventoryStatus    interface{} `json:"inventory_status"`
	PricingAnalytics   interface{} `json:"pricing_analytics"`
	CustomerBehavior   interface{} `json:"customer_behavior"`
	OperationalHealth  interface{} `json:"operational_health"`
	TopProducts        interface{} `json:"top_products"`
	CategoryRevenue    interface{} `json:"category_revenue"`
}

// GetCoefficients - Handler for fetching latest regression model coefficients
func (h *DashboardHandler) GetCoefficients(c echo.Context) error {
	logging.LogInfo("GetCoefficients: init")
	data, err := h.dashboardService.GetCoefficients()
	if err != nil {
		logging.LogError("GetCoefficients: Failed to fetch coefficients: " + err.Error())
		return c.JSON(http.StatusInternalServerError, DashboardResponse{
			Success: false,
			Error:   "Failed to fetch coefficients: " + err.Error(),
		})
	}

	logging.LogInfo("GetCoefficients: success")
	return c.JSON(http.StatusOK, DashboardResponse{
		Success: true,
		Message: "Latest regression model coefficients retrieved successfully",
		Data:    data,
	})
}

// GetAnalytics - Comprehensive dashboard analytics
func (h *DashboardHandler) GetAnalytics(c echo.Context) error {
	logging.LogInfo("GetAnalytics: init")
	// Parse query parameters with defaults
	daysParam := c.QueryParam("days")
	days := 30 // default
	if daysParam != "" {
		if parsedDays, err := strconv.Atoi(daysParam); err == nil && parsedDays > 0 {
			days = parsedDays
		}
	}

	limitParam := c.QueryParam("limit")
	limit := 10 // default
	if limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Fetch all analytics data
	analytics := ComprehensiveAnalytics{}

	// Model Performance
	logging.LogInfo("GetAnalytics: fetching model performance")
	if modelData, err := h.dashboardService.ViewModelPerformance(); err != nil {
		logging.LogError("GetAnalytics: Failed to fetch model performance data: " + err.Error())
		return c.JSON(http.StatusInternalServerError, DashboardResponse{
			Success: false,
			Error:   "Failed to fetch model performance data: " + err.Error(),
		})
	} else {
		analytics.ModelPerformance = modelData
	}

	// Sales Analytics
	logging.LogInfo("GetAnalytics: fetching sales analytics")
	if salesData, err := h.dashboardService.AnalyseSales(days); err != nil {
		logging.LogError("GetAnalytics: Failed to fetch sales analytics: " + err.Error())
		return c.JSON(http.StatusInternalServerError, DashboardResponse{
			Success: false,
			Error:   "Failed to fetch sales analytics: " + err.Error(),
		})
	} else {
		analytics.SalesAnalytics = salesData
	}

	// Inventory Status
	logging.LogInfo("GetAnalytics: fetching inventory status")
	if inventoryData, err := h.dashboardService.ManageInventory(); err != nil {
		logging.LogError("GetAnalytics: Failed to fetch inventory data: " + err.Error())
		return c.JSON(http.StatusInternalServerError, DashboardResponse{
			Success: false,
			Error:   "Failed to fetch inventory data: " + err.Error(),
		})
	} else {
		analytics.InventoryStatus = inventoryData
	}

	// Pricing Analytics
	logging.LogInfo("GetAnalytics: fetching pricing analytics")
	if pricingData, err := h.dashboardService.AnalyseDynamicPricing(days); err != nil {
		logging.LogError("GetAnalytics: Failed to fetch pricing analytics: " + err.Error())
		return c.JSON(http.StatusInternalServerError, DashboardResponse{
			Success: false,
			Error:   "Failed to fetch pricing analytics: " + err.Error(),
		})
	} else {
		analytics.PricingAnalytics = pricingData
	}

	// Customer Behavior
	logging.LogInfo("GetAnalytics: fetching customer behavior")
	if customerData, err := h.dashboardService.AnalyseCustomerBehavior(); err != nil {
		logging.LogError("GetAnalytics: Failed to fetch customer behavior data: " + err.Error())
		return c.JSON(http.StatusInternalServerError, DashboardResponse{
			Success: false,
			Error:   "Failed to fetch customer behavior data: " + err.Error(),
		})
	} else {
		analytics.CustomerBehavior = customerData
	}

	// Operational Health
	logging.LogInfo("GetAnalytics: fetching operational health")
	if healthData, err := h.dashboardService.ViewOperationalHealth(); err != nil {
		logging.LogError("GetAnalytics: Failed to fetch operational health data: " + err.Error())
		return c.JSON(http.StatusInternalServerError, DashboardResponse{
			Success: false,
			Error:   "Failed to fetch operational health data: " + err.Error(),
		})
	} else {
		analytics.OperationalHealth = healthData
	}

	// Top Products
	logging.LogInfo("GetAnalytics: fetching top products")
	if topProducts, err := h.dashboardService.GetTopSellingProducts(limit); err != nil {
		logging.LogError("GetAnalytics: Failed to fetch top products: " + err.Error())
		return c.JSON(http.StatusInternalServerError, DashboardResponse{
			Success: false,
			Error:   "Failed to fetch top products: " + err.Error(),
		})
	} else {
		analytics.TopProducts = topProducts
	}

	// Category Revenue
	logging.LogInfo("GetAnalytics: fetching category revenue")
	if categoryRevenue, err := h.dashboardService.GetRevenueByCategory(days); err != nil {
		logging.LogError("GetAnalytics: Failed to fetch category revenue: " + err.Error())
		return c.JSON(http.StatusInternalServerError, DashboardResponse{
			Success: false,
			Error:   "Failed to fetch category revenue: " + err.Error(),
		})
	} else {
		analytics.CategoryRevenue = categoryRevenue
	}

	logging.LogInfo("GetAnalytics: success")
	return c.JSON(http.StatusOK, DashboardResponse{
		Success: true,
		Message: "Dashboard analytics retrieved successfully",
		Data:    analytics,
	})
}

// GetModelPerformance - Specific model performance endpoint
func (h *DashboardHandler) GetModelPerformance(c echo.Context) error {
	logging.LogInfo("GetModelPerformance: init")
	data, err := h.dashboardService.ViewModelPerformance()
	if err != nil {
		logging.LogError("GetModelPerformance: Failed to fetch model performance: " + err.Error())
		return c.JSON(http.StatusInternalServerError, DashboardResponse{
			Success: false,
			Error:   "Failed to fetch model performance: " + err.Error(),
		})
	}

	logging.LogInfo("GetModelPerformance: success")
	return c.JSON(http.StatusOK, DashboardResponse{
		Success: true,
		Message: "Model performance data retrieved successfully",
		Data:    data,
	})
}

// GetSalesAnalytics - Specific sales analytics endpoint
func (h *DashboardHandler) GetSalesAnalytics(c echo.Context) error {
	logging.LogInfo("GetSalesAnalytics: init")
	daysParam := c.QueryParam("days")
	days := 30
	if daysParam != "" {
		if parsedDays, err := strconv.Atoi(daysParam); err == nil && parsedDays > 0 {
			days = parsedDays
		}
	}

	data, err := h.dashboardService.AnalyseSales(days)
	if err != nil {
		logging.LogError("GetSalesAnalytics: Failed to fetch sales analytics: " + err.Error())
		return c.JSON(http.StatusInternalServerError, DashboardResponse{
			Success: false,
			Error:   "Failed to fetch sales analytics: " + err.Error(),
		})
	}

	logging.LogInfo("GetSalesAnalytics: success")
	return c.JSON(http.StatusOK, DashboardResponse{
		Success: true,
		Message: "Sales analytics retrieved successfully",
		Data:    data,
	})
}

// GetInventoryStatus - Specific inventory status endpoint
func (h *DashboardHandler) GetInventoryStatus(c echo.Context) error {
	logging.LogInfo("GetInventoryStatus: init")
	data, err := h.dashboardService.ManageInventory()
	if err != nil {
		logging.LogError("GetInventoryStatus: Failed to fetch inventory status: " + err.Error())
		return c.JSON(http.StatusInternalServerError, DashboardResponse{
			Success: false,
			Error:   "Failed to fetch inventory status: " + err.Error(),
		})
	}

	logging.LogInfo("GetInventoryStatus: success")
	return c.JSON(http.StatusOK, DashboardResponse{
		Success: true,
		Message: "Inventory status retrieved successfully",
		Data:    data,
	})
}

// GetPricingAnalytics - Specific pricing analytics endpoint
func (h *DashboardHandler) GetPricingAnalytics(c echo.Context) error {
	logging.LogInfo("GetPricingAnalytics: init")
	daysParam := c.QueryParam("days")
	days := 7 // shorter default for pricing changes
	if daysParam != "" {
		if parsedDays, err := strconv.Atoi(daysParam); err == nil && parsedDays > 0 {
			days = parsedDays
		}
	}

	data, err := h.dashboardService.AnalyseDynamicPricing(days)
	if err != nil {
		logging.LogError("GetPricingAnalytics: Failed to fetch pricing analytics: " + err.Error())
		return c.JSON(http.StatusInternalServerError, DashboardResponse{
			Success: false,
			Error:   "Failed to fetch pricing analytics: " + err.Error(),
		})
	}

	logging.LogInfo("GetPricingAnalytics: success")
	return c.JSON(http.StatusOK, DashboardResponse{
		Success: true,
		Message: "Pricing analytics retrieved successfully",
		Data:    data,
	})
}

// GetCustomerBehavior - Specific customer behavior endpoint
func (h *DashboardHandler) GetCustomerBehavior(c echo.Context) error {
	logging.LogInfo("GetCustomerBehavior: init")
	data, err := h.dashboardService.AnalyseCustomerBehavior()
	if err != nil {
		logging.LogError("GetCustomerBehavior: Failed to fetch customer behavior data: " + err.Error())
		return c.JSON(http.StatusInternalServerError, DashboardResponse{
			Success: false,
			Error:   "Failed to fetch customer behavior data: " + err.Error(),
		})
	}

	logging.LogInfo("GetCustomerBehavior: success")
	return c.JSON(http.StatusOK, DashboardResponse{
		Success: true,
		Message: "Customer behavior data retrieved successfully",
		Data:    data,
	})
}

// GetOperationalHealth - Specific operational health endpoint
func (h *DashboardHandler) GetOperationalHealth(c echo.Context) error {
	logging.LogInfo("GetOperationalHealth: init")
	data, err := h.dashboardService.ViewOperationalHealth()
	if err != nil {
		logging.LogError("GetOperationalHealth: Failed to fetch operational health data: " + err.Error())
		return c.JSON(http.StatusInternalServerError, DashboardResponse{
			Success: false,
			Error:   "Failed to fetch operational health data: " + err.Error(),
		})
	}

	logging.LogInfo("GetOperationalHealth: success")
	return c.JSON(http.StatusOK, DashboardResponse{
		Success: true,
		Message: "Operational health data retrieved successfully",
		Data:    data,
	})
}

// GetTopProducts - Top selling products endpoint
func (h *DashboardHandler) GetTopProducts(c echo.Context) error {
	logging.LogInfo("GetTopProducts: init")
	limitParam := c.QueryParam("limit")
	limit := 10
	if limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	data, err := h.dashboardService.GetTopSellingProducts(limit)
	if err != nil {
		logging.LogError("GetTopProducts: Failed to fetch top products: " + err.Error())
		return c.JSON(http.StatusInternalServerError, DashboardResponse{
			Success: false,
			Error:   "Failed to fetch top products: " + err.Error(),
		})
	}

	logging.LogInfo("GetTopProducts: success")
	return c.JSON(http.StatusOK, DashboardResponse{
		Success: true,
		Message: "Top products retrieved successfully",
		Data:    data,
	})
}

// GetCategoryRevenue - Category revenue breakdown endpoint
func (h *DashboardHandler) GetCategoryRevenue(c echo.Context) error {
	logging.LogInfo("GetCategoryRevenue: init")
	daysParam := c.QueryParam("days")
	days := 30
	if daysParam != "" {
		if parsedDays, err := strconv.Atoi(daysParam); err == nil && parsedDays > 0 {
			days = parsedDays
		}
	}

	data, err := h.dashboardService.GetRevenueByCategory(days)
	if err != nil {
		logging.LogError("GetCategoryRevenue: Failed to fetch category revenue: " + err.Error())
		return c.JSON(http.StatusInternalServerError, DashboardResponse{
			Success: false,
			Error:   "Failed to fetch category revenue: " + err.Error(),
		})
	}

	logging.LogInfo("GetCategoryRevenue: success")
	return c.JSON(http.StatusOK, DashboardResponse{
		Success: true,
		Message: "Category revenue data retrieved successfully",
		Data:    data,
	})
}
