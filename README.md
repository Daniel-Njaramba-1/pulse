Focus on electronics

- Price Volatility - electronics have frequent price fluctuations - supply chain issues, dynamic pricing is highly valuable
- Competitive pricing pressure
- Clear Feature Relevance - Demand Score, Inventory Ratio, Competitive Index
- High Perceived Value - Customers are price sensitive with electronics

Price = Intercept + (Feature1 * Coefficient1) + (Feature2 * Coefficient2) + ... + (FeatureN * CoefficientN)

Training Model - readjusting coefficients based on previous sales data and other metrics
In the TrainNewModel method of ModelService, the system:
- Collects historical pricing features from the database
- Gets historical product prices
- Creates a new pricing model instance
- Trains it using these features and prices

The training itself happens in the Train method of PricingModel where:

Each feature and corresponding price are fed into the regression model
The model learns the relationship between features and prices
After training, new coefficients are calculated and saved

golang trigger for adjustprice after sales outside transaction

Adjusted Price = Base Price × (1 + Σ(coefficient_i × feature_i))

So I should lean into transparency with customers:
informing them of upcoming price changes - price adjustment will run on a cron job
informing them of a product with good sentiment scores from reviews
informing them of a freshly restocked product
informing them of a product selling fast

Unique Angle - Transparency
Informing of price trends using price adjustment table - draw graph


Core Performance Metrics
1. Model Performance Tracking

R-squared trend over time - Track model accuracy across different training iterations
MSE/RMSE trends - Monitor prediction error patterns
MAE progression - Track absolute error improvements

2. Sales Analytics

Sales volume over time (daily/weekly/monthly)
Revenue trends by product/category/brand
Sales velocity by product - crucial for your pricing model
Top performing products (by volume and revenue)
Sales conversion funnel (cart → order → payment → sale)

3. Inventory Management

Current stock levels across all products
Stock turnover rates by product/category
Low stock alerts (products below threshold)
Restock frequency patterns
Days since last restock distribution

4. Dynamic Pricing Intelligence

Price adjustment frequency and magnitude
Base price vs adjusted price comparisons
Price elasticity impact (price changes vs sales response)
Revenue impact of price adjustments
Price adjustment distribution by product categories

Additional Strategic Dashboards
5. Customer Behavior

Review ratings distribution and trends
Wishlist vs purchase conversion rates
Customer acquisition and retention metrics
Average order value trends

6. Product Performance Matrix

Rating vs sales volume scatter plot
Price vs demand elasticity analysis
Category performance comparison
Brand performance metrics

7. Operational Health

Order status distribution (pending/completed/failed/cancelled)
Payment success rates by method
Stock-out frequency and duration
System performance metrics (model run frequency, data freshness)

Real-time Monitoring
8. Live Metrics

Current active carts and conversion probability
Real-time sales feed
Inventory alerts and critical stock levels
Recent price adjustments and their immediate impact













