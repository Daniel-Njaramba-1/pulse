from dataclasses import dataclass
from datetime import datetime
import pandas as pd
import numpy as np
from sklearn.linear_model import LinearRegression
from sklearn.metrics import r2_score, mean_squared_error, mean_absolute_error
from sqlalchemy import text

from db import get_db_engine
import logging

engine = get_db_engine()

@dataclass
class PricingFeatures:
    product_id: int
    days_since_last_sale: int 
    sales_velocity: float 
    total_sales_count: int 
    total_sales_value: float 
    category_percentile: float 
    review_score: float 
    wishlist_to_sales_ratio: float 
    days_since_restock: int
    last_model_run: datetime

@dataclass
class ModelCoefficients: 
    model_version: str 
    training_date: datetime
    sample_size: int
    r_squared: float 
    mse: float 
    rmse: float 
    mae: float 
    days_since_last_sale_coef: float
    sales_velocity_coef: float
    total_sales_count_coef: float 
    total_sales_value_coef: float 
    category_percentile_coef: float
    review_score_coef: float
    wishlist_to_sales_ratio_coef: float
    days_since_restock_coef: float

# Setup logging to file
logging.basicConfig(
    filename='model_training.log',
    level=logging.INFO,
    format='%(asctime)s %(levelname)s %(message)s'
)

# Model training and persistence

def get_all_training_data():
    """Fetch pricing features for all products and their current adjusted_price."""
    query = """
    SELECT 
        pf.*,
        pm.adjusted_price, pm.base_price
    FROM pricing_features pf
    JOIN product_metrics pm ON pf.product_id = pm.product_id
    WHERE pm.adjusted_price IS NOT NULL
    """
    
    df = pd.read_sql(query, engine)
    logging.info("Fetched training data from DB:\n%s", df.to_string())
    feature_columns = [
        'days_since_last_sale',
        'sales_velocity',
        'total_sales_count',
        'total_sales_value',
        'category_percentile',
        'review_score',
        'wishlist_to_sales_ratio',
        'days_since_restock'
    ]
    X = df[feature_columns]
    y = df['adjusted_price'] / df['base_price']
    logging.info("Input features (X):\n%s", X.to_string())
    logging.info("Target variable (y):\n%s", y.to_string())
    return X, y

def get_model_coefficients():
    """Get the most recent model coefficients from the database."""
    query = """
    SELECT * FROM price_model_coefficients
    ORDER BY training_date DESC
    LIMIT 1
    """
    df = pd.read_sql(query, engine)
    logging.info("Fetched model coefficients from DB:\n%s", df.to_string())
    if df.empty:
        return None
    row = df.iloc[0]
    return ModelCoefficients(
        model_version=row['model_version'],
        training_date=row['training_date'],
        sample_size=row['sample_size'],
        r_squared=row['r_squared'],
        mse=row['mse'],
        rmse=row['rmse'],
        mae=row['mae'],
        days_since_last_sale_coef=row['days_since_last_sale_coef'],
        sales_velocity_coef=row['sales_velocity_coef'],
        total_sales_count_coef=row['total_sales_count_coef'],
        total_sales_value_coef=row['total_sales_value_coef'],
        category_percentile_coef=row['category_percentile_coef'],
        review_score_coef=row['review_score_coef'],
        wishlist_to_sales_ratio_coef=row['wishlist_to_sales_ratio_coef'],
        days_since_restock_coef=row['days_since_restock_coef']
    )

def save_model_coefficients(coefficients: ModelCoefficients):
    """Save new model coefficients to the database."""
    query = """
    INSERT INTO price_model_coefficients (
        model_version, training_date, sample_size, r_squared, mse, rmse, mae,
        days_since_last_sale_coef, sales_velocity_coef, total_sales_count_coef,
        total_sales_value_coef, category_percentile_coef, review_score_coef,
        wishlist_to_sales_ratio_coef, days_since_restock_coef
    ) VALUES (
        :model_version, :training_date, :sample_size, :r_squared,
        :mse, :rmse, :mae, :days_since_last_sale_coef,
        :sales_velocity_coef, :total_sales_count_coef,
        :total_sales_value_coef, :category_percentile_coef,
        :review_score_coef, :wishlist_to_sales_ratio_coef,
        :days_since_restock_coef
    )
    """
    with engine.connect() as conn:
        conn.execute(text(query), {
            'model_version': coefficients.model_version,
            'training_date': coefficients.training_date,
            'sample_size': coefficients.sample_size,
            'r_squared': coefficients.r_squared,
            'mse': coefficients.mse,
            'rmse': coefficients.rmse,
            'mae': coefficients.mae,
            'days_since_last_sale_coef': coefficients.days_since_last_sale_coef,
            'sales_velocity_coef': coefficients.sales_velocity_coef,
            'total_sales_count_coef': coefficients.total_sales_count_coef,
            'total_sales_value_coef': coefficients.total_sales_value_coef,
            'category_percentile_coef': coefficients.category_percentile_coef,
            'review_score_coef': coefficients.review_score_coef,
            'wishlist_to_sales_ratio_coef': coefficients.wishlist_to_sales_ratio_coef,
            'days_since_restock_coef': coefficients.days_since_restock_coef
        })
        conn.commit()
    logging.info("Saved model coefficients to DB:\n%s", coefficients)

