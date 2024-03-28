package signing

import (
	"context"
	"pismo-dev/external/types"
)

type SigningInterface interface {
	GetUserWalletAddress(ctx context.Context, userId, networkId string) (types.SigningSvcResponse, error)
}
