package ogmintservice

import (
	"context"
	"fmt"
	apiLogger "pismo-dev/api/logger"
	commontypes "pismo-dev/commonpkg/types"
	"pismo-dev/constants"
	"pismo-dev/internal/appconfig"
	"pismo-dev/internal/models"
	"pismo-dev/internal/types"
	"pismo-dev/pkg/cache"
	kafkaclient "pismo-dev/pkg/kafka/client"
	"pismo-dev/pkg/logger"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type OGMintSvc struct {
	ServiceName         string
	Queues              map[string]commontypes.QueueWrapperConfig
	CacheW              *cache.CacheWrapper[string, string]
	MintConfigs         map[string]commontypes.MintConfig
	KafkaProducerClient *kafkaclient.ProducerClient
	services            RequiredServices
	repos               RequiredRepos
	AmplitudeEventTopic string
}

func NewOGMintSvc(queues map[string]commontypes.QueueWrapperConfig, cacheW *cache.CacheWrapper[string, string], mintConfigs map[string]commontypes.MintConfig, kafkaProducerClient *kafkaclient.ProducerClient, amplitudeEventTopic string) *OGMintSvc {
	return &OGMintSvc{
		ServiceName:         "OGMintService",
		Queues:              queues,
		CacheW:              cacheW,
		MintConfigs:         mintConfigs,
		KafkaProducerClient: kafkaProducerClient,
		AmplitudeEventTopic: amplitudeEventTopic,
	}
}

func (svc *OGMintSvc) SetRequiredServices(services RequiredServices) {
	svc.services = services
}

func (svc *OGMintSvc) SetRequiredRepos(repos RequiredRepos) {
	svc.repos = repos
}

func (svc *OGMintSvc) CheckTotalMints(ctx context.Context, contractAddress string) (int, error) {
	count := svc.repos.OrderMetaDataRepo.GetTotalOGMintCount(contractAddress)
	if count < 0 {
		logger.Error("OGMintFlow:Error fetching OG Mint Count", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.DB_ERROR].ErrorCode})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.DB_ERROR].HttpStatus, constants.ERROR_TYPES[constants.DB_ERROR].ErrorCode, "DB Error")
		return -1, fmt.Errorf("Error fetching OG Mint Count")
	}
	return count, nil

}

func (svc *OGMintSvc) GetUserMints(ctx context.Context, contractAddress string, userId string) (int, error) {
	count := svc.repos.OrderMetaDataRepo.GetUserNftCollectionOrderCountByOrderTypeAndStatus(userId, contractAddress, constants.CHECK_ALLOW_OG_MINT_STATUS, constants.OG_MINT)
	if count < 0 {
		logger.Error("OGMintFlow:Error fetching User  OG Mint Count", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.DB_ERROR].ErrorCode})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.DB_ERROR].HttpStatus, constants.ERROR_TYPES[constants.DB_ERROR].ErrorCode, "DB Error")
		return -1, fmt.Errorf("Error fetching User OG Mint Count")
	}
	return count, nil

}

