BEGIN;
ALTER TABLE IF EXISTS order_metadata 
    DROP COLUMN IF EXISTS retry_count;

DROP INDEX IF EXISTS idx_order_type;
DROP INDEX IF EXISTS idx_status;
COMMIT;