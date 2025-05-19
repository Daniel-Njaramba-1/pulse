import logging
from dotenv import load_dotenv
from flask import Flask, request, jsonify
from flask_cors import CORS

from features import process_features
from learning import train_model, adjust_price_for_all_products, adjust_price_for_product
# Load environment variables
load_dotenv()

# Configure logging
logging.basicConfig(
    format='%(asctime)s-%(message)s'
)
logger = logging.getLogger(__name__)


# Flask app setup
app = Flask(__name__)
CORS(app)

@app.route('/compute_features', methods=['POST'])
def compute_features():
    data = request.get_json()
    product_id = data.get('product_id')
    if product_id is None:
        return jsonify({"error": "Missing product_id"}), 400
    result = process_features(product_id)
    return jsonify({"result": result})

@app.route('/train_model', methods=['POST'])
def train_pricing_model():
    """Train pricing model with available data"""
    try:
        logger.info("Starting model training")
        model_coefficients = train_model()
        return jsonify({
            "status": "success",
            "message": "Model training completed successfully",
            "model_version": model_coefficients.model_version,
            "metrics": {
                "r_squared": model_coefficients.r_squared,
                "rmse": model_coefficients.rmse,
                "mae": model_coefficients.mae,
                "sample_size": model_coefficients.sample_size
            }
        }), 200
    except Exception as e:
        logger.exception(f"Error training model: {str(e)}")
        return jsonify({"status": "error", "message": str(e)}), 500

# Adjust price for a single product
@app.route('/adjust-price/<int:product_id>', methods=['POST'])
def adjust_product_price(product_id):
    """Adjust price for one product"""
    try:
        logger.info(f"Starting price adjustment for product {product_id}")
        result = adjust_price_for_product(product_id)
        return jsonify({"result": result}), 200
    except Exception as e:
        logger.exception(f"Error adjusting price for product {product_id}: {str(e)}")
        return jsonify({"error": str(e)}), 500

# Adjust prices for all products
@app.route('/adjust-prices', methods=['POST'])
def adjust_all_prices():
    """Adjust all prices"""
    try:
        logger.info("Starting adjusting all prices")
        result = adjust_price_for_all_products()
        return jsonify({"result": result}), 200
    except Exception as e:
        logger.exception(f"Error adjusting all prices: {str(e)}")
        return jsonify({"error": str(e)}), 500

# Set up file handler for logging if not already present
if not any(isinstance(h, logging.FileHandler) for h in logger.handlers):
    file_handler = logging.FileHandler("feature_functions.log")
    file_handler.setFormatter(logging.Formatter('%(asctime)s-%(levelname)s-%(message)s'))
    logger.addHandler(file_handler)
    logger.setLevel(logging.INFO)



if __name__ == '__main__':
    app.run(debug=True, port=5872)
