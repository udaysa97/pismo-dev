package types

import (
	"pismo-dev/internal/models"
	"pismo-dev/internal/types"
)

type Status string

var (
	StatusSuccess Status = "success"
	StatusError   Status = "error"
)

const (
	BADREQUEST    = "ER-TECH-0001"
	RUNTIMEERROR  = "ER-TECH-0005"
	UNPROCESSABLE = "ER-TECH-0013"
)

type ErrorResponseWithCount struct {
	Items any    `json:"items,omitempty"`
	Count uint64 `json:"count,omitempty"`
}

type ErrorResponse struct {
	Code      int    `json:"code,omitempty"`
	ErrorCode string `json:"errorCode,omitempty"`
	Message   string `json:"message,omitempty"`
	TraceID   string `json:"trace_id,omitempty"`
	Status    Status `json:"status,omitempty"`
}
type ResponseDTO[T any] struct {
	Status  Status         `json:"status,omitempty"`
	Success bool           `json:"success"`
	OrderId string         `json:"order_id,omitempty"`
	Error   *ErrorResponse `json:"error,omitempty"`
	Result  *T             `json:"data,omitempty"`
}

type SampleResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type EndpointResponse struct {
	Status Status
	Error  *ErrorResponse
	Result map[string]any
}

type EstimationResponseWrapperInterface struct {
	Estimation EstimateResponse `json:"estimation"`
}

type EstimateResponse struct {
	TransactionFee    []types.FMEstimateTxFee  `json:"transactionFee"`
	IsGsnRequired     bool                     `json:"isGsnRequired"`
	IsGsnPossible     bool                     `json:"isGsnPossible"`
	GsnWithdrawTokens []types.GsnWithdrawToken `json:"gsnWithdrawTokens"`
	OrderId           string                   `json:"orderId"`
}

type ExecuteResponse struct {
	OrderId string `json:"orderId"`
}

type MintOGResponse struct {
	OrderId string `json:"orderId"`
}

type ExecuteResponseDTO struct {
	Data ExecuteResponse `json:"data"`
}

type OrderDetailsResponse struct {
	models.OrderDetails
}

type TransactionDetailsResponse struct {
	models.TransactionDetails
}

type GetUserCollectionMintCountsResponse struct {
	UserMints  int `json:"user_mints"`
	TotalMints int `json:"total_mints"`
}
