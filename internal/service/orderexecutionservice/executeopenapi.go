package orderexecutionservice

import (
	"context"
	"fmt"
	apiLogger "pismo-dev/api/logger"
	apiTypes "pismo-dev/api/types"
	"pismo-dev/constants"
	"pismo-dev/internal/appconfig"
	"pismo-dev/internal/models"
	"pismo-dev/internal/types"
	"pismo-dev/pkg/logger"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

func (svc *OrderExecutionSvc) OpenApiExecuteOrder(ctx context.Context, orderDetails apiTypes.OpenApiExecuteRequest) (types.FMExecuteResult, error) {

	// Converting and checking if valid int quantity sent by user
	floatAmount, err := strconv.ParseFloat(orderDetails.Quantity, 64)
	if err != nil {
		logger.Error("OpenAPIExecuteFlow:Unable to convert requested amount to float", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "userId": orderDetails.UserId, "error": err.Error()})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "Invalid Quantity To Transfer")
		return types.FMExecuteResult{}, err
	}
	integerAmount := int(floatAmount)
	if floatAmount != float64(int(floatAmount)) {
		logger.Error("OpenAPIExecuteFlow:User entered amount in decimals", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.QUANTITY_ERROR].ErrorCode, "userId": orderDetails.UserId})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.QUANTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.QUANTITY_ERROR].ErrorCode, "Invalid Quantity To Transfer")
		return types.FMExecuteResult{}, fmt.Errorf("decimal Amount error")
	}
	if integerAmount == 0 {
		logger.Error("OpenAPIMintFlow:User entered 0 amount", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.QUANTITY_ERROR].ErrorCode, "userId": orderDetails.UserId})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.QUANTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.QUANTITY_ERROR].ErrorCode, "Invalid Quantity To Transfer")
		return types.FMExecuteResult{}, fmt.Errorf("0 Amount error")
	}
	_, err = svc.CheckOpenApiUserEligible(ctx, orderDetails, floatAmount)
	if err != nil {
		return types.FMExecuteResult{}, err
	}

	var wg sync.WaitGroup

	networkDetailsChan := make(chan types.NetworkChannel)
	nftDetailsChan := make(chan types.DQLChannel)
	wg.Add(2)
	// calling functions using go routines
	go svc.DqlNetworkInteraction(&wg, networkDetailsChan, ctx, orderDetails.NetworkId)
	go svc.DQLNftInteraction(&wg, nftDetailsChan, ctx, orderDetails.NFTTokenId)
	// Creating listeners for our channels

	networkResult := <-networkDetailsChan
	nftDetailsResult := <-nftDetailsChan
	// Waiting for all go routines to finish
	wg.Wait()
	// close all channels to enable reading data

	close(networkDetailsChan)
	close(nftDetailsChan)

	// Process each channel data

	var networkName string

	err = networkResult.Error
	if err != nil {
		logger.Error("ExecuteFlow: Error Fetching network details", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "orderDetails": orderDetails, "error": err.Error()})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, err.Error())
		return types.FMExecuteResult{}, err
	}
	networkName = networkResult.Result

	var nftDetails types.DQLEntityResponseInterface

	err = nftDetailsResult.Error
	if err != nil {
		logger.Error("ExecuteFlow: Error while querying data from DQL", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "DQLResponse": nftDetailsResult.Result, "request": orderDetails})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, err.Error())
		return types.FMExecuteResult{}, err
	}
	nftDetails = nftDetailsResult.Result

	orderMetaDataObj := svc.prepareOpenAPIOrderMetaData(orderDetails)
	orderId := orderMetaDataObj.OrderId
	savedOrder := svc.repos.OrderMetaDataRepo.InsertOrderMetadata(&orderMetaDataObj)
	if savedOrder.Error != nil {
		logger.Error("OpenAPIExecuteFlow: Error saving orderMetaData", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "error": savedOrder.Error.Error(), "orderData": orderMetaDataObj})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.PROCESS_ERROR].HttpStatus, constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "error creating orderMetaData")
		return types.FMExecuteResult{}, fmt.Errorf("error creating orderMetaData")
	}
	//TODO: Pass repo on init and send object to repo

	var payloadForFM types.OpenAPINFTOrder
	payloadForFM.UserId = orderDetails.UserId
	payloadForFM.NetworkID = orderDetails.NetworkId
	payloadForFM.JobId = orderId
	payloadForFM.Sponsor = orderDetails.IsSponsored
	payloadForFM.Payload = types.OpenAPIExecutePayload{
		Quantity:           strconv.Itoa(integerAmount),
		NftContractAddress: orderDetails.CollectionAddress,
		NftId:              orderDetails.NFTId,
		NftType:            orderDetails.ErcType,
		NetworkId:          orderDetails.NetworkId,
		SenderAddress:      orderDetails.SenderAddress,
		RecepientAddress:   orderDetails.RecipientWalletAddress,
	}
	payloadForFM.Deadline = time.Now().Add(time.Duration(appconfig.FM_EXECUTE_DEADLINE) * time.Minute).Unix()
	// if orderDetails.IsGsnRequired {
	// 	payloadForFM.Payload.GsnIncludeMaxAmount = orderDetails.GsnIncludeMaxAmount
	// 	payloadForFM.Payload.GsnIncludenetworkId = orderDetails.GsnIncludeNetworkId
	// 	if len(orderDetails.GsnIncludeToken) > 0 {
	// 		payloadForFM.Payload.GsnIncludeToken = orderDetails.GsnIncludeToken
	// 		dqlResponse, err := svc.services.DQLSvc.GetTokenByAddress(ctx, payloadForFM.Payload.GsnIncludeToken, payloadForFM.Payload.GsnIncludenetworkId)
	// 		if err != nil || len(dqlResponse.Entities) == 0 {
	// 			logger.Error("Error in Fetching DQL response for Execute", dqlResponse, err)
	// 		} else {
	// 			decimals := dqlResponse.Entities[0].Details.Decimals
	// 			payloadForFM.Payload.GsnIncludeMaxAmount = multiplyByDecimals(orderDetails.GsnIncludeMaxAmount, decimals)
	// 		}
	// 	} else {
	// 		networkDetails, err := svc.services.DQLSvc.GetEntityById(ctx, orderDetails.GsnIncludeNetworkId, false)
	// 		if err != nil {
	// 			logger.Error("OrderExecutionService: error while getting GSN network details", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "error": err, "networkId": orderDetails.GsnIncludeNetworkId})
	// 			apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "Could not fetch network details")
	// 			return types.FMExecuteResult{}, err
	// 		}
	// 		if (networkDetails.Details == types.DQLDetails{}) {
	// 			logger.Error("OrderExecutionService | Empty network Details", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "network": orderDetails.GsnIncludeNetworkId})
	// 			apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "Downstream Error: GSN Network info not found")
	// 			return types.FMExecuteResult{}, fmt.Errorf("OrderEstimationService | Network Not found")
	// 		}
	// 		tokenDetails, err := svc.services.DQLSvc.GetEntityById(ctx, networkDetails.Details.NativeCurrencyId, false)
	// 		if err != nil {
	// 			logger.Error("OrderExecutionService: error while getting GSN Token details", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "error": err, "tokenDetails": networkDetails.Details.NativeCurrencyId})
	// 			apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "Could not fetch native token details")
	// 			return types.FMExecuteResult{}, err
	// 		}
	// 		decimals := tokenDetails.Details.Decimals
	// 		payloadForFM.Payload.GsnIncludeMaxAmount = multiplyByDecimals(orderDetails.GsnIncludeMaxAmount, decimals)
	// 		payloadForFM.Payload.GsnIncludeToken = ""
	// 	}
	// }
	logger.Info("Payload for Execute", map[string]interface{}{"context": ctx, "payload": payloadForFM})
	responseFromFM, err := svc.services.FlowManagerSvc.ExecuteOpenAPIOrder(ctx, payloadForFM, string(constants.NFT_TRANSFER_ORDER), orderDetails.VendorId)
	if err != nil {
		dbResponse := svc.repos.OrderMetaDataRepo.UpdateOrderStatusByOrderId(orderId, constants.FAILED)
		if dbResponse.Error != nil {
			logger.Error("OpenAPIExecuteFlow: !Immediate attention! failed to update Order Status. ", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "Response": dbResponse.Error.Error()})
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
		logger.Error("OpenAPIExecuteFlow: Call to FM failed", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "FMResponse": responseFromFM, "request": orderDetails})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, errMsg)
		return types.FMExecuteResult{}, err
	}

	logger.Info(`OrderExecutionService response from flow manager:`, map[string]interface{}{"responseFromFM": responseFromFM})
	err = svc.NotifyTokenTransfer(orderId, orderDetails.UserId, orderDetails.NFTId, orderDetails.ErcType, constants.NFT_TRANSFER_PURPOSE, orderDetails.NetworkId, orderDetails.VendorId)
	if err != nil {
		logger.Error("Error Notifying PS", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "user": orderDetails.UserId, "orderId": orderId})
	}
	amplitudeEvent := types.AmplitudeTransferEventInterface{
		AppName:          constants.AMPLITUDE_APP_NAME,
		Network:          orderDetails.NetworkId,
		Token:            nftDetails.Details.Name,
		NftType:          "",
		TokenId:          orderDetails.NFTId,
		ErcType:          orderDetails.ErcType,
		Type:             constants.AMPLITUDE_TRANSFER_EVENT_TYPE,
		OrderId:          orderId,
		Status:           constants.AMPLITUDE_STATUS_MAPPING[orderMetaDataObj.Status],
		DeviceId:         "",
		UserId:           orderDetails.UserId,
		DeviceType:       "",
		ReceiversAddress: orderDetails.RecipientWalletAddress,
		TokenCount:       integerAmount,
		CollectionId:     orderDetails.CollectionId,
		CollectionName:   nftDetails.Details.CollectionName,
		Chain:            networkName,
		VendorId:         orderDetails.VendorId,
	}
	err = svc.NotifyAmplitudeEvent(orderId, orderDetails.UserId, amplitudeEvent)
	if err != nil {
		logger.Error("Error Notifying Amplitude Service OpenAPI", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "user": orderDetails.UserId, "orderId": orderId})
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

