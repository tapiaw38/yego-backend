ALTER TABLE orders ADD COLUMN IF NOT EXISTS delivery_user_id VARCHAR(255);
ALTER TABLE orders ADD COLUMN IF NOT EXISTS delivery_accepted_at TIMESTAMP WITH TIME ZONE;

CREATE INDEX IF NOT EXISTS idx_orders_delivery_user_id ON orders(delivery_user_id);
