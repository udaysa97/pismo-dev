package ogmintservice

import (
	"pismo-dev/external/DQL"
	"pismo-dev/external/crossmint"
	"pismo-dev/external/nftport"
	"pismo-dev/external/signing"
	"pismo-dev/internal/repository/ordermetadata"
	"pismo-dev/internal/repository/transactiondata"
)

type RequiredServices struct {
	NFTPORTSvc   nftport.NftPortInterface
	CrossMintSvc crossmint.CrossMintInterface
	DQLSvc       DQL.DQLInterface
	SigningSvc   signing.SigningInterface
}

type RequiredRepos struct {
	OrderMetaDataRepo   ordermetadata.OrderMetadataRepositoryInterface
	TransactionDataRepo transactiondata.TransactionDataRepositoryInterface
}
