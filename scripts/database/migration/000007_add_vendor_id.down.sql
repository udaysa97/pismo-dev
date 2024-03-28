BEGIN;
ALTER TABLE order_metadata 
    DROP COLUMN IF EXISTS "vendor_id";
COMMIT;