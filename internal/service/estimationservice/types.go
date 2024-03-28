package estimationservice

import (
	"pismo-dev/external/DQL"
	"pismo-dev/external/FM"
	"pismo-dev/external/portfolio"
	"pismo-dev/external/signing"
	"pismo-dev/internal/repository/ordermetadata"
)

type RequiredServices struct {
	FlowManagerSvc FM.FMInterface
	PortfolioSvc   portfolio.PortfolioInterface
	DQLSvc         DQL.DQLInterface
	SigningSvc     signing.SigningInterface
}

type RequiredRepos struct {
	OrderMetaDataRepo ordermetadata.OrderMetadataRepositoryInterface
}
