package initializer

import (
	commontypes "pismo-dev/commonpkg/types"
	"pismo-dev/constants"
	"pismo-dev/internal/appconfig"
	"pismo-dev/internal/repository"
	"pismo-dev/internal/service"
	"pismo-dev/internal/service/estimationservice"
	"pismo-dev/internal/service/ogmintservice"
	"pismo-dev/internal/service/orderexecutionservice"
	"pismo-dev/internal/service/orderexecutiontracker"
	"pismo-dev/internal/service/otpservice"
	"pismo-dev/internal/service/reconcileservice"

	"pismo-dev/pkg/cache"
	kafkaclient "pismo-dev/pkg/kafka/client"
)

func InitServices(repositories *repository.Repositories, cacheWstring *cache.CacheWrapper[string, string], cacheWint *cache.CacheWrapper[string, int], kafkaProducer *kafkaclient.ProducerClient, kafkaConsumer *kafkaclient.ConsumerClient, queueWrapperInterface map[string]commontypes.QueueWrapperConfig) (services *service.Service) {
	options := []service.ServiceOption{
		service.WithEstimationSvc(kafkaProducer),
		service.WithExecutionSvc(kafkaProducer, queueWrapperInterface[constants.RECONQUEUE].WrapperObj, cacheWstring),
		service.WithOtpSvc(cacheWstring, cacheWint),
		service.WithDQLSvc(),
		service.WithFlowManagerSvc(),
		service.WithPortfolioSvc(),
		service.WithAuthSvc(),
		service.WithSigningSvc(),
		service.WithEmailSvc(),
		service.WithNftPortSvc(),
		service.WithCrossMintSvc(),
		service.WithOrderSvc(repositories.OrderMetadataRepo),
		service.WithTransactionDataSvc(repositories.TransactionDataRepo),
		service.WithReconcileSvc(queueWrapperInterface[constants.RECONQUEUE].WrapperObj),
		service.WithJobExecTrackerSvc(kafkaProducer, kafkaConsumer),
		service.WithOgMintSvc(queueWrapperInterface, cacheWstring, kafkaProducer),
	}
	services = service.NewService(options...)
	if services.EstimationSvc != nil {
		services.EstimationSvc.SetRequiredServices(estimationservice.RequiredServices{
			FlowManagerSvc: services.FlowManagerSvc,
			PortfolioSvc:   services.PortfolioSvc,
			DQLSvc:         services.DQLSvc,
			SigningSvc:     services.SigningSvc,
		})
		services.EstimationSvc.SetRequiredRepos(estimationservice.RequiredRepos{
			OrderMetaDataRepo: repositories.OrderMetadataRepo,
		})
	}
	if services.OtpSvc != nil {
		services.OtpSvc.SetRequiredServices(otpservice.RequiredServices{
			AuthSvc:  services.AuthSvc,
			EmailSvc: services.EmailSvc,
		})
	}

	if services.OrderExecTrackerSvc != nil {
		services.OrderExecTrackerSvc.SetRequiredRepos(orderexecutiontracker.RequiredRepos{
			OrderMetaDataRepo:   repositories.OrderMetadataRepo,
			TransactionDataRepo: repositories.TransactionDataRepo,
		})
		go services.OrderExecTrackerSvc.InitJobTrackerConsumer()
	}

	if services.ReconcileSvc != nil {
		services.ReconcileSvc.SetRequiredServices(reconcileservice.RequiredServices{
			FlowManagerSvc: services.FlowManagerSvc,
		})
		services.ReconcileSvc.SetRequiredRepos(reconcileservice.RequiredRepos{
			OrderMetaDataRepo:   repositories.OrderMetadataRepo,
			TransactionDataRepo: repositories.TransactionDataRepo,
		})
		go services.ReconcileSvc.InitiateSqsConsumer()
	}

	if services.OgMintSvc != nil {
		services.OgMintSvc.SetRequiredServices(ogmintservice.RequiredServices{
			SigningSvc:   services.SigningSvc,
			NFTPORTSvc:   services.NftPortSvc,
			CrossMintSvc: services.CrossMintSvc,
			DQLSvc:       services.DQLSvc,
		})
		services.OgMintSvc.SetRequiredRepos(ogmintservice.RequiredRepos{
			OrderMetaDataRepo:   repositories.OrderMetadataRepo,
			TransactionDataRepo: repositories.TransactionDataRepo,
		})
		mintConfigs := appconfig.GetMintConfigs()
		for key, _ := range mintConfigs {
			go services.OgMintSvc.InitiateSqsConsumer(key)
		}
	}

	if services.ExecutionSvc != nil {
		services.ExecutionSvc.SetRequiredServices(orderexecutionservice.RequiredServices{
			FlowManagerSvc: services.FlowManagerSvc,
			PortfolioSvc:   services.PortfolioSvc,
			DQLSvc:         services.DQLSvc,
			OtpSvc:         services.OtpSvc,
			ReconcileSvc:   services.ReconcileSvc,
			SigningSvc:     services.SigningSvc,
			OgmintSvc:      services.OgMintSvc,
		})
		services.ExecutionSvc.SetRequiredRepos(orderexecutionservice.RequiredRepos{
			OrderMetaDataRepo: repositories.OrderMetadataRepo,
		})
	}

	return
}
