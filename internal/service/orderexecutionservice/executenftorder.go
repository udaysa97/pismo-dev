package orderexecutionservice

import (
	"context"
	"fmt"
	apiLogger "pismo-dev/api/logger"
	"pismo-dev/constants"
	"pismo-dev/internal/appconfig"
	"pismo-dev/internal/models"
	"pismo-dev/internal/types"
	"pismo-dev/pkg/cache"
	kafkaclient "pismo-dev/pkg/kafka/client"
	"pismo-dev/pkg/logger"
	"pismo-dev/pkg/queue"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type OrderExecutionSvc struct {
	ServiceName         string
	services            RequiredServices
	repos               RequiredRepos
	KafkaProducerClient *kafkaclient.ProducerClient
	TokenTransferTopic  string
	AmplitudeEventTopic string
	SqsQueueWrapper     *queue.QueueWrapper
	SQSDelaySeconds     int
	CacheW              *cache.CacheWrapper[string, string]
}

func NewOrderExecutionSvc(kafkaProducerClient *kafkaclient.ProducerClient, tokenTransferTopic string, amplitudeEventTopic string, SqsQueueWrapper queue.QueueWrapper, sqsDelaySeconds int, cacheW *cache.CacheWrapper[string, string]) *OrderExecutionSvc {
	return &OrderExecutionSvc{
		ServiceName:         "NFTOrderExecutionSvc",
		KafkaProducerClient: kafkaProducerClient,
		TokenTransferTopic:  tokenTransferTopic,
		AmplitudeEventTopic: amplitudeEventTopic,
		SqsQueueWrapper:     &SqsQueueWrapper,
		SQSDelaySeconds:     sqsDelaySeconds,
		CacheW:              cacheW,
	}
}

func (svc *OrderExecutionSvc) SetRequiredServices(services RequiredServices) {
	svc.services = services
}

func (svc *OrderExecutionSvc) SetRequiredRepos(repos RequiredRepos) {
	svc.repos = repos
}

func (svc *OrderExecutionSvc) CheckUserEligible(ctx context.Context, userDetails types.UserDetailsInterface, nftTokenDetails types.NFTTransferDetailsInterface, floatAmount float64) (bool, error) {
	if !svc.services.OtpSvc.CheckEligibility(ctx, userDetails) {
		logger.Error("CheckEligibility: User not eligible.", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.VALIDATION_ERROR].ErrorCode, "user": userDetails.Id})
		return false, fmt.Errorf("user not eligible")
	}
	if match, err := svc.services.OtpSvc.MatchReloginPin(ctx, userDetails); !match {
		logger.Error("CheckEligibility: Login Pin not matched", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.VALIDATION_ERROR].ErrorCode, "user": userDetails.Id, "error": err})
		return false, fmt.Errorf("login Pin not matched")
	}
	if match := svc.services.OtpSvc.MatchTransferPin(ctx, userDetails); !match {
		logger.Error("CheckEligibility: transfer OTP not matched", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.VALIDATION_ERROR].ErrorCode, "user": userDetails.Id})
		return false, fmt.Errorf("transfer OTP does not match")
	}
	tokenBalance, err := svc.services.PortfolioSvc.GetUserBalance(ctx, nftTokenDetails.NftId, nftTokenDetails.NetworkId, userDetails.Id)
	if err != nil {
		logger.Error("ExecuteFlow: can't retrieve user balance from portfolio svc", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "error": err})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, err.Error())
		return false, err
	}
	floatBalance, err := strconv.ParseFloat(tokenBalance.Result.Rows[0].Quantity, 64)
	if err != nil {
		logger.Error("ExecuteFlow:Unable to convert user balance to float", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "userId": userDetails.Id, "error": err.Error()})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, err.Error())
		return false, err
	}
	if floatAmount <= 0 {
		logger.Error("ExecuteFlow:User entered amount less than or equal to 0", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.QUANTITY_ERROR].ErrorCode, "userId": userDetails.Id})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.QUANTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.QUANTITY_ERROR].ErrorCode, "Amount cannot be zero")
		return false, fmt.Errorf("0 Amount error")
	}
	if tokenBalance.Result.Rows[0].EntityId == nftTokenDetails.NftId && floatBalance >= floatAmount {
		return true, nil
	} else {
		logger.Error("User Not eligible for transfer", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.QUANTITY_ERROR].ErrorCode, "portfolioResponse": tokenBalance})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.QUANTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.QUANTITY_ERROR].ErrorCode, "user Not eligible for transfer")
		return false, fmt.Errorf("user Not eligible for transfer")
	}
}