func (svc *OrderExecutionSvc) CheckOpenApiUserEligible(ctx context.Context, nftTokenDetails apiTypes.OpenApiExecuteRequest, floatAmount float64) (bool, error) {
	tokenBalance, err := svc.services.PortfolioSvc.GetUserBalance(ctx, nftTokenDetails.NFTTokenId, nftTokenDetails.NetworkId, nftTokenDetails.UserId)
	if err != nil {
		logger.Error("Unable to get user balance", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "userId": nftTokenDetails.UserId})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "Insufficient User balance")
		return false, err
	}
	floatBalance, err := strconv.ParseFloat(tokenBalance.Result.Rows[0].Quantity, 64)
	if err != nil {
		logger.Error("Unable to convert user balance to float", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "userId": nftTokenDetails.UserId, "error": err.Error()})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, err.Error())
		return false, err
	}
	if floatAmount <= 0 {
		logger.Error("User entered amount less than or equal to 0", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.QUANTITY_ERROR].ErrorCode, "userId": nftTokenDetails.UserId})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.QUANTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.QUANTITY_ERROR].ErrorCode, "Amount cannot be zero")
		return false, fmt.Errorf("0 Amount error")
	}
	if tokenBalance.Result.Rows[0].EntityId == nftTokenDetails.NFTTokenId && floatBalance >= floatAmount {
		return true, nil //TODO: Check locked balance logic with portfolio svc
	} else {
		logger.Error("User Not eligible for transfer", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.QUANTITY_ERROR].ErrorCode, "portfolioResponse": tokenBalance})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.QUANTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.QUANTITY_ERROR].ErrorCode, "User Not eligible for transfer")
		return false, fmt.Errorf("user Not eligible for transfer")
	}
}

