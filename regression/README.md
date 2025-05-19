Backend - Dynamic Pricing - Multi Variable Linear Regression
Golang + Python

Golang API request → Python modeling service → Model training → Coefficient storage → Price adjustment → Response to API

Adjusted Price = Base Price + (Σ(coefficient_i × feature_i))

1. Build Pricing Features
2. Collect Training Data
3. Normalize
   Fit Initial Model
4. Evaluate Model Performance
5. (Optional) Refine Model Coefficients
6. Apply Model to Adjust Prices
7. Store Results


Model Refinement:
Start with basic linear regression
Evaluate using R² and MSE
If below threshold, try regularization techniques (Ridge/Lasso)
Consider polynomial features for important variables
Save the best model automatically

Initially coefficients are 0
Adjusted Price == Base Price

Normalization




Formula 1: Multiplicative Model
Adjusted Price = Base Price × (1 + Σ(coefficient_i × feature_i))

Formula 2: Additive Model
Adjusted Price = Base Price + (Σ(coefficient_i × feature_i))

1. Scaling Behavior
Multiplicative Model:

Price changes scale proportionally with the base price
A 10% increase applies equally to both low and high-priced items
Features produce percentage adjustments to the base price

Additive Model:

Price changes are absolute dollar amounts regardless of base price
A $10 increase has different proportional impact on low vs high-priced items
Features produce fixed dollar adjustments

2. Impact on Different Price Points
For example, with a +$10 coefficient impact:

Multiplicative: $100 item → +10% → $110; $1000 item → +10% → $1100
Additive: $100 item → +$10 → $110; $1000 item → +$10 → $1010

3. Coefficient Interpretation
Multiplicative:

Coefficients represent percentage adjustments
A coefficient of 0.05 means "increase price by 5% per unit of feature"

Additive:

Coefficients represent absolute dollar adjustments
A coefficient of 5 means "add $5 per unit of feature"

# Update the adjusted price in the database
    update_query = """
    UPDATE product_metrics 
    SET 
        adjusted_price = :adjusted_price,
        price_adjustment_ratio = :ratio,
        price_adjustment_bounded = :bounded
    WHERE product_id = :product_id
    """