def train_model(): 
    """Train a linear regression model using the pricing features."""
    X, y = get_all_training_data()
    if len(X) == 0:
        raise ValueError("No training data available")
    model = LinearRegression()
    model.fit(X, y)
    y_pred = model.predict(X)
    logging.info("Model predictions (y_pred):\n%s", pd.Series(y_pred).to_string())
    r_squared = r2_score(y, y_pred)
    mse = mean_squared_error(y, y_pred)
    rmse = np.sqrt(mse)
    mae = mean_absolute_error(y, y_pred)
    new_coefficients = ModelCoefficients(
        model_version=f"v{datetime.now().strftime('%Y%m%d_%H%M%S')}",
        training_date=datetime.now(),
        sample_size=int(len(X)),
        r_squared=float(r_squared),
        mse=float(mse),
        rmse=float(rmse),
        mae=float(mae),
        days_since_last_sale_coef=float(model.coef_[0]),
        sales_velocity_coef=float(model.coef_[1]),
        total_sales_count_coef=float(model.coef_[2]),
        total_sales_value_coef=float(model.coef_[3]),
        category_percentile_coef=float(model.coef_[4]),
        review_score_coef=float(model.coef_[5]),
        wishlist_to_sales_ratio_coef=float(model.coef_[6]),
        days_since_restock_coef=float(model.coef_[7])
    )
    logging.info("Trained model coefficients:\n%s", new_coefficients)
    save_model_coefficients(new_coefficients)
    return new_coefficients

def adjust_price_for_product(product_id, min_ratio=0.8, max_ratio=1.2):
    """Adjust price for a single product using the latest model coefficients with bounds.
    
    Args:
        product_id: The ID of the product to adjust
        min/max adjusted price = +-20% base price
    
    Returns:
        The adjusted price
    """
    # Get latest model coefficients
    coefficients = get_model_coefficients()
    if coefficients is None:
        raise ValueError("No model coefficients found in database")
    
    # Get product features and base price
    query = """
    SELECT 
        pf.*,
        pm.base_price,
        pm.adjusted_price
    FROM pricing_features pf
    JOIN product_metrics pm ON pf.product_id = pm.product_id
    WHERE pf.product_id = %(product_id)s
    """
    df = pd.read_sql(query, engine, params={'product_id': product_id})
    
    if df.empty:
        raise ValueError(f"No data found for product {product_id}")
    
    # Prepare features in the same order as training
    features = [
        df['days_since_last_sale'].iloc[0],
        df['sales_velocity'].iloc[0],
        df['total_sales_count'].iloc[0],
        df['total_sales_value'].iloc[0],
        df['category_percentile'].iloc[0],
        df['review_score'].iloc[0],
        df['wishlist_to_sales_ratio'].iloc[0],
        df['days_since_restock'].iloc[0]
    ]
    
    # Calculate price adjustment ratio using model coefficients
    adjustment_ratio = (
        coefficients.days_since_last_sale_coef * features[0] +
        coefficients.sales_velocity_coef * features[1] +
        coefficients.total_sales_count_coef * features[2] +
        coefficients.total_sales_value_coef * features[3] +
        coefficients.category_percentile_coef * features[4] +
        coefficients.review_score_coef * features[5] +
        coefficients.wishlist_to_sales_ratio_coef * features[6] +
        coefficients.days_since_restock_coef * features[7]
    )
    
    # Bound the ratio to +-20% of base price
    bounded_ratio = min(max(adjustment_ratio, min_ratio), max_ratio)
    base_price = df['base_price'].iloc[0]
    old_adjusted_price = df['adjusted_price'].iloc[0] if not pd.isnull(df['adjusted_price'].iloc[0]) else None
    adjusted_price = base_price * bounded_ratio

    # Convert to native Python float to avoid PostgreSQL schema interpretation issues
    adjusted_price = float(adjusted_price)
    old_adjusted_price = float(old_adjusted_price) if old_adjusted_price is not None else None

    # Update the adjusted price in the database
    update_query = """
    UPDATE product_metrics 
    SET adjusted_price = :adjusted_price
    WHERE product_id = :product_id
    """

    # Log the price adjustment in a separate table
    insert_log_query = """
    INSERT INTO price_adjustments (product_id, old_price, new_price, model_version)
    VALUES (:product_id, :old_price, :new_price, :model_version)
    """

    with engine.connect() as conn:
        conn.execute(text(update_query), {
            'adjusted_price': adjusted_price,
            'product_id': product_id
        })
        conn.execute(text(insert_log_query), {
            'product_id': product_id,
            'old_price': old_adjusted_price,
            'new_price': adjusted_price,
            'model_version': coefficients.model_version
        })
        conn.commit()
    
    logging.info(f"Adjusted price for product {product_id}: {adjusted_price:.2f} (old: {old_adjusted_price}, ratio: {bounded_ratio:.4f})")
    return adjusted_price

