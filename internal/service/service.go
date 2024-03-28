package service

import (
	"pismo-dev/external/DQL"
	"pismo-dev/external/FM"
	"pismo-dev/external/auth"
	"pismo-dev/external/crossmint"
	"pismo-dev/external/email"
	"pismo-dev/external/nftport"
	"pismo-dev/external/portfolio"
	"pismo-dev/external/signing"
	"pismo-dev/internal/appconfig"
	"pismo-dev/internal/repository/ordermetadata"
	"pismo-dev/internal/repository/transactiondata"
	"pismo-dev/internal/service/estimationservice"
	"pismo-dev/internal/service/ogmintservice"
	"pismo-dev/internal/service/orderexecutionservice"
	"pismo-dev/internal/service/orderexecutiontracker"
	"pismo-dev/internal/service/orderservice"
	"pismo-dev/internal/service/otpservice"
	"pismo-dev/internal/service/reconcileservice"
	"pismo-dev/internal/service/transactiondataservice"
	"pismo-dev/pkg/cache"

	commontypes "pismo-dev/commonpkg/types"
	"pismo-dev/pkg/httpclient"
	"pismo-dev/pkg/httpclient/nethttp/client"
	httpclientDriver "pismo-dev/pkg/httpclient/nethttp/drivers"
	kafkaclient "pismo-dev/pkg/kafka/client"
	"pismo-dev/pkg/queue"

	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
)

type Service struct {
	EstimationSvc       estimationservice.EstimationServiceInterface
	ExecutionSvc        orderexecutionservice.ExecutionServiceInterface
	OtpSvc              otpservice.OTPServiceInterface
	DQLSvc              DQL.DQLInterface
	PortfolioSvc        portfolio.PortfolioInterface
	FlowManagerSvc      FM.FMInterface
	SigningSvc          signing.SigningInterface
	ReconcileSvc        reconcileservice.ReconcileServiceInterface
	AuthSvc             auth.AuthInterface
	EmailSvc            email.EmailInterface
	OrderSvc            orderservice.OrderServiceInterface
	TransactionDataSvc  transactiondataservice.TransactionServiceInterface
	OrderExecTrackerSvc orderexecutiontracker.OrderExecutionTrackerInterface
	NftPortSvc          nftport.NftPortInterface
	CrossMintSvc        crossmint.CrossMintInterface
	OgMintSvc           ogmintservice.OGMintServiceInterface
}

type ServiceOption func(*Service) error

func NewService(opts ...ServiceOption) *Service {
	service := &Service{}

	// Loop through each option
	for _, opt := range opts {
		// Call the option giving the instantiated
		// *service as the argument
		err := opt(service)
		if err != nil {
			return nil
		}
	}

	return service
}

func WithEstimationSvc(kafkaProducerClient *kafkaclient.ProducerClient) ServiceOption {
	return func(svc *Service) error {
		svc.EstimationSvc = estimationservice.NewEstimationSvc(kafkaProducerClient)
		return nil
	}
}

func WithExecutionSvc(kafkaProducerClient *kafkaclient.ProducerClient, sqsqw queue.QueueWrapper, cacheWstring *cache.CacheWrapper[string, string]) ServiceOption {
	return func(svc *Service) error {
		svc.ExecutionSvc = orderexecutionservice.NewOrderExecutionSvc(kafkaProducerClient, appconfig.TOKEN_TRANSFER_TOPIC, appconfig.AMPLITUDE_EVENT_TOPIC, sqsqw, appconfig.SQS_DELAY_SECONDS, cacheWstring)
		return nil
	}
}

func WithOtpSvc(cacheWstring *cache.CacheWrapper[string, string], cacheWint *cache.CacheWrapper[string, int]) ServiceOption {
	return func(svc *Service) error {
		svc.OtpSvc = otpservice.NewOTPSvc(cacheWstring, cacheWint)
		return nil
	}
}

func WithPortfolioSvc() ServiceOption {
	portfolioHttpClient := httptrace.WrapClient(client.NewCustomNetHttpClient())
	httpClient := httpclient.NewHttpClientWrapper("HttpClient", httpclientDriver.NewNetHttpClient(portfolioHttpClient))
	return func(svc *Service) error {
		svc.PortfolioSvc = portfolio.NewPortfolioSvc(httpClient, appconfig.PORTFOLIO_SERVICE_URL, appconfig.OKTO_VPC_SECRET)
		return nil
	}
}

func WithDQLSvc() ServiceOption {
	DQLHttpClient := httptrace.WrapClient(client.NewCustomNetHttpClient())
	httpClient := httpclient.NewHttpClientWrapper("HttpClient", httpclientDriver.NewNetHttpClient(DQLHttpClient))
	return func(svc *Service) error {
		svc.DQLSvc = DQL.NewDQLSvc(httpClient, appconfig.DQL_SERVICE_URL, appconfig.OKTO_VPC_SECRET)
		return nil
	}
}

