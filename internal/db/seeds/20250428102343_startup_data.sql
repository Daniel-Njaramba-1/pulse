-- +goose Up
-- +goose StatementBegin

-- Mock data for brands table
INSERT INTO brands (name, description) VALUES
('Samsung', 'South Korean multinational conglomerate.'),
('Apple', 'American multinational technology company.'),
('Google', 'Leading American company on Android and Browsers.');

-- Seed categories
INSERT INTO categories (name, description) VALUES
('Laptops', 'Portable personal computers.'),
('Phones', 'Mobile phones with advanced computing capabilities.'),
('Watches', 'Wearable devices with smart features.');

-- Seed products (15 total: 5 per category)
INSERT INTO products (category_id, brand_id, name, description, is_active) VALUES
-- Laptops
((SELECT id FROM categories WHERE name = 'Laptops'), (SELECT id FROM brands WHERE name = 'Apple'), 'MacBook Air', 'Lightweight and efficient.', TRUE),
((SELECT id FROM categories WHERE name = 'Laptops'), (SELECT id FROM brands WHERE name = 'Apple'), 'MacBook Pro', 'High-performance M-series laptop.', TRUE),
((SELECT id FROM categories WHERE name = 'Laptops'), (SELECT id FROM brands WHERE name = 'Samsung'), 'Galaxy Book2', 'ultraportable Samsung laptop.', TRUE),
((SELECT id FROM categories WHERE name = 'Laptops'), (SELECT id FROM brands WHERE name = 'Samsung'), 'Galaxy Book3 Pro', 'High-end Samsung laptop.', TRUE),
((SELECT id FROM categories WHERE name = 'Laptops'), (SELECT id FROM brands WHERE name = 'Google'), 'Pixelbook Go', 'Google Chromebook.', TRUE);

-- Phones
INSERT INTO products (category_id, brand_id, name, description, is_active) VALUES
((SELECT id FROM categories WHERE name = 'Phones'), (SELECT id FROM brands WHERE name = 'Apple'), 'iPhone 15', 'Latest iPhone from Apple.', TRUE),
((SELECT id FROM categories WHERE name = 'Phones'), (SELECT id FROM brands WHERE name = 'Apple'), 'iPhone SE', 'Compact and affordable iPhone.', TRUE),
((SELECT id FROM categories WHERE name = 'Phones'), (SELECT id FROM brands WHERE name = 'Samsung'), 'Galaxy S23', 'Flagship Samsung smartphone.', TRUE),
((SELECT id FROM categories WHERE name = 'Phones'), (SELECT id FROM brands WHERE name = 'Samsung'), 'Galaxy A54', 'Mid-range Samsung smartphone.', TRUE),
((SELECT id FROM categories WHERE name = 'Phones'), (SELECT id FROM brands WHERE name = 'Google'), 'Pixel 8', 'Android phone with AI features.', TRUE);

-- Watches
INSERT INTO products (category_id, brand_id, name, description, is_active) VALUES
((SELECT id FROM categories WHERE name = 'Watches'), (SELECT id FROM brands WHERE name = 'Apple'), 'Apple Watch Series 9', 'Latest Apple Watch.', TRUE),
((SELECT id FROM categories WHERE name = 'Watches'), (SELECT id FROM brands WHERE name = 'Apple'), 'Apple Watch SE', 'Affordable smartwatch option.', TRUE),
((SELECT id FROM categories WHERE name = 'Watches'), (SELECT id FROM brands WHERE name = 'Samsung'), 'Galaxy Watch 6', 'Latest Samsung smartwatch.', TRUE),
((SELECT id FROM categories WHERE name = 'Watches'), (SELECT id FROM brands WHERE name = 'Samsung'), 'Galaxy Watch Active 2', 'Fitness-focused smartwatch.', TRUE),
((SELECT id FROM categories WHERE name = 'Watches'), (SELECT id FROM brands WHERE name = 'Google'), 'Pixel Watch 2', 'Smart Google wearable.', TRUE);

