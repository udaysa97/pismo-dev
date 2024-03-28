BEGIN;
ALTER TABLE IF EXISTS order_metadata 
    ADD COLUMN IF NOT EXISTS retry_count INT DEFAULT 0;

DROP INDEX IF EXISTS idx_order_type;
DROP INDEX IF EXISTS idx_status;
CREATE INDEX idx_order_type ON order_metadata (order_type);
CREATE INDEX idx_status ON order_metadata (status);
COMMIT;