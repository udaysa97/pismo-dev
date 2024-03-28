package portfolio

import (
	"context"
	"pismo-dev/external/types"
)

type PortfolioInterface interface {
	GetUserBalance(ctx context.Context, nftId string, networkId string, userId string) (types.TokenBalanceResponse, error)
}
