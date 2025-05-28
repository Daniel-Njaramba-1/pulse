import logging
from datetime import datetime
from sqlalchemy import text

from db import get_db_engine
import nltk
from nltk.sentiment import SentimentIntensityAnalyzer

engine = get_db_engine()

# Download VADER lexicon if not already present
try:
    nltk.data.find('sentiment/vader_lexicon.zip')
except LookupError:
    nltk.download('vader_lexicon')

# Configure logging
logger = logging.getLogger(__name__)

def get_days_since_last_sale(product_id: int) -> int | None:
    """Calculate days since the last sale of a product."""
    logger.info(f"Called get_days_since_last_sale with product_id={product_id}")
    query = """
        SELECT last_sale
        FROM product_metrics
        WHERE product_id = :product_id
    """
    with engine.connect() as conn:
        result = conn.execute(text(query), {"product_id": product_id}).fetchone()
        if result and result[0]:
            last_sale = result[0]
            days_since = (datetime.now() - last_sale).days
            logger.info(f"get_days_since_last_sale result: {days_since}")
            return days_since
        logger.info("get_days_since_last_sale result: None")
        return None
    
def get_sales_velocity(product_id: int, days: int = 30) -> float:
    """Calculate sales velocity (sales per day) over the last N days."""
    logger.info(f"Called get_sales_velocity with product_id={product_id}, days={days}")
    query = """
        SELECT COUNT(*) as sale_count
        FROM sales
        WHERE product_id = :product_id
        AND created_at >= NOW() - INTERVAL ':days days'
    """
    with engine.connect() as conn:
        result = conn.execute(text(query), {"product_id": product_id, "days": days}).fetchone()
        velocity = result[0] / days if result and result[0] else 0.0
        logger.info(f"get_sales_velocity result: {velocity}")
        return velocity

def get_total_sales_count(product_id: int) -> int:
    """Get total number of sales for a product."""
    logger.info(f"Called get_total_sales_count with product_id={product_id}")
    query = """
        SELECT COUNT(*) as total_sales
        FROM sales
        WHERE product_id = :product_id
    """
    with engine.connect() as conn:
        result = conn.execute(text(query), {"product_id": product_id}).fetchone()
        count = result[0] if result else 0
        logger.info(f"get_total_sales_count result: {count}")
        return count

def get_total_sales_value(product_id: int) -> float:
    """Get total sales value for a product."""
    logger.info(f"Called get_total_sales_value with product_id={product_id}")
    query = """
        SELECT SUM(sale_price * quantity) as total_value
        FROM sales
        WHERE product_id = :product_id
    """
    with engine.connect() as conn:
        result = conn.execute(text(query), {"product_id": product_id}).fetchone()
        value = float(result[0]) if result and result[0] else 0.0
        logger.info(f"get_total_sales_value result: {value}")
        return value

def get_category_percentile(product_id: int) -> float:
    """Calculate product's sales percentile within its category."""
    logger.info(f"Called get_category_percentile with product_id={product_id}")
    query = """
        WITH category_sales AS (
            SELECT p.id, p.category_id, COALESCE(SUM(s.sale_price * s.quantity), 0) as total_sales
            FROM products p
            LEFT JOIN sales s ON p.id = s.product_id
            GROUP BY p.id, p.category_id
        ),
        category_rankings AS (
            SELECT 
                id,
                category_id,
                total_sales,
                PERCENT_RANK() OVER (PARTITION BY category_id ORDER BY total_sales) as sales_percentile
            FROM category_sales
        )
        SELECT sales_percentile
        FROM category_rankings
        WHERE id = :product_id
    """
    with engine.connect() as conn:
        result = conn.execute(text(query), {"product_id": product_id}).fetchone()
        percentile = float(result[0]) if result and result[0] is not None else 0.0
        logger.info(f"get_category_percentile result: {percentile}")
        return percentile

