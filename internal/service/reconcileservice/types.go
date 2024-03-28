package reconcileservice

import (
	"pismo-dev/external/FM"
	"pismo-dev/internal/repository/ordermetadata"
	"pismo-dev/internal/repository/transactiondata"
)

type RequiredServices struct {
	FlowManagerSvc FM.FMInterface
}

type RequiredRepos struct {
	OrderMetaDataRepo   ordermetadata.OrderMetadataRepositoryInterface
	TransactionDataRepo transactiondata.TransactionDataRepositoryInterface
}