def adjust_price_for_all_products(min_ratio=0.8, max_ratio=1.2):
    """Adjust price for all products using the latest model coefficients with bounds.

    Args:
        min/max adjusted price = +-20% base price

    Returns:
        DataFrame with product_id, adjusted_price
    """
    # Get latest model coefficients
    coefficients = get_model_coefficients()
    if coefficients is None:
        raise ValueError("No model coefficients found in database")

    # Get all products' features and base prices
    query = """
    SELECT 
        pf.*,
        pm.base_price,
        pm.adjusted_price
    FROM pricing_features pf
    JOIN product_metrics pm ON pf.product_id = pm.product_id
    """
    df = pd.read_sql(query, engine)

    if df.empty:
        raise ValueError("No product data found")

    # Calculate adjustment ratios for all products
    adjustment_ratios = (
        coefficients.days_since_last_sale_coef * df['days_since_last_sale'] +
        coefficients.sales_velocity_coef * df['sales_velocity'] +
        coefficients.total_sales_count_coef * df['total_sales_count'] +
        coefficients.total_sales_value_coef * df['total_sales_value'] +
        coefficients.category_percentile_coef * df['category_percentile'] +
        coefficients.review_score_coef * df['review_score'] +
        coefficients.wishlist_to_sales_ratio_coef * df['wishlist_to_sales_ratio'] +
        coefficients.days_since_restock_coef * df['days_since_restock']
    )

    # Bound the ratio to +-20% of base price
    bounded_ratios = adjustment_ratios.clip(lower=min_ratio, upper=max_ratio)
    df['bounded_ratio'] = bounded_ratios
    df['adjusted_price'] = df['base_price'] * df['bounded_ratio']

    # Track which products were bounded
    df['price_bounded'] = (adjustment_ratios != bounded_ratios)

    # Update adjusted prices and log adjustments in the database
    update_query = """
    UPDATE product_metrics 
    SET adjusted_price = :adjusted_price
    WHERE product_id = :product_id
    """

    insert_log_query = """
    INSERT INTO price_adjustments (product_id, old_price, new_price, model_version)
    VALUES (:product_id, :old_price, :new_price, :model_version)
    """

    with engine.connect() as conn:
        for _, row in df.iterrows():
            old_adjusted_price = row['adjusted_price'] if not pd.isnull(row['adjusted_price']) else None
            new_adjusted_price = row['adjusted_price']
            conn.execute(text(update_query), {
                'adjusted_price': new_adjusted_price,
                'product_id': row['product_id']
            })
            conn.execute(text(insert_log_query), {
                'product_id': row['product_id'],
                'old_price': old_adjusted_price,
                'new_price': new_adjusted_price,
                'model_version': coefficients.model_version
            })
        conn.commit()

    bounded_count = df['price_bounded'].sum()
    if bounded_count > 0:
        logging.info(f"{bounded_count} products had price adjustments bounded ({bounded_count/len(df)*100:.1f}%)")

    logging.info(f"Adjusted prices for {len(df)} products")
    return df[['product_id', 'adjusted_price']].to_dict(orient='records')