func (svc *OGMintSvc) PushOrderForProcessing(ctx context.Context, userId, collectionId, networkId, orderType, customUriString string) (string, error) {
	// Create waitGroup and channels to execute parallel calls
	var wg sync.WaitGroup

	signingServicechannel := make(chan types.SigningServiceChannel)
	networkChannel := make(chan types.NetworkChannel)
	collectionDetailsChannel := make(chan types.CollectionDqlChannel)
	wg.Add(3)
	// calling functions on different go routines

	go svc.signingServiceInteraction(&wg, signingServicechannel, ctx, userId, networkId)
	go svc.DqlNetworkInteraction(&wg, networkChannel, ctx, networkId)
	go svc.DqlCollectionInteraction(&wg, collectionDetailsChannel, ctx, collectionId)
	// Creating listeners for our channels

	signingServiceResult := <-signingServicechannel
	networkResult := <-networkChannel
	collectionDetailsResult := <-collectionDetailsChannel
	// Waiting for all go routines to finish
	wg.Wait()
	// close all channels to enable reading data

	close(signingServicechannel)
	close(networkChannel)
	close(collectionDetailsChannel)
	// Process each channel data

	var networkName string

	err := networkResult.Error
	if err != nil {
		logger.Error("ExecuteFlow: Error Fetching network details", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "collectionId": collectionId, "error": err.Error()})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, err.Error())
		return "", fmt.Errorf("DQL Error")
	}
	networkName = networkResult.Result

	var walletAddress string

	err = signingServiceResult.Error
	if err != nil {
		logger.Error("OGMintFlow: Error Fetching user address", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "user": userId, "error": err.Error()})
		if strings.Contains(signingServiceResult.Result.Error.Message, constants.WALLET_NOT_BACKEDUP_SS_ERROR) {
			apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.WALLET_NOT_BACKED_UP_ERROR].HttpStatus, constants.ERROR_TYPES[constants.WALLET_NOT_BACKED_UP_ERROR].ErrorCode, "Wallet not backed up by user")
		} else {
			apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "Cannot Fetch User Wallet Address")
		}
		return "", err
	}
	walletAddress = signingServiceResult.Result.Address

	var collectionDetails types.DQLNftCollectionResponseInterface

	err = collectionDetailsResult.Error
	collectionDetails = collectionDetailsResult.Result
	if err != nil || (collectionDetails.Details == types.DQLCollectionDetails{}) || (collectionDetails.Details.ContractMetadata == types.DQLCollectionMetaData{}) {
		logger.Error("OGMintFlow:Error fetching collection details", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "collectionId": collectionId})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "Internal Error")
		return "", fmt.Errorf("DQL Error")
	}

	// Run validations
	if networkId != collectionDetails.Details.NetworkID {
		logger.Error("OGMintFlow:Error fetching collection details: Network Mismatch", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "collectionId": collectionId})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "Internal Error")
		return "", fmt.Errorf("Network Mismatch Error")
	}
	if len(collectionDetails.Details.ContractMetadata.VendorName) == 0 {
		logger.Error("OGMintFlow:Error fetching collection details: Vendor not provided in collection", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "collectionId": collectionId})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "Internal Error")
		return "", fmt.Errorf("Network Mismatch Error")
	}
	existingMints, err := svc.CheckTotalMints(ctx, collectionDetails.Details.ContractAddress)
	if err != nil {
		return "", fmt.Errorf("Total Mint count error")
	}
	mintLimit := 0
	mintLimit, _ = strconv.Atoi(collectionDetails.Details.ContractMetadata.NftMintLimit)
	if existingMints >= mintLimit {
		logger.Error("OGMintFlow:Offer expired for collection", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.COLLECTION_LIMIT_ERROR].ErrorCode, "collectionId": collectionId})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.COLLECTION_LIMIT_ERROR].HttpStatus, constants.ERROR_TYPES[constants.COLLECTION_LIMIT_ERROR].ErrorCode, "Collection limit exceeded")
		return "", fmt.Errorf("Collection Limit Exceed")
	}
	userMints, err := svc.GetUserMints(ctx, collectionDetails.Details.ContractAddress, userId)
	if err != nil {
		return "", fmt.Errorf("Total Mint count error")
	}
	userLimit := 0
	userLimit, _ = strconv.Atoi(collectionDetails.Details.ContractMetadata.MintLimitPerWallet)
	if userMints >= userLimit {
		logger.Error("OGMintFlow:User Already minted NFT", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.USER_OG_LIMIT_ERROR].ErrorCode, "collectionId": collectionId, "userId": userId})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.USER_OG_LIMIT_ERROR].HttpStatus, constants.ERROR_TYPES[constants.USER_OG_LIMIT_ERROR].ErrorCode, "User Mint limit exceeded")
		return "", fmt.Errorf("User Collection Limit Exceed")
	}

	orderMetaDataObj := svc.prepareOrderMetaData(userId, networkId, collectionDetails.Details.NftType, collectionDetails.Details.ContractAddress, orderType)
	orderId := orderMetaDataObj.OrderId
	savedOrder := svc.repos.OrderMetaDataRepo.InsertOrderMetadata(&orderMetaDataObj)
	if savedOrder.Error != nil {
		logger.Error("OGMintFlow: Error saving orderMetaData", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "error": savedOrder.Error.Error(), "orderData": orderMetaDataObj})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.PROCESS_ERROR].HttpStatus, constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "error creating orderMetaData")
		return "", fmt.Errorf("error creating orderMetaData")
	}
	queueObj := svc.Queues[constants.COLLECTION_VENDOR_MAPPING[strings.ToUpper(strings.Trim(collectionDetails.Details.ContractMetadata.VendorName, " "))]+networkId]
	queueWrapper := queueObj.WrapperObj
	queueMessage := types.OGMintQueueMessage{
		OrderId:     orderId,
		Type:        constants.MINT,
		RetryCount:  0,
		UserId:      userId,
		Status:      constants.CREATED,
		ToAddress:   walletAddress,
		NetworkId:   networkId,
		MetaDataURI: collectionDetails.Details.ContractMetadata.NftMetadataURI,
		OrderType:   orderType,
		VendorName:  strings.ToUpper(strings.Trim(collectionDetails.Details.ContractMetadata.VendorName, " ")),
	}
	if len(collectionDetails.Details.ContractMetadata.CrossmintCollectionId) > 0 {
		queueMessage.ContractIdentifier = collectionDetails.Details.ContractMetadata.CrossmintCollectionId
	} else {
		queueMessage.ContractIdentifier = collectionDetails.Details.ContractAddress
	}
	if len(customUriString) > 0 {
		queueMessage.MetaDataURI = customUriString
	}
	queueWrapper.PublishMessage(ctx, queueMessage, orderId, queueObj.QueueProps.SQSDelaySeconds)
	amplitudeEvent := types.AmplitudeMintEventInterface{
		AppName:          constants.AMPLITUDE_APP_NAME,
		Network:          networkId,
		NftType:          strings.ToUpper(collectionDetails.Details.NftType),
		TokenId:          "",
		Type:             constants.AMPLITUDE_MINT_EVENT_TYPE,
		Status:           constants.AMPLITUDE_STATUS_MAPPING[orderMetaDataObj.Status],
		Product:          "nft",
		CollectionName:   collectionDetails.Details.CollectionName,
		CollectionId:     collectionId,
		TokenCount:       1,
		ReceiversAddress: walletAddress,
		Chain:            networkName,
	}
	err = svc.NotifyAmplitudeEvent(orderId, userId, amplitudeEvent)
	if err != nil {
		logger.Error("Error Notifying Amplitude Service", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "user": userId, "orderId": orderId})
	}
	return orderId, nil

}

