package orderexecutiontracker

import (
	"pismo-dev/internal/repository/ordermetadata"
	"pismo-dev/internal/repository/transactiondata"
)

type RequiredRepos struct {
	OrderMetaDataRepo   ordermetadata.OrderMetadataRepositoryInterface
	TransactionDataRepo transactiondata.TransactionDataRepositoryInterface
}