-- Seed data for product_metrics table - creating metrics for all products
INSERT INTO product_metrics (product_id, base_price, adjusted_price) VALUES
-- Laptops
((SELECT id FROM products WHERE name = 'MacBook Air'), 1199.99, 1199.99),
((SELECT id FROM products WHERE name = 'MacBook Pro'), 1999.99, 1199.99),
((SELECT id FROM products WHERE name = 'Galaxy Book2'), 899.99, 899.99),
((SELECT id FROM products WHERE name = 'Galaxy Book3 Pro'), 1399.99, 1399.99),
((SELECT id FROM products WHERE name = 'Pixelbook Go'), 849.99, 849.99);

-- Phones
INSERT INTO product_metrics (product_id, base_price, adjusted_price) VALUES
((SELECT id FROM products WHERE name = 'iPhone 15'), 1099.99, 1099.99),
((SELECT id FROM products WHERE name = 'iPhone SE'), 429.99, 429.99),
((SELECT id FROM products WHERE name = 'Galaxy S23'), 999.99, 999.99),
((SELECT id FROM products WHERE name = 'Galaxy A54'), 449.99, 449.99),
((SELECT id FROM products WHERE name = 'Pixel 8'), 699.99, 699.99);

-- Watches
INSERT INTO product_metrics (product_id, base_price, adjusted_price) VALUES
((SELECT id FROM products WHERE name = 'Apple Watch Series 9'), 399.99, 399.99),
((SELECT id FROM products WHERE name = 'Apple Watch SE'), 249.99, 249.99),
((SELECT id FROM products WHERE name = 'Galaxy Watch 6'), 329.99, 329.99),
((SELECT id FROM products WHERE name = 'Galaxy Watch Active 2'), 249.99, 249.99),
((SELECT id FROM products WHERE name = 'Pixel Watch 2'), 349.99, 349.99);

-- Seed data for stocks table
INSERT INTO stocks (product_id, quantity) VALUES
-- Laptops
((SELECT id FROM products WHERE name = 'MacBook Air'), 30),
((SELECT id FROM products WHERE name = 'MacBook Pro'), 25),
((SELECT id FROM products WHERE name = 'Galaxy Book2'), 20),
((SELECT id FROM products WHERE name = 'Galaxy Book3 Pro'), 15),
((SELECT id FROM products WHERE name = 'Pixelbook Go'), 22);

--Phones
INSERT INTO stocks (product_id, quantity) VALUES
((SELECT id FROM products WHERE name = 'iPhone 15'), 75),
((SELECT id FROM products WHERE name = 'iPhone SE'), 60),
((SELECT id FROM products WHERE name = 'Galaxy S23'), 50),
((SELECT id FROM products WHERE name = 'Galaxy A54'), 65),
((SELECT id FROM products WHERE name = 'Pixel 8'), 45);

-- Watches
INSERT INTO stocks (product_id, quantity) VALUES
((SELECT id FROM products WHERE name = 'Apple Watch Series 9'), 40),
((SELECT id FROM products WHERE name = 'Apple Watch SE'), 45),
((SELECT id FROM products WHERE name = 'Galaxy Watch 6'), 35),
((SELECT id FROM products WHERE name = 'Galaxy Watch Active 2'), 30),
((SELECT id FROM products WHERE name = 'Pixel Watch 2'), 25);

-- Seed data for stocking history
INSERT INTO stock_history (product_id, event_type, quantity_change, quantity_after) VALUES 
-- Laptops
((SELECT id FROM products WHERE name = 'MacBook Air'), 'restock', 30, 30),
((SELECT id FROM products WHERE name = 'MacBook Pro'), 'restock', 25, 25),
((SELECT id FROM products WHERE name = 'Galaxy Book2'), 'restock', 20, 20),
((SELECT id FROM products WHERE name = 'Galaxy Book3 Pro'),'restock', 15, 15),
((SELECT id FROM products WHERE name = 'Pixelbook Go'), 'restock', 22, 22);

