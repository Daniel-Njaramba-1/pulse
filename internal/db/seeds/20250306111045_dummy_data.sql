INSERT INTO price_model_coefficients (
            model_version, training_date, sample_size, r_squared,
            intercept, sales_count_coef, sales_value_coef, sales_velocity_coef,
            days_since_sale_coef, category_rank_coef, category_percentile_coef,
            review_score_coef, wishlist_ratio_coef, days_in_stock_coef, seasonal_factor_coef
        ) VALUES (
            1, :training_date, :sample_size, :r_squared,
            :intercept, :sales_count_coef, :sales_value_coef, :sales_velocity_coef,
            :days_since_sale_coef, :category_rank_coef, :category_percentile_coef,
            :review_score_coef, :wishlist_ratio_coef, :days_in_stock_coef, :seasonal_factor_coef
        )