func (svc *OrderExecutionSvc) OpenApiMintOrder(ctx context.Context, orderDetails apiTypes.OpenAPINFTMintOrder) (types.FMExecuteResult, error) {

	// Converting and checking if valid int quantity sent by user
	floatAmount, err := strconv.ParseFloat(orderDetails.Quantity, 64)
	if err != nil {
		logger.Error("OpenAPIMintFlow:Unable to convert requested amount to float", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "userId": orderDetails.UserId, "error": err.Error()})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "Invalid Quantity To Transfer")
		return types.FMExecuteResult{}, err
	}
	integerAmount := int(floatAmount)
	if floatAmount != float64(int(floatAmount)) {
		logger.Error("OpenAPIMintFlow:User entered amount in decimals", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.QUANTITY_ERROR].ErrorCode, "userId": orderDetails.UserId})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.QUANTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.QUANTITY_ERROR].ErrorCode, "Invalid Quantity To Transfer")
		return types.FMExecuteResult{}, fmt.Errorf("decimal Amount error")
	}
	if integerAmount == 0 {
		logger.Error("OpenAPIMintFlow:User entered 0 amount", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.QUANTITY_ERROR].ErrorCode, "userId": orderDetails.UserId})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.QUANTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.QUANTITY_ERROR].ErrorCode, "Invalid Quantity To Transfer")
		return types.FMExecuteResult{}, fmt.Errorf("0 Amount error")
	}

	var wg sync.WaitGroup

	networkChannel := make(chan types.NetworkChannel)
	collectionDetailsChannel := make(chan types.CollectionDqlChannel)
	wg.Add(2)
	// calling functions on different go routines

	go svc.DqlNetworkInteraction(&wg, networkChannel, ctx, orderDetails.NetworkId)
	go svc.services.OgmintSvc.DqlCollectionInteraction(&wg, collectionDetailsChannel, ctx, orderDetails.CollectionId)
	// Creating listeners for our channels

	networkResult := <-networkChannel
	collectionDetailsResult := <-collectionDetailsChannel
	// Waiting for all go routines to finish
	wg.Wait()
	// close all channels to enable reading data

	close(networkChannel)
	close(collectionDetailsChannel)
	// Process each channel data

	var networkName string

	err = networkResult.Error
	if err != nil {
		logger.Error("ExecuteFlow: Error Fetching network details openAPI", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "networkId": orderDetails.NetworkId, "error": err.Error()})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, err.Error())
		return types.FMExecuteResult{}, fmt.Errorf("DQL Error")
	}
	networkName = networkResult.Result

	var collectionDetails types.DQLNftCollectionResponseInterface

	err = collectionDetailsResult.Error
	collectionDetails = collectionDetailsResult.Result
	if err != nil || (collectionDetails.Details == types.DQLCollectionDetails{}) || (collectionDetails.Details.ContractMetadata == types.DQLCollectionMetaData{}) {
		logger.Error("OGMintFlow:Error fetching collection details openAPI", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "collectionId": orderDetails.CollectionId})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "Internal Error")
		return types.FMExecuteResult{}, fmt.Errorf("DQL Error")
	}

	orderMetaDataObj := svc.prepareOpenAPIMintOrderMetaData(orderDetails)
	orderId := orderMetaDataObj.OrderId
	savedOrder := svc.repos.OrderMetaDataRepo.InsertOrderMetadata(&orderMetaDataObj)
	if savedOrder.Error != nil {
		logger.Error("OpenAPIMintFlow: Error saving orderMetaData", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "error": savedOrder.Error.Error(), "orderData": orderMetaDataObj})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.PROCESS_ERROR].HttpStatus, constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "error creating orderMetaData")
		return types.FMExecuteResult{}, fmt.Errorf("error creating orderMetaData")
	}
	//TODO: Pass repo on init and send object to repo
	var payloadForFM types.OpenAPINFTOrder
	payloadForFM.UserId = orderDetails.UserId
	payloadForFM.NetworkID = orderDetails.NetworkId
	payloadForFM.JobId = orderId
	payloadForFM.Sponsor = orderDetails.IsSponsored
	payloadData := types.OpenAPIMintPayload{
		Quantity:           strconv.Itoa(integerAmount),
		NftContractAddress: orderDetails.CollectionAddress,
		NftType:            orderDetails.ErcType,
		NetworkId:          orderDetails.NetworkId,
		SenderAddress:      orderDetails.SenderAddress,
		RecepientAddress:   orderDetails.RecipientAddress,
		MetaData: types.NftMetaData{
			Description:    orderDetails.MetaData.Description,
			CollectionName: orderDetails.CollectionName,
			NftName:        orderDetails.MetaData.NftName,
			Uri:            orderDetails.MetaData.Uri,
			Properties:     orderDetails.MetaData.Properties,
		},
	}

	payloadForFM.Payload = payloadData
	payloadForFM.Deadline = time.Now().Add(time.Duration(appconfig.FM_EXECUTE_DEADLINE) * time.Minute).Unix()
	logger.Info("Payload for Execute", map[string]interface{}{"context": ctx, "payload": payloadForFM})
	responseFromFM, err := svc.services.FlowManagerSvc.ExecuteOpenAPIOrder(ctx, payloadForFM, string(constants.NFT_MINT_ORDER), orderDetails.VendorId)
	if err != nil {
		dbResponse := svc.repos.OrderMetaDataRepo.UpdateOrderStatusByOrderId(orderId, constants.FAILED)
		if dbResponse.Error != nil {
			logger.Error("OpenAPIMintFlow: !Immediate attention! failed to update Order Status. ", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "Response": dbResponse.Error.Error()})
		}
		errMsg := "Error executing Mint for NFT"
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
		logger.Error("OpenAPIMintFlow: Call to FM failed", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "FMResponse": responseFromFM, "request": orderDetails})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, errMsg)
		return types.FMExecuteResult{}, err
	}

	logger.Info(`OrderExecutionService response from flow manager:`, map[string]interface{}{"responseFromFM": responseFromFM})
	err = svc.NotifyTokenTransfer(orderId, orderDetails.UserId, "", orderDetails.ErcType, constants.NFT_MINT_PURPOSE, orderDetails.NetworkId, orderDetails.VendorId)
	if err != nil {
		logger.Error("Error Notifying PS", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "user": orderDetails.UserId, "orderId": orderId})
	}

	message := types.SQSJobInterface{
		JobId: orderId,
	}
	err = svc.SqsQueueWrapper.PublishMessage(ctx, message, orderId, svc.SQSDelaySeconds)
	if err != nil {
		logger.Error("Error publishing message to SQS: %s", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "error": err.Error()})
	}

	amplitudeEvent := types.AmplitudeMintEventInterface{
		AppName:          constants.AMPLITUDE_APP_NAME,
		Network:          orderDetails.NetworkId,
		NftType:          strings.ToUpper(orderDetails.ErcType),
		TokenId:          "",
		Type:             constants.AMPLITUDE_MINT_EVENT_TYPE,
		Status:           constants.AMPLITUDE_STATUS_MAPPING[orderMetaDataObj.Status],
		Product:          "nft",
		CollectionName:   collectionDetails.Details.CollectionName,
		CollectionId:     orderDetails.CollectionId,
		TokenCount:       1,
		ReceiversAddress: orderDetails.RecipientAddress,
		Chain:            networkName,
		VendorId:         orderDetails.VendorId,
	}
	err = svc.services.OgmintSvc.NotifyAmplitudeEvent(orderId, orderDetails.UserId, amplitudeEvent)
	if err != nil {
		logger.Error("Error Notifying Amplitude Service", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "user": orderDetails.UserId, "orderId": orderId})
	}

	return responseFromFM, nil
}

