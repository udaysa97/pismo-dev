package orderexecutionservice

import (
	"context"
	apiTypes "pismo-dev/api/types"
	"pismo-dev/internal/types"
)

type ExecutionServiceInterface interface {
	ExecuteOrder(ctx context.Context, userDetails types.UserDetailsInterface, nftTokenDetails types.NFTTransferDetailsInterface) (types.FMExecuteResult, error)
	SetRequiredServices(services RequiredServices)
	SetRequiredRepos(repos RequiredRepos)
	OpenApiExecuteOrder(ctx context.Context, nftRequest apiTypes.OpenApiExecuteRequest) (types.FMExecuteResult, error)
	OpenApiMintOrder(ctx context.Context, nftRequest apiTypes.OpenAPINFTMintOrder) (types.FMExecuteResult, error)
}
