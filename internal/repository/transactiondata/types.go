package transactiondata

import "github.com/google/uuid"

type TransactionRequest struct {
	OrderId        uuid.UUID      `json:"order_id"`
	TxHash         string         `json:"tx_hash"`
	Status         string         `json:"status"`
	OrderTxType    string         `json:"order_tx_type"`
	PayloadType    string         `json:"payload_type"`
	GasUsed        int32          `json:"gas_used"`
	GasPrice       int32          `json:"gas_price"`
	TokenTransfers map[string]any `json:"token_transfers"`
}
