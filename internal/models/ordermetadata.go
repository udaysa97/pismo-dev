package models

import (
	"time"
)

type OrderMetadata struct {
	OrderId           string    `json:"order_id" gorm:"type:uuid;default:uuid_generate_v4()"`
	UserId            string    `json:"user_id"`
	VendorId          string    `json:"vendor_id"`
	Status            string    `json:"status"`
	NetworkId         string    `json:"network_id"`
	EntityType        string    `json:"entity_type"`
	EntityAddress     string    `json:"entity_address"`
	NftId             string    `json:"nft_id"`
	Count             string    `json:"count"`
	OrderType         string    `json:"order_type"`
	Slippage          string    `json:"slippage"`
	ExecutionResponse string    `json:"execution_response"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	RetryCount        int       `json:"retry_count"`
}

type OrderMetadataWithTx struct {
	OrderId           string    `json:"order_id" gorm:"type:uuid;default:uuid_generate_v4()"`
	UserId            string    `json:"user_id"`
	VendorId          string    `json:"vendor_id"`
	Status            string    `json:"status"`
	NetworkId         string    `json:"network_id"`
	EntityType        string    `json:"entity_type"`
	EntityAddress     string    `json:"entity_address"`
	NftId             string    `json:"nft_id"`
	Count             string    `json:"count"`
	OrderType         string    `json:"order_type"`
	Slippage          string    `json:"slippage"`
	ExecutionResponse string    `json:"execution_response"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	RetryCount        int       `json:"retry_count"`
	TxHash            string    `json:"tx_hash"`
}

type OrderDetails struct {
	Count        int64                 `json:"count"`
	OrderDetails []OrderMetadataWithTx `json:"order_details"`
}

type Pagination struct {
	Limit     int    `json:"limit"`
	Page      int    `json:"page"`
	Sort      string `json:"sort"`
	Direction string `json:"direction"`
}
