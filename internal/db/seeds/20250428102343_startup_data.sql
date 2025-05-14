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
INSERT INTO stocks (product_id, quantity, first_stocked_date) VALUES
-- Laptops
((SELECT id FROM products WHERE name = 'MacBook Air'), 30, NOW()),
((SELECT id FROM products WHERE name = 'MacBook Pro'), 25, NOW()),
((SELECT id FROM products WHERE name = 'Galaxy Book2'), 20, NOW()),
((SELECT id FROM products WHERE name = 'Galaxy Book3 Pro'), 15, NOW()),
((SELECT id FROM products WHERE name = 'Pixelbook Go'), 22, NOW());

--Phones
INSERT INTO stocks (product_id, quantity, first_stocked_date) VALUES
((SELECT id FROM products WHERE name = 'iPhone 15'), 75, NOW()),
((SELECT id FROM products WHERE name = 'iPhone SE'), 60, NOW()),
((SELECT id FROM products WHERE name = 'Galaxy S23'), 50, NOW()),
((SELECT id FROM products WHERE name = 'Galaxy A54'), 65, NOW()),
((SELECT id FROM products WHERE name = 'Pixel 8'), 45, NOW());

-- Watches
INSERT INTO stocks (product_id, quantity, first_stocked_date) VALUES
((SELECT id FROM products WHERE name = 'Apple Watch Series 9'), 40, NOW()),
((SELECT id FROM products WHERE name = 'Apple Watch SE'), 45, NOW()),
((SELECT id FROM products WHERE name = 'Galaxy Watch 6'), 35, NOW()),
((SELECT id FROM products WHERE name = 'Galaxy Watch Active 2'), 30, NOW()),
((SELECT id FROM products WHERE name = 'Pixel Watch 2'), 25, NOW());

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
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