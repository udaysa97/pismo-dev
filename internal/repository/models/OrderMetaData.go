package models

import (
	"time"
)

type OrderMetaData struct {
	OrderId           string    `json:"orderId"`
	UserId            string    `json:"userId"`
	Status            string    `json:"status"`
	EntityType        string    `json:"entityType"`
	EntityId          string    `json:"entityId"`
	NetworkId         string    `json:"networkId"`
	Transactions      []string  `json:"transactions"`
	GasFee            string    `json:"gasFee"`
	Slippage          string    `json:"slippage"`
	OrderType         string    `json:"orderType"`
	ExecutionResponse any       `json:"executionResponse"`
	ContractAddress   string    `json:"contractAddress"`
	ErcType           string    `json:"ercType"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
	RetryCount        int       `json:"retryCount"`
}