func (svc *OrderExecutionSvc) ExecuteOrder(ctx context.Context, userDetails types.UserDetailsInterface, nftTokenDetails types.NFTTransferDetailsInterface) (types.FMExecuteResult, error) {

	// Converting and checking if valid int quantity sent by user
	floatAmount, err := strconv.ParseFloat(nftTokenDetails.Amount, 64)
	if err != nil {
		logger.Error("ExecuteFlow:Unable to convert requested amount to float", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "userId": userDetails.Id, "error": err.Error()})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "Invalid Quantity To Transfer")
		return types.FMExecuteResult{}, err
	}
	integerAmount := int(floatAmount)
	if floatAmount != float64(int(floatAmount)) {
		logger.Error("ExecuteFlow:User entered amount in decimals", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.QUANTITY_ERROR].ErrorCode, "userId": userDetails.Id})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.QUANTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.QUANTITY_ERROR].ErrorCode, "Invalid Quantity To Transfer")
		return types.FMExecuteResult{}, fmt.Errorf("decimal Amount error")
	}
	// Check user eligibility based on balance
	_, err = svc.CheckUserEligible(ctx, userDetails, nftTokenDetails, floatAmount)
	if err != nil {
		return types.FMExecuteResult{}, err
	}

	// Create waitGroup and channels to execute parallel calls
	var wg sync.WaitGroup
	signingServiceChan := make(chan types.SigningServiceChannel)
	networkDetailsChan := make(chan types.NetworkChannel)
	nftDetailsChan := make(chan types.DQLChannel)
	wg.Add(3)
	// calling functions using go routines
	go svc.signingServiceInteraction(&wg, signingServiceChan, ctx, userDetails, nftTokenDetails)
	go svc.DqlNetworkInteraction(&wg, networkDetailsChan, ctx, nftTokenDetails.NetworkId)
	go svc.DQLNftInteraction(&wg, nftDetailsChan, ctx, nftTokenDetails.NftId)
	// Creating listeners for our channels

	signingServiceResult := <-signingServiceChan
	networkResult := <-networkDetailsChan
	nftDetailsResult := <-nftDetailsChan
	// Waiting for all go routines to finish
	wg.Wait()
	// close all channels to enable reading data
	close(signingServiceChan)
	close(networkDetailsChan)
	close(nftDetailsChan)

	// Process each channel data

	err = signingServiceResult.Error
	if err != nil {
		logger.Error("ExecuteFlow: Error Fetching user address", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "nftDetails": nftTokenDetails, "error": err.Error()})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, err.Error())
		return types.FMExecuteResult{}, err
	}
	userDetails.UserWalletAddress = signingServiceResult.Result.Address

	var networkName string

	err = networkResult.Error
	if err != nil {
		logger.Error("ExecuteFlow: Error Fetching network details", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "nftDetails": nftTokenDetails, "error": err.Error()})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, err.Error())
		return types.FMExecuteResult{}, err
	}
	networkName = networkResult.Result

	var nftDetails types.DQLEntityResponseInterface

	err = nftDetailsResult.Error
	if err != nil {
		logger.Error("ExecuteFlow: Error while querying data from DQL", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "DQLResponse": nftDetailsResult.Result, "request": nftTokenDetails})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, err.Error())
		return types.FMExecuteResult{}, err
	}
	nftDetails = nftDetailsResult.Result

	// Run validations on data received
	if userDetails.UserWalletAddress == nftTokenDetails.DestinationWalletAddress {
		logger.Error("ExecuteFlow: Self Transfer detected", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.SELF_TRANSFER_ERROR].ErrorCode, "request": nftTokenDetails, "jobId": nftTokenDetails.OrderId, "userId": userDetails.Id})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.SELF_TRANSFER_ERROR].HttpStatus, constants.ERROR_TYPES[constants.SELF_TRANSFER_ERROR].ErrorCode, "Cannot Transfer to self")
		return types.FMExecuteResult{}, fmt.Errorf("Self Transfer detected")
	}

	if nftDetails.Details.NetworkId != nftTokenDetails.NetworkId {
		logger.Error("ExecuteFlow: Network ID mismatch", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "DQLResponse": nftDetails, "request": nftTokenDetails})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "network ID mismatch")
		return types.FMExecuteResult{}, fmt.Errorf("network ID mismatch")
	}
	if nftDetails.Details.ErcType != nftTokenDetails.ErcType {
		logger.Error("ExecuteFlow: ErcType mismatch", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "DQLResponse": nftDetails, "request": nftTokenDetails})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "ErcType mismatch")
		return types.FMExecuteResult{}, fmt.Errorf("ErcType mismatch")
	}

	orderMetaDataObj := svc.prepareOrderMetaData(userDetails, nftTokenDetails, nftDetails)
	orderId := orderMetaDataObj.OrderId
	savedOrder := svc.repos.OrderMetaDataRepo.InsertOrderMetadata(&orderMetaDataObj)
	if savedOrder.Error != nil {
		logger.Error("ExecuteFlow: Error saving orderMetaData", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "error": savedOrder.Error.Error(), "orderData": orderMetaDataObj})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.PROCESS_ERROR].HttpStatus, constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "error creating orderMetaData")
		return types.FMExecuteResult{}, fmt.Errorf("error creating orderMetaData")
	}
	//TODO: Pass repo on init and send object to repo

	var payloadForFM types.FMExecuteRequest
	payloadForFM.UserId = userDetails.Id
	payloadForFM.FlowType = constants.FM_FLOW_TYPE
	payloadForFM.Operation = constants.FM_EXECUTE_OPERATION
	payloadForFM.JobId = orderId
	payloadForFM.Payload = types.ExecutePayload{
		Amount:             strconv.Itoa(integerAmount),
		NftContractAddress: nftDetails.Details.Address,
		NftId:              nftDetails.Details.TokenId,
		NftType:            nftDetails.Details.ErcType,
		NetworkId:          nftDetails.Details.NetworkId,
		SenderAddress:      userDetails.UserWalletAddress,
		RecepientAddress:   nftTokenDetails.DestinationWalletAddress,
		IsGsnRequired:      nftTokenDetails.IsGsnRequired,
		Deadline:           fmt.Sprintf("%d", time.Now().Add(time.Duration(appconfig.FM_EXECUTE_DEADLINE)*time.Minute).Unix()),
	}
	if nftTokenDetails.IsGsnRequired {
		payloadForFM.Payload.GsnIncludeMaxAmount = nftTokenDetails.GsnIncludeMaxAmount
		payloadForFM.Payload.GsnIncludenetworkId = nftTokenDetails.GsnIncludeNetworkId
		if len(nftTokenDetails.GsnIncludeToken) > 0 {
			payloadForFM.Payload.GsnIncludeToken = nftTokenDetails.GsnIncludeToken
			dqlResponse, err := svc.services.DQLSvc.GetTokenByAddress(ctx, payloadForFM.Payload.GsnIncludeToken, payloadForFM.Payload.GsnIncludenetworkId)
			if err != nil || len(dqlResponse.Entities) == 0 {
				logger.Error("Error in Fetching DQL response for Execute", dqlResponse, err)
			} else {
				decimals := dqlResponse.Entities[0].Details.Decimals
				payloadForFM.Payload.GsnIncludeMaxAmount = multiplyByDecimals(nftTokenDetails.GsnIncludeMaxAmount, decimals)
			}
		} else {
			networkDetails, err := svc.services.DQLSvc.GetEntityById(ctx, nftTokenDetails.GsnIncludeNetworkId, false)
			if err != nil {
				logger.Error("OrderExecutionService: error while getting GSN network details", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "error": err, "networkId": nftTokenDetails.GsnIncludeNetworkId})
				apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "Could not fetch network details")
				return types.FMExecuteResult{}, err
			}
			if (networkDetails.Details == types.DQLDetails{}) {
				logger.Error("OrderExecutionService | Empty network Details", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "network": nftTokenDetails.GsnIncludeNetworkId})
				apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "Downstream Error: GSN Network info not found")
				return types.FMExecuteResult{}, fmt.Errorf("OrderEstimationService | Network Not found")
			}
			tokenDetails, err := svc.services.DQLSvc.GetEntityById(ctx, networkDetails.Details.NativeCurrencyId, false)
			if err != nil {
				logger.Error("OrderExecutionService: error while getting GSN Token details", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "error": err, "tokenDetails": networkDetails.Details.NativeCurrencyId})
				apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "Could not fetch native token details")
				return types.FMExecuteResult{}, err
			}
			decimals := tokenDetails.Details.Decimals
			payloadForFM.Payload.GsnIncludeMaxAmount = multiplyByDecimals(nftTokenDetails.GsnIncludeMaxAmount, decimals)
			payloadForFM.Payload.GsnIncludeToken = ""
		}
	}
	logger.Info("Payload for Execute", map[string]interface{}{"context": ctx, "payload": payloadForFM})
	responseFromFM, err := svc.services.FlowManagerSvc.ExecuteOrder(ctx, payloadForFM)
	if err != nil {
		dbResponse := svc.repos.OrderMetaDataRepo.UpdateOrderStatusByOrderId(orderId, constants.FAILED)
		if dbResponse.Error != nil {
			logger.Error("ExecuteFlow: !Immediate attention! failed to update Order Status. ", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "Response": dbResponse.Error.Error()})
		}
		errMsg := "Error executing transfer for NFT"
		fmErrorMsgSet := false
		if (responseFromFM.Error != types.Error{} && responseFromFM.Error.ErrorCode != "") {
			if msg, ok := constants.FE_ERROR_CODES_MAPPING[responseFromFM.Error.ErrorCode]; ok {
				errMsg = fmt.Sprintf("%s:%s", errMsg, msg.Message)
				fmErrorMsgSet = true
			}
		}
		if (!fmErrorMsgSet && responseFromFM.Error != types.Error{} && responseFromFM.Error.Message != "") {
			errMsg = fmt.Sprintf("%s:%s", errMsg, responseFromFM.Error.Message)
		}
		logger.Error("ExecuteFlow: Call to FM failed", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "FMResponse": responseFromFM, "request": nftTokenDetails})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, errMsg)
		return types.FMExecuteResult{}, err
	}

	logger.Info(`OrderExecutionService response from flow manager:`, map[string]interface{}{"responseFromFM": responseFromFM})
	err = svc.NotifyTokenTransfer(orderId, userDetails.Id, nftTokenDetails.NftId, nftDetails.Details.ErcType, constants.NFT_TRANSFER_PURPOSE, nftTokenDetails.NetworkId, "")
	if err != nil {
		logger.Error("Error Notifying PS", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "user": userDetails.Id, "orderId": orderId})
	}
	amplitudeEvent := types.AmplitudeTransferEventInterface{
		AppName:          constants.AMPLITUDE_APP_NAME,
		Network:          nftTokenDetails.NetworkId,
		Token:            nftDetails.Details.Name,
		NftType:          "",
		TokenId:          nftTokenDetails.NftId,
		ErcType:          nftDetails.Details.ErcType,
		Type:             constants.AMPLITUDE_TRANSFER_EVENT_TYPE,
		OrderId:          orderId,
		Status:           constants.AMPLITUDE_STATUS_MAPPING[orderMetaDataObj.Status],
		DeviceId:         userDetails.DeviceId,
		UserId:           userDetails.Id,
		DeviceType:       userDetails.Source,
		ReceiversAddress: nftTokenDetails.DestinationWalletAddress,
		TokenCount:       integerAmount,
		CollectionId:     nftDetails.Details.CollectionId,
		CollectionName:   nftDetails.Details.CollectionName,
		Chain:            networkName,
	}
	err = svc.NotifyAmplitudeEvent(orderId, userDetails.Id, amplitudeEvent)
	if err != nil {
		logger.Error("Error Notifying Amplitude Service", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "user": userDetails.Id, "orderId": orderId})
	}

	message := types.SQSJobInterface{
		JobId: orderId,
	}
	err = svc.SqsQueueWrapper.PublishMessage(ctx, message, orderId, svc.SQSDelaySeconds)
	if err != nil {
		logger.Error("Error publishing message to SQS: %s", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "error": err.Error()})
	}

	return responseFromFM, nil
}

