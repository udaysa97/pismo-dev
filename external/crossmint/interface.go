package crossmint

import (
	"context"
	"pismo-dev/external/types"
)

type CrossMintInterface interface {
	MintNft(ctx context.Context, chain, contractIdentifier, metadataUri, toAddress string, isBYOC bool) (types.CrossMintMintResponse, error)
	FetchMintStatus(ctx context.Context, chain, txHash, contractIdentifier string) (types.CrossMintStatusResponse, error)
}
