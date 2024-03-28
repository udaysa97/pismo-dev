package types

import "pismo-dev/constants"

type OGMintQueueMessage struct {
	RetryCount         int                `json:"retry_count"`
	Status             string             `json:"status"`
	UserId             string             `json:"user_id"`
	OrderId            string             `json:"order_id"`
	NetworkId          string             `json:"network_id"`
	TxHash             string             `json:"tx_hash"`
	TxIdentifier       string             `json:"tx_identifier"`
	Type               constants.OGOPTYPE `json:"type"`
	ToAddress          string             `json:"to_address"`
	MetaDataURI        string             `json:"metadata_uri"`
	ContractIdentifier string             `json:"contract_identifier"`
	ErrorMessage       string             `json:"error_message"`
	OrderType          string             `json:"order_type"`
	VendorName         string             `json:"vendor_name"`
}