func (svc *OrderExecutionSvc) prepareOpenAPIOrderMetaData(orderDetails apiTypes.OpenApiExecuteRequest) models.OrderMetadata {
	userId, _ := uuid.Parse(orderDetails.UserId)
	networkId, _ := uuid.Parse(orderDetails.NetworkId)
	return models.OrderMetadata{
		OrderId:           uuid.New().String(),
		UserId:            userId.String(),
		NetworkId:         networkId.String(),
		Status:            constants.CREATED,
		VendorId:          orderDetails.VendorId,
		EntityType:        strings.ToUpper(orderDetails.ErcType),
		EntityAddress:     orderDetails.CollectionAddress,
		NftId:             orderDetails.NFTId,
		Count:             orderDetails.Quantity,
		OrderType:         orderDetails.OperationType,
		Slippage:          appconfig.SLIPPAGE,
		ExecutionResponse: "",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
}

func (svc *OrderExecutionSvc) prepareOpenAPIMintOrderMetaData(orderDetails apiTypes.OpenAPINFTMintOrder) models.OrderMetadata {
	userId, _ := uuid.Parse(orderDetails.UserId)
	networkId, _ := uuid.Parse(orderDetails.NetworkId)
	return models.OrderMetadata{
		OrderId:           uuid.New().String(),
		UserId:            userId.String(),
		NetworkId:         networkId.String(),
		Status:            constants.CREATED,
		VendorId:          orderDetails.VendorId,
		EntityType:        strings.ToUpper(orderDetails.ErcType),
		EntityAddress:     orderDetails.CollectionAddress,
		Count:             orderDetails.Quantity,
		OrderType:         orderDetails.OperationType,
		Slippage:          appconfig.SLIPPAGE,
		ExecutionResponse: "",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
}
