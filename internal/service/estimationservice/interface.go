package estimationservice

import (
	"context"
	"pismo-dev/internal/types"
)

type EstimationServiceInterface interface {
	CalculateEstimate(ctx context.Context, userDetails types.UserDetailsInterface, nftTokenDetails types.NFTTransferDetailsInterface) (types.FMEstimateResultOutput, error)
	SetRequiredServices(services RequiredServices)
	SetRequiredRepos(repos RequiredRepos)
}
