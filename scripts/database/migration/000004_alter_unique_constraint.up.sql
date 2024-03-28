BEGIN;
DROP INDEX IF EXISTS uidx_transaction_data;
CREATE UNIQUE INDEX uidx_transaction_data ON transaction_data (order_id, tx_hash);
COMMIT;