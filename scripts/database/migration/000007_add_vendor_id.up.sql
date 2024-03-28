BEGIN;
ALTER TABLE order_metadata 
    ADD COLUMN "vendor_id" varchar NULL;
COMMIT;