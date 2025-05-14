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













