package FM

import (
	"context"
	"pismo-dev/internal/types"
)

type FMInterface interface {
	GetEstimate(ctx context.Context, payload types.FMEstimateRequest) (types.FMEstimateResult, error)
	ExecuteOrder(ctx context.Context, payload types.FMExecuteRequest) (types.FMExecuteResult, error)
	ExecuteOpenAPIOrder(ctx context.Context, payload types.OpenAPINFTOrder, orderType string, vendorId string) (types.FMExecuteResult, error)
	GetStatus(ctx context.Context, jobId string) (types.FMGetStatusResult, error)
}
