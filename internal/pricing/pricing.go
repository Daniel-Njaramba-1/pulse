package pricing

import (
    "fmt"
    "log"
    "math"

    "gonum.org/v1/gonum/stat"
	"github.com/sajari/regression"
)

type PricingModel struct {
	// Coefficients for the linear regression model
	Intercept   float64 // Base price component
	StockCoeff  float64 // How stock affects price
	SalesCoeff  float64 // How sales velocity affects price
	
	// Historical data points
	StockLevels []float64
	SalesRates  []float64
	Prices      []float64
	
	// Current state
	CurrentStock      int
	CurrentSalesRate  float64  // Sales per day/hour
	InitialBasePrice  float64
	
	// Minimum data points required for regression
	MinDataPoints int
}

// NewPricingModel creates a new pricing model with initial parameters
func NewPricingModel(basePrice float64, initialStock int) *PricingModel {
    return &PricingModel{
        Intercept:        basePrice,  // Start with base price as intercept
        StockCoeff:       -0.01,      // Default: price decreases as stock increases
        SalesCoeff:       0.05,       // Default: price increases as sales increase
        StockLevels:      []float64{},
        SalesRates:       []float64{},
        Prices:           []float64{},
        CurrentStock:     initialStock,
        CurrentSalesRate: 0,
        InitialBasePrice: basePrice,
        MinDataPoints:    5,          // Require more data for stability
    }
}

// AddDataPoint adds a new observation to the historical data
func (pm *PricingModel) AddDataPoint(stock int, salesRate float64, price float64) {
    pm.StockLevels = append(pm.StockLevels, float64(stock))
    pm.SalesRates = append(pm.SalesRates, salesRate)
    pm.Prices = append(pm.Prices, price)
    
    // Update current state
    pm.CurrentStock = stock
    pm.CurrentSalesRate = salesRate
    
    // Re-fit the model if we have enough data points
    if len(pm.StockLevels) >= pm.MinDataPoints {
        err := pm.FitModel()
        if err != nil {
            log.Printf("Warning: %v. Using simple pricing formula instead.", err)
        }
    }
}

// checkVariance returns true if there's enough variance in the data
func (pm *PricingModel) checkVariance(data []float64, threshold float64) bool {
    if len(data) < 2 {
        return false
    }
    
    // Calculate mean
    sum := 0.0
    for _, val := range data {
        sum += val
    }
    mean := sum / float64(len(data))
    
    // Calculate variance
    varSum := 0.0
    for _, val := range data {
        diff := val - mean
        varSum += diff * diff
    }
    variance := varSum / float64(len(data))
    
    return variance > threshold
}

// FitModel recalculates the linear regression coefficients based on historical data
func (pm *PricingModel) FitModel() error {
    if len(pm.StockLevels) < pm.MinDataPoints {
        return fmt.Errorf("not enough data points for regression, need at least %d", pm.MinDataPoints)
    }

    // Check for sufficient variance in data
    stockVarianceThreshold := 1.0
    salesVarianceThreshold := 0.01
    
    if !pm.checkVariance(pm.StockLevels, stockVarianceThreshold) {
        return fmt.Errorf("insufficient variance in stock levels")
    }
    
    if !pm.checkVariance(pm.SalesRates, salesVarianceThreshold) {
        return fmt.Errorf("insufficient variance in sales rates")
    }
    
    // Create weights (all equal for now)
    weights := make([]float64, len(pm.StockLevels))
    for i := range weights {
        weights[i] = 1.0
    }
    
    // Set up matrices for multiple regression
    // We'll use a simpler approach with separate regressions for now
    
    // First, regress price against stock levels
    origin := false
    alphaStock, betaStock := stat.LinearRegression(pm.StockLevels, pm.Prices, weights, origin)
    
    // Second, regress price against sales rates
    alphaSales, betaSales := stat.LinearRegression(pm.SalesRates, pm.Prices, weights, origin)
    
    // Check for NaN values
    if math.IsNaN(alphaStock) || math.IsNaN(betaStock) {
        return fmt.Errorf("stock regression produced NaN values")
    }
    
    if math.IsNaN(alphaSales) || math.IsNaN(betaSales) {
        return fmt.Errorf("sales regression produced NaN values")
    }
    
    // Update coefficients
    pm.Intercept = (alphaStock + alphaSales) / 2
    pm.StockCoeff = betaStock
    pm.SalesCoeff = betaSales
    
    // Apply reasonable constraints to coefficients
    // Stock coefficient should generally be negative (higher stock → lower price)
    if pm.StockCoeff > 0 {
        pm.StockCoeff = math.Min(pm.StockCoeff, 0.05) // Allow small positive effect
    } else {
        pm.StockCoeff = math.Max(pm.StockCoeff, -0.1) // Limit negative effect
    }
    
    // Sales coefficient should generally be positive (higher sales → higher price)
    if pm.SalesCoeff < 0 {
        pm.SalesCoeff = math.Max(pm.SalesCoeff, -0.1) // Allow small negative effect
    } else {
        pm.SalesCoeff = math.Min(pm.SalesCoeff, 0.3) // Limit positive effect
    }
    
    return nil
}

