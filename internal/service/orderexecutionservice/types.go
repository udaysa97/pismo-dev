package orderexecutionservice

import (
	"pismo-dev/external/DQL"
	"pismo-dev/external/FM"
	"pismo-dev/external/portfolio"
	"pismo-dev/external/signing"
	"pismo-dev/internal/repository/ordermetadata"
	"pismo-dev/internal/service/ogmintservice"
	"pismo-dev/internal/service/otpservice"
	"pismo-dev/internal/service/reconcileservice"
)

type RequiredServices struct {
	FlowManagerSvc FM.FMInterface
	PortfolioSvc   portfolio.PortfolioInterface
	DQLSvc         DQL.DQLInterface
	ReconcileSvc   reconcileservice.ReconcileServiceInterface
	OtpSvc         otpservice.OTPServiceInterface
	SigningSvc     signing.SigningInterface
	OgmintSvc      ogmintservice.OGMintServiceInterface
}

type RequiredRepos struct {
	OrderMetaDataRepo ordermetadata.OrderMetadataRepositoryInterface
}
