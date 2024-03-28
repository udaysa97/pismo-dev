package models

type TransactionData struct {
	OrderId        string      `json:"order_id"`
	TxHash         string      `json:"tx_hash"`
	Status         string      `json:"status"` //TransactionState
	OrderTxType    string      `json:"order_tx_type"`
	PayloadType    string      `json:"payload_type"`
	GasUsed        string      `json:"gas_used"`
	GasPrice       string      `json:"gas_price"`
	TokenTransfers interface{} `json:"token_transfers,omitempty" gorm:"type:jsonb;default:'[]';not null"`
}

type TransactionDetails struct {
	Count              int64             `json:"count"`
	TransactionDetails []TransactionData `json:"transaction_details"`
}
