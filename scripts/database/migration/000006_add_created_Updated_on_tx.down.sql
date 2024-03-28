BEGIN;
ALTER TABLE transaction_data 
    DROP COLUMN IF EXISTS "created_at";
ALTER TABLE transaction_data 
    DROP COLUMN IF EXISTS "updated_at";
COMMIT;