func (svc *OGMintSvc) prepareOrderMetaData(userId, networkId, ercType, collectionAddress, orderType string) models.OrderMetadata {
	return models.OrderMetadata{
		OrderId:           uuid.New().String(),
		UserId:            userId,
		NetworkId:         networkId,
		Status:            constants.CREATED,
		EntityType:        strings.ToUpper(ercType),
		EntityAddress:     collectionAddress,
		NftId:             "TBD",
		Count:             "1",
		OrderType:         orderType,
		Slippage:          appconfig.SLIPPAGE,
		ExecutionResponse: "",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		RetryCount:        0,
	}
}

func (svc *OGMintSvc) NotifyAmplitudeEvent(orderId, userId string, event types.AmplitudeMintEventInterface) error {
	if len(event.NftType) == 0 {
		event.NftType = "ERC721"
	}

	amplitudeNotification := types.AmplitudeEventInterface{
		UserId:          userId,
		Eventtype:       constants.AMPLITUDE_EVENT_NAME,
		EventProperties: event,
	}
	return svc.KafkaProducerClient.Produce(svc.AmplitudeEventTopic, orderId, amplitudeNotification)
}

func (svc *OGMintSvc) signingServiceInteraction(wg *sync.WaitGroup, channel chan types.SigningServiceChannel, ctx context.Context, userId, networkId string) {
	defer wg.Done()
	signingSvcResponse, err := svc.services.SigningSvc.GetUserWalletAddress(ctx, userId, networkId)
	channel <- types.SigningServiceChannel{
		Result: signingSvcResponse,
		Error:  err,
	}

}

func (svc *OGMintSvc) DqlCollectionInteraction(wg *sync.WaitGroup, channel chan types.CollectionDqlChannel, ctx context.Context, identifier string) {
	defer wg.Done()
	collectionDetails, err := svc.services.DQLSvc.GetNftCollectionById(ctx, identifier)
	channel <- types.CollectionDqlChannel{
		Result: collectionDetails,
		Error:  err,
	}
}

func (svc *OGMintSvc) DqlNetworkInteraction(wg *sync.WaitGroup, channel chan types.NetworkChannel, ctx context.Context, identifier string) {
	defer wg.Done()
	var dqlError error
	networkName, found, err := svc.CacheW.Driver.Get(identifier, constants.NETWORK_CACHE_PURPOSE)
	if err != nil || !found {
		if err != nil {
			logger.Info("Network Cache not found error", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.INVALID_OTP_ERROR].ErrorCode, "network": identifier, "error": err.Error()})
		}
		if !found {
			logger.Info("Network Cache not found", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.INVALID_OTP_ERROR].ErrorCode, "network": identifier})
		}
		var dqlResponse types.DQLEntityResponseInterface
		dqlResponse, dqlError = svc.services.DQLSvc.GetEntityById(ctx, identifier, false)
		svc.CacheW.Driver.SetEx(identifier, dqlResponse.Details.Name, constants.NETWORK_CACHE_PURPOSE, time.Duration(appconfig.NETWORK_CACHE_TTL)*time.Minute)
		networkName = dqlResponse.Details.Name
	}
	channel <- types.NetworkChannel{
		Result: networkName,
		Error:  dqlError,
	}
}
