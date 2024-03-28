package query

type TransactionData struct {
	OrderId        string `json:"order_id"`
	TxHash         string `json:"tx_hash"`
	Status         string `json:"status"`
	OrderTxType    string `json:"order_tx_type"`
	PayloadType    string `json:"payload_type"`
	GasUsed        string `json:"gas_used"`
	GasPrice       string `json:"gas_price"`
	TokenTransfers string `json:"token_transfers"`
}