func WithFlowManagerSvc() ServiceOption {
	FMHttpClient := httptrace.WrapClient(client.NewCustomNetHttpClient())
	httpClient := httpclient.NewHttpClientWrapper("HttpClient", httpclientDriver.NewNetHttpClient(FMHttpClient))
	return func(svc *Service) error {
		svc.FlowManagerSvc = FM.NewFMSvc(httpClient, appconfig.FM_SERVICE_URL, appconfig.DAPP_SERVICE_URL, appconfig.OKTO_VPC_SECRET) //TOCHANGE HERE
		return nil
	}
}

func WithAuthSvc() ServiceOption {
	FMHttpClient := httptrace.WrapClient(client.NewCustomNetHttpClient())
	httpClient := httpclient.NewHttpClientWrapper("HttpClient", httpclientDriver.NewNetHttpClient(FMHttpClient))
	return func(svc *Service) error {
		svc.AuthSvc = auth.NewAuthSvc(httpClient, appconfig.AUTH_SERVICE_BASE_URL, appconfig.CEFI_VPC_SECRET) //TOCHANGE HERE
		return nil
	}
}

func WithEmailSvc() ServiceOption {
	FMHttpClient := httptrace.WrapClient(client.NewCustomNetHttpClient())
	httpClient := httpclient.NewHttpClientWrapper("HttpClient", httpclientDriver.NewNetHttpClient(FMHttpClient))
	return func(svc *Service) error {
		svc.EmailSvc = email.NewEmailSvc(httpClient, appconfig.EMAIL_SERVICE_BASE_URL, appconfig.COMMUNICATION_VPC_AUTHORIZATION_SECRET) //TOCHANGE HERE
		return nil
	}
}

func WithReconcileSvc(sqsqw queue.QueueWrapper) ServiceOption {
	return func(svc *Service) error {
		svc.ReconcileSvc = reconcileservice.NewReconcileSvc(sqsqw, appconfig.SQS_WAIT_TIME, appconfig.SQS_VISIBILITY_TIMEOUT, appconfig.MAX_NUMBER_SQS_MESSAGE)
		return nil
	}
}

func WithOrderSvc(ordermetadata ordermetadata.OrderMetadataRepositoryInterface) ServiceOption {
	return func(svc *Service) error {
		svc.OrderSvc = orderservice.NewOrderService(ordermetadata)
		return nil

	}
}

func WithTransactionDataSvc(transactiondata transactiondata.TransactionDataRepositoryInterface) ServiceOption {
	return func(svc *Service) error {
		svc.TransactionDataSvc = transactiondataservice.NewTransactionDataService(transactiondata)
		return nil

	}
}

func WithJobExecTrackerSvc(kafkaProducerClient *kafkaclient.ProducerClient, kafkaConsumerClient *kafkaclient.ConsumerClient) ServiceOption {
	return func(svc *Service) error {
		svc.OrderExecTrackerSvc = orderexecutiontracker.NewOrderExecutionTrackerSvc(kafkaProducerClient, kafkaConsumerClient, appconfig.AMPLITUDE_EVENT_TOPIC, appconfig.FM_CONSUMER_TOPIC)
		return nil

	}
}

func WithSigningSvc() ServiceOption {
	FMHttpClient := httptrace.WrapClient(client.NewCustomNetHttpClient())
	httpClient := httpclient.NewHttpClientWrapper("HttpClient", httpclientDriver.NewNetHttpClient(FMHttpClient))
	return func(svc *Service) error {
		svc.SigningSvc = signing.NewSigningSvc(httpClient, appconfig.SIGNING_SERVICE_BASE_URL, appconfig.SIGNING_VPC_SECRET) //TOCHANGE HERE
		return nil
	}
}

func WithNftPortSvc() ServiceOption {
	FMHttpClient := httptrace.WrapClient(client.NewCustomNetHttpClient())
	httpClient := httpclient.NewHttpClientWrapper("HttpClient", httpclientDriver.NewNetHttpClient(FMHttpClient))
	return func(svc *Service) error {
		svc.NftPortSvc = nftport.NewNftPortSvc(httpClient, appconfig.NFTPORT_BASE_URL, appconfig.NFTPORT_API_KEY) //TOCHANGE HERE
		return nil
	}
}

func WithCrossMintSvc() ServiceOption {
	CrossMintHttpClient := client.NewCustomNetHttpClient()
	httpClient := httpclient.NewHttpClientWrapper("HttpClient", httpclientDriver.NewNetHttpClient(CrossMintHttpClient))
	return func(svc *Service) error {
		svc.CrossMintSvc = crossmint.NewCrossMintSvc(httpClient, appconfig.CROSSMINT_BASE_URL, appconfig.CROSSMINT_CLIENT_KEY, appconfig.CROSSMINT_PROJECT_KEY)
		return nil
	}
}

func WithOgMintSvc(queueWrapperInterface map[string]commontypes.QueueWrapperConfig, cacheW *cache.CacheWrapper[string, string], kafkaProducerClient *kafkaclient.ProducerClient) ServiceOption {
	return func(svc *Service) error {
		mintConfigs := appconfig.GetMintConfigs()
		svc.OgMintSvc = ogmintservice.NewOGMintSvc(queueWrapperInterface, cacheW, mintConfigs, kafkaProducerClient, appconfig.AMPLITUDE_EVENT_TOPIC)
		return nil

	}
}
