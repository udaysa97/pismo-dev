BEGIN;
ALTER TABLE order_metadata 
    ALTER COLUMN "count" TYPE INT,
    ALTER COLUMN "count" SET NOT NULL;
COMMIT;