// RecordSale records a new sale, updates stock and sales rate, and recalculates pricing
func (pm *PricingModel) RecordSale(quantity int) (float64, error) {
    if quantity <= 0 {
        return 0, fmt.Errorf("quantity must be positive")
    }
    
    if pm.CurrentStock < quantity {
        return 0, fmt.Errorf("insufficient stock: have %d, requested %d", pm.CurrentStock, quantity)
    }
    
    // Update stock
    pm.CurrentStock -= quantity
    
    // Update sales rate (simple moving average)
    if len(pm.SalesRates) == 0 {
        pm.CurrentSalesRate = float64(quantity)
    } else {
        pm.CurrentSalesRate = (pm.CurrentSalesRate*0.8) + (float64(quantity)*0.2)
    }
    
    // Calculate new price
    price := pm.CalculatePrice()
    
    // Add a new data point with updated state
    pm.AddDataPoint(pm.CurrentStock, pm.CurrentSalesRate, price)
    
    return price, nil
}

// AddStock increases the current stock level
func (pm *PricingModel) AddStock(quantity int) (float64, error) {
    if quantity <= 0 {
        return 0, fmt.Errorf("quantity must be positive")
    }
    
    // Update stock
    pm.CurrentStock += quantity
    
    // Calculate new price
    price := pm.CalculatePrice()
    
    // Add a new data point with updated state
    pm.AddDataPoint(pm.CurrentStock, pm.CurrentSalesRate, price)
    
    return price, nil
}

// CalculatePrice computes the current price based on the regression model
func (pm *PricingModel) CalculatePrice() float64 {
    // If we don't have enough data for regression yet, or if there was an error fitting the model,
    // use a simple formula
    if len(pm.StockLevels) < pm.MinDataPoints || 
       math.IsNaN(pm.Intercept) || 
       math.IsNaN(pm.StockCoeff) || 
       math.IsNaN(pm.SalesCoeff) {
        // Simple initial model: reduce price as stock increases, increase with sales
        stockFactor := 1.0 - (float64(pm.CurrentStock) * 0.001)
        salesFactor := 1.0 + (pm.CurrentSalesRate * 0.01)
        
        // Ensure we don't go below 70% or above 130% of base price
        adjustmentFactor := math.Max(0.7, math.Min(1.3, stockFactor * salesFactor))
        return pm.InitialBasePrice * adjustmentFactor
    }
    
    // Use the regression model
    price := pm.Intercept + 
        (pm.StockCoeff * float64(pm.CurrentStock)) + 
        (pm.SalesCoeff * pm.CurrentSalesRate)
    
    // Ensure price doesn't go negative or too low (at least 50% of base price)
    minPrice := pm.InitialBasePrice * 0.5
    maxPrice := pm.InitialBasePrice * 1.5
    
    // Apply constraints
    price = math.Max(minPrice, price)
    price = math.Min(maxPrice, price)
    
    return price
}

// GetModelState returns the current state of the pricing model
func (pm *PricingModel) GetModelState() map[string]interface{} {
    // Check for NaN values to avoid JSON serialization errors
    intercept := pm.Intercept
    stockCoeff := pm.StockCoeff
    salesCoeff := pm.SalesCoeff
    
    if math.IsNaN(intercept) {
        intercept = pm.InitialBasePrice
    }
    
    if math.IsNaN(stockCoeff) {
        stockCoeff = -0.01
    }
    
    if math.IsNaN(salesCoeff) {
        salesCoeff = 0.05
    }
    
    return map[string]interface{}{
        "intercept": intercept,
        "stockCoefficient": stockCoeff,
        "salesCoefficient": salesCoeff,
        "currentStock": pm.CurrentStock,
        "currentSalesRate": pm.CurrentSalesRate,
        "currentPrice": pm.CalculatePrice(),
        "basePrice": pm.InitialBasePrice,
        "dataPoints": len(pm.StockLevels),
        "usingRegressionModel": len(pm.StockLevels) >= pm.MinDataPoints && 
                               !math.IsNaN(pm.Intercept) && 
                               !math.IsNaN(pm.StockCoeff) && 
                               !math.IsNaN(pm.SalesCoeff),
    }
}

// GetHistory returns the historical data for visualization
func (pm *PricingModel) GetHistory() map[string]interface{} {
    // Filter out any NaN values before returning
    stockLevels := make([]float64, 0, len(pm.StockLevels))
    salesRates := make([]float64, 0, len(pm.SalesRates))
    prices := make([]float64, 0, len(pm.Prices))
    
    for i := range pm.StockLevels {
        if !math.IsNaN(pm.StockLevels[i]) && !math.IsNaN(pm.SalesRates[i]) && !math.IsNaN(pm.Prices[i]) {
            stockLevels = append(stockLevels, pm.StockLevels[i])
            salesRates = append(salesRates, pm.SalesRates[i])
            prices = append(prices, pm.Prices[i])
        }
    }
    
    return map[string]interface{}{
        "stockLevels": stockLevels,
        "salesRates": salesRates,
        "prices": prices,
    }
}

// ResetModel resets the model to its initial state
func (pm *PricingModel) ResetModel(basePrice float64, initialStock int) {
    pm.Intercept = basePrice
    pm.StockCoeff = -0.01
    pm.SalesCoeff = 0.05
    pm.StockLevels = []float64{}
    pm.SalesRates = []float64{}
    pm.Prices = []float64{}
    pm.CurrentStock = initialStock
    pm.CurrentSalesRate = 0
    pm.InitialBasePrice = basePrice
}

func TestData() {
	regression.PowCross(5,2.3)
}