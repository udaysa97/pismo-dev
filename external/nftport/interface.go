package nftport

import (
	"context"
	"pismo-dev/external/types"
)

type NftPortInterface interface {
	MintNft(ctx context.Context, chain, contract_address, metadata_uri, to_address string) (types.NftPortMintResponse, error)
	FetchMintStatus(ctx context.Context, chain, txHash string) (types.NftPortStatusResponse, error)
}