--Phones
INSERT INTO stock_history (product_id, event_type, quantity_change, quantity_after) VALUES
((SELECT id FROM products WHERE name = 'iPhone 15'), 'restock', 75, 75),
((SELECT id FROM products WHERE name = 'iPhone SE'), 'restock', 60, 60),
((SELECT id FROM products WHERE name = 'Galaxy S23'), 'restock', 50, 50),
((SELECT id FROM products WHERE name = 'Galaxy A54'), 'restock', 65, 65),
((SELECT id FROM products WHERE name = 'Pixel 8'), 'restock', 45, 45);

-- Watches
INSERT INTO stock_history (product_id, event_type, quantity_change, quantity_after) VALUES
((SELECT id FROM products WHERE name = 'Apple Watch Series 9'), 'restock', 40, 40),
((SELECT id FROM products WHERE name = 'Apple Watch SE'), 'restock', 45, 45),
((SELECT id FROM products WHERE name = 'Galaxy Watch 6'), 'restock', 35, 35),
((SELECT id FROM products WHERE name = 'Galaxy Watch Active 2'), 'restock', 30, 30),
((SELECT id FROM products WHERE name = 'Pixel Watch 2'), 'restock', 25, 25);


-- Seed data for price_model_coefficients
INSERT INTO price_model_coefficients (
  model_version,
  training_date,
  sample_size,
  days_since_last_sale_coef,
  sales_velocity_coef,
  total_sales_count_coef,
  total_sales_value_coef,
  category_percentile_coef,
  review_score_coef,
  wishlist_to_sales_ratio_coef,
  days_since_restock_coef
) VALUES
('v1.0', NOW(), 0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM price_model_coefficients WHERE model_version = 'v1.0';

DELETE FROM stock_history WHERE product_id IN (
  SELECT id FROM products WHERE name IN (
    'MacBook Air', 'MacBook Pro', 'Galaxy Book2', 'Galaxy Book3 Pro', 'Pixelbook Go',
    'iPhone 15', 'iPhone SE', 'Galaxy S23', 'Galaxy A54', 'Pixel 8',
    'Apple Watch Series 9', 'Apple Watch SE', 'Galaxy Watch 6', 'Galaxy Watch Active 2', 'Pixel Watch 2'
  )
);

DELETE FROM stocks WHERE product_id IN (
  SELECT id FROM products WHERE name IN (
    'MacBook Air', 'MacBook Pro', 'Galaxy Book2', 'Galaxy Book3 Pro', 'Pixelbook Go',
    'iPhone 15', 'iPhone SE', 'Galaxy S23', 'Galaxy A54', 'Pixel 8',
    'Apple Watch Series 9', 'Apple Watch SE', 'Galaxy Watch 6', 'Galaxy Watch Active 2', 'Pixel Watch 2'
  )
);
DELETE FROM product_metrics WHERE product_id IN (
  SELECT id FROM products WHERE name IN (
    'MacBook Air', 'MacBook Pro', 'Galaxy Book2', 'Galaxy Book3 Pro', 'Pixelbook Go',
    'iPhone 15', 'iPhone SE', 'Galaxy S23', 'Galaxy A54', 'Pixel 8',
    'Apple Watch Series 9', 'Apple Watch SE', 'Galaxy Watch 6', 'Galaxy Watch Active 2', 'Pixel Watch 2'
  )
);
DELETE FROM products WHERE name IN (
  'MacBook Air', 'MacBook Pro', 'Galaxy Book2', 'Galaxy Book3 Pro', 'Pixelbook Go',
  'iPhone 15', 'iPhone SE', 'Galaxy S23', 'Galaxy A54', 'Pixel 8',
  'Apple Watch Series 9', 'Apple Watch SE', 'Galaxy Watch 6', 'Galaxy Watch Active 2', 'Pixel Watch 2'
);

DELETE FROM categories WHERE name IN ('Laptops', 'Phones', 'Watches');
DELETE FROM brands WHERE name IN ('Samsung', 'Apple', 'Google');
-- +goose StatementEnd