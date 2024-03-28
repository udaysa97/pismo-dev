BEGIN;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS order_metadata (
	"order_id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
	"user_id" uuid NOT NULL,
	"status" varchar NOT NULL,
	"network_id" uuid NOT NULL,
	"entity_type" varchar NOT NULL DEFAULT 'NFT',
	"entity_address" varchar NOT NULL,
	"nft_id" varchar NOT NULL,
	"count" int NOT NULL,
	"order_type" varchar NOT NULL,
	"slippage" varchar NOT NULL,
	"execution_response" text NULL,
	"created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	"updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX uidx_order_metadata ON order_metadata (order_id, user_id);

CREATE TABLE IF NOT EXISTS transaction_data (
	"order_id" uuid NOT NULL,
	"tx_hash" varchar NOT NULL,
	"status" varchar NOT NULL,
	"order_tx_type" varchar NOT NULL,
	"payload_type" varchar NOT NULL,
	"gas_used" int NOT NULL,
	"gas_price" int NOT NULL,
	"token_transfers" jsonb NULL,
	PRIMARY KEY (order_id, tx_hash)
);

CREATE UNIQUE INDEX uidx_transaction_data ON transaction_data (order_id);
COMMIT;