def get_review_score(product_id: int) -> float:
    """Get average review score for a product, enhanced with sentiment analysis."""

    logger.info(f"Called get_review_score with product_id={product_id}")

    # Get numeric average rating
    average_rating_query = """
        SELECT average_rating
        FROM product_metrics
        WHERE product_id = :product_id
    """

    # Get review texts
    review_query = """
        SELECT review_text
        FROM reviews
        WHERE product_id = :product_id
        AND review_text IS NOT NULL
        AND review_text != ''
    """

    with engine.connect() as conn:
        # Numeric average
        result = conn.execute(text(average_rating_query), {"product_id": product_id}).fetchone()
        numeric_score = float(result[0]) if result and result[0] is not None else 0.0

        # Sentiment analysis
        review_texts = [row[0] for row in conn.execute(text(review_query), {"product_id": product_id}).fetchall()]
        if review_texts:
            sia = SentimentIntensityAnalyzer()
            sentiment_scores = [sia.polarity_scores(text)["compound"] for text in review_texts]
            avg_sentiment = sum(sentiment_scores) / len(sentiment_scores)
            # Normalize sentiment (-1 to 1) to (0 to 5) scale
            sentiment_score = (avg_sentiment + 1) * 2.5
            # Weighted average: 70% numeric, 30% sentiment
            final_score = 0.7 * numeric_score + 0.3 * sentiment_score
            logger.info(f"get_review_score result (with sentiment): {final_score}")
            return final_score
        else:
            logger.info(f"get_review_score result (numeric only): {numeric_score}")
            return numeric_score

def get_wishlist_to_sales_ratio(product_id: int) -> float:
    """Calculate ratio of wishlist count to total sales."""
    logger.info(f"Called get_wishlist_to_sales_ratio with product_id={product_id}")
    query = """
        SELECT 
            pm.wishlist_count,
            COUNT(s.id) as total_sales
        FROM product_metrics pm
        LEFT JOIN sales s ON pm.product_id = s.product_id
        WHERE pm.product_id = :product_id
        GROUP BY pm.wishlist_count
    """
    with engine.connect() as conn:
        result = conn.execute(text(query), {"product_id": product_id}).fetchone()
        if result and result[0] and result[1]:
            ratio = float(result[0]) / float(result[1]) if float(result[1]) > 0 else 0.0
            logger.info(f"get_wishlist_to_sales_ratio result: {ratio}")
            return ratio
        logger.info("get_wishlist_to_sales_ratio result: 0.0")
        return 0.0

def get_days_since_restock(product_id: int) -> int | None:
    """Calculate days since the last restock of a product."""
    logger.info(f"Called get_days_since_restock with product_id={product_id}")
    query = """
        SELECT created_at
        FROM stock_history
        WHERE product_id = :product_id
        AND event_type = 'restock'
        ORDER BY created_at DESC
        LIMIT 1
    """
    with engine.connect() as conn:
        result = conn.execute(text(query), {"product_id": product_id}).fetchone()
        if result and result[0]:
            last_restock = result[0]
            days_since = (datetime.now() - last_restock).days
            logger.info(f"get_days_since_restock result: {days_since}")
            return days_since
        logger.info("get_days_since_restock result: None")
        return None

def compute_all_features(product_id: int) -> dict:
    """Compute all pricing features for a product."""
    logger.info(f"Called compute_all_features with product_id={product_id}")
    features = {
        "product_id": product_id,
        "days_since_last_sale": get_days_since_last_sale(product_id),
        "sales_velocity": get_sales_velocity(product_id),
        "total_sales_count": get_total_sales_count(product_id),
        "total_sales_value": get_total_sales_value(product_id),
        "category_percentile": get_category_percentile(product_id),
        "review_score": get_review_score(product_id),
        "wishlist_to_sales_ratio": get_wishlist_to_sales_ratio(product_id),
        "days_since_restock": get_days_since_restock(product_id)
    }
    logger.info(f"compute_all_features result: {features}")
    return features

def process_features(product_id: int) -> str:
    """Compute all pricing features and save to the pricing_features table."""
    logger.info(f"Called process_features with product_id={product_id}")
    features = compute_all_features(product_id)
    columns = ', '.join(features.keys())
    placeholders = ', '.join([f":{k}" for k in features.keys()])
    update_assignments = ', '.join([f"{k} = :{k}" for k in features.keys() if k != 'product_id'])

    select_query = "SELECT 1 FROM pricing_features WHERE product_id = :product_id"
    insert_query = f"""
        INSERT INTO pricing_features ({columns})
        VALUES ({placeholders})
    """
    update_query = f"""
        UPDATE pricing_features
        SET {update_assignments}
        WHERE product_id = :product_id
    """

    try:
        with engine.begin() as conn:
            exists = conn.execute(text(select_query), {"product_id": product_id}).fetchone()
            if exists:
                conn.execute(text(update_query), features)
                logger.info("process_features: updated existing record")
            else:
                conn.execute(text(insert_query), features)
                logger.info("process_features: inserted new record")
        return "success"
    except Exception as e:
        logger.error(f"process_features error: {e}")
        return f"error: {e}"