func (svc *OrderExecutionSvc) prepareOrderMetaData(userDetails types.UserDetailsInterface, nftTokenDetails types.NFTTransferDetailsInterface, nftDetails types.DQLEntityResponseInterface) models.OrderMetadata {
	userId, _ := uuid.Parse(userDetails.Id)
	networkId, _ := uuid.Parse(nftTokenDetails.NetworkId)
	return models.OrderMetadata{
		OrderId:           uuid.New().String(),
		UserId:            userId.String(),
		NetworkId:         networkId.String(),
		Status:            constants.CREATED,
		EntityType:        strings.ToUpper(nftDetails.Details.ErcType),
		EntityAddress:     nftDetails.Details.Address,
		NftId:             nftTokenDetails.NftId,
		Count:             nftTokenDetails.Amount,
		OrderType:         nftTokenDetails.TransferType,
		Slippage:          appconfig.SLIPPAGE,
		ExecutionResponse: "",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
}

func (svc *OrderExecutionSvc) NotifyTokenTransfer(jobId, userId, entityId, entityType, orderType, networkId, vendorId string) error {
	transferNotification := types.TransferNotificationKafkaDataInterface{
		JobId:  jobId,
		UserId: userId,
		//EntityId:   entityId,
		//EntityType: entityType,
		OrderType: orderType,
		NetworkId: networkId,
		VendorId:  vendorId,
	}
	return svc.KafkaProducerClient.Produce(svc.TokenTransferTopic, jobId, transferNotification)
}

func (svc *OrderExecutionSvc) NotifyAmplitudeEvent(orderId string, userId string, event types.AmplitudeTransferEventInterface) error {
	amplitudeNotification := types.AmplitudeEventInterface{
		UserId:          userId,
		Eventtype:       constants.AMPLITUDE_EVENT_NAME,
		EventProperties: event,
	}
	return svc.KafkaProducerClient.Produce(svc.AmplitudeEventTopic, orderId, amplitudeNotification)
}

func (svc *OrderExecutionSvc) signingServiceInteraction(wg *sync.WaitGroup, channel chan types.SigningServiceChannel, ctx context.Context, userDetails types.UserDetailsInterface, nftTokenDetails types.NFTTransferDetailsInterface) {
	defer wg.Done()
	signingSvcResponse, err := svc.services.SigningSvc.GetUserWalletAddress(ctx, userDetails.Id, nftTokenDetails.NetworkId)
	channel <- types.SigningServiceChannel{
		Result: signingSvcResponse,
		Error:  err,
	}
}

func (svc *OrderExecutionSvc) DqlNetworkInteraction(wg *sync.WaitGroup, channel chan types.NetworkChannel, ctx context.Context, identifier string) {
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

func (svc *OrderExecutionSvc) DQLNftInteraction(wg *sync.WaitGroup, channel chan types.DQLChannel, ctx context.Context, identifier string) {
	defer wg.Done()
	dqlResponse, err := svc.services.DQLSvc.GetEntityById(ctx, identifier, true)
	channel <- types.DQLChannel{
		Result: dqlResponse,
		Error:  err,
	}
}
