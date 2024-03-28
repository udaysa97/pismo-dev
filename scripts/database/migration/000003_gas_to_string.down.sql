BEGIN;
ALTER TABLE transaction_data 
    ALTER COLUMN "gas_used" TYPE INT,
    ALTER COLUMN "gas_used" SET NOT NULL,
    ALTER COLUMN "gas_price" TYPE INT,
    ALTER COLUMN "gas_price" SET NOT NULL;
COMMIT;