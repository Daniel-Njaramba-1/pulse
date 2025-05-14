-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION notify_price_adjustment()
RETURNS TRIGGER AS $$
BEGIN 
    PERFORM pg_notify('price_adjustment', row_to_json(NEW)::text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_notify_price_adjustment
AFTER INSERT OR UPDATE ON price_adjustments
FOR EACH ROW 
EXECUTE FUNCTION notify_price_adjustment();

CREATE OR REPLACE FUNCTION notify_sale()
RETURNS TRIGGER AS $$
BEGIN 
    PERFORM pg_notify('sale', row_to_json(NEW)::text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_notify_sale
AFTER INSERT OR UPDATE ON sales
FOR EACH ROW 
EXECUTE FUNCTION notify_sale();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trigger_notify_price_adjustment ON price_adjustments;
DROP FUNCTION IF EXISTS notify_price_adjustment();

DROP TRIGGER IF EXISTS trigger_notify_sale ON sales;
DROP FUNCTION IF EXISTS notify_sale();
-- +goose StatementEnd
