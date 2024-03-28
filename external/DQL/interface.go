package DQL

import (
	"context"
	"pismo-dev/internal/types"
)

type DQLInterface interface {
	GetEntityById(ctx context.Context, entityId string, isNft bool) (types.DQLEntityResponseInterface, error)
	GetTokenByAddress(ctx context.Context, address string, networkId string) (types.DQLByAddressResponseInterface, error)
	GetNftCollectionById(ctx context.Context, entityId string) (types.DQLNftCollectionResponseInterface, error)
	GetNftByCollectionAndTokenId(ctx context.Context, collectionAddress string, nftTokenId string, networkId string) (types.DQLByAddressResponseInterface, error)
}
