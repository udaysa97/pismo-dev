package ogmintservice

import (
	"context"
	"encoding/json"
	"fmt"
	"pismo-dev/constants"
	"pismo-dev/internal/models"
	"pismo-dev/internal/types"
	"pismo-dev/pkg/logger"
	"time"
)

func (svc *OGMintSvc) InitiateSqsConsumer(key string) {
	queueWrapper := svc.Queues[key].WrapperObj
	queueProps := svc.Queues[key].QueueProps
	queueWrapper.InitiateQueueConsumption(queueProps.SQSVisibilityTimeout, queueProps.SQSWaitTime, queueProps.SQSMaxMessage, svc.ProcessSQSMessage, true)
}

func (svc *OGMintSvc) ProcessSQSMessage(message string) bool {
	var sqsMessage types.OGMintQueueMessage
	ctx := context.TODO()
	err := json.Unmarshal([]byte(message), &sqsMessage)
	if err != nil {
		logger.Error("NFT Order processing error:", map[string]interface{}{"message": message, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode})
		return false
	}
	if len(sqsMessage.OrderId) <= 0 {
		logger.Error("NFT Order processing error: No orderId in message", map[string]interface{}{"message": message, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode})
		return true
	}
	logger.Info("Process-OGMINT: Started reconcile for jobId", map[string]interface{}{"orderId": sqsMessage.OrderId})
	return svc.processOrder(ctx, sqsMessage)
}

func (svc *OGMintSvc) processOrder(ctx context.Context, order types.OGMintQueueMessage) bool {
	orderDetails, err := svc.GetOrderDetails(order.OrderId)
	if err != nil {
		logger.Error("Process-OGMINT: Error fetching order", map[string]interface{}{"orderId": order.OrderId, "error": err.Error()})
		return svc.CreateQueueEntryFromExistingQueue(ctx, order, order.Status, true, order.RetryCount+1, order.TxHash, order.Type, "DBError: Could not get order details")
	}
	if order.RetryCount > svc.MintConfigs[constants.COLLECTION_VENDOR_MAPPING[order.VendorName]+order.NetworkId].MaxRetries {
		logger.Error("Process-OGMINT:Retry exceeded:Deleting Message", map[string]interface{}{"orderId": order.OrderId})
		errorMessage := "Retry Exceeded: " + order.ErrorMessage
		svc.repos.OrderMetaDataRepo.UpdateOrderStatusForFailedOrder(order.OrderId, constants.FAILED, errorMessage)
		amplitudeEvent := types.AmplitudeMintEventInterface{
			AppName:    constants.AMPLITUDE_APP_NAME,
			Network:    order.NetworkId,
			TokenId:    "",
			Type:       constants.AMPLITUDE_MINT_EVENT_TYPE,
			Status:     constants.AMPLITUDE_STATUS_MAPPING[constants.FAILED],
			Product:    constants.AMPLITUDE_PRODUCT_TYPE,
			TokenCount: 1,
		}
		svc.NotifyAmplitudeEvent(order.OrderId, order.UserId, amplitudeEvent)
		return true
	}
	if orderDetails.Status != order.Status {
		logger.Error("Process-OGMINT: Stale entry found", map[string]interface{}{"orderId": order.OrderId, "order": order, "orderDetails": orderDetails})
		return true
	}
	_, err = svc.CheckLockAllowed(svc.GetLockKeyAndTTL(order.NetworkId, order.Type, order.VendorName))
	if err != nil {
		//logger.Info("Process-OGMINT: Lock Not Acquired", map[string]interface{}{"orderId": order.OrderId, "error": err.Error()})
		return false
	}

	if order.Type == constants.MINT {
		return svc.Mint(ctx, order, orderDetails)

	} else {
		return svc.Status(ctx, order, orderDetails)
	}
}

func (svc *OGMintSvc) Mint(ctx context.Context, order types.OGMintQueueMessage, orderDetails models.OrderMetadata) bool {
	var txIdentifier string
	var err error
	if order.VendorName == constants.NFTPORTVENDOR {
		txHash, err := svc.MintNFTNftPort(order)
		if err != nil {
			logger.Error("Process-OGMINT:Error Calling MINT API in NftPort", map[string]interface{}{"orderId": order.OrderId, "error": err.Error(), "order": order})
			svc.CreateQueueEntryFromExistingQueue(ctx, order, order.Status, true, order.RetryCount+1, "", order.Type, err.Error())
			return true
		}
		txIdentifier = txHash
		svc.createNewTransaction(order.OrderId, txHash, order.OrderType)
	} else {
		txId, err := svc.MintNFTCrossMint(order)
		if err != nil {
			logger.Error("Process-OGMINT:Error Calling MINT API in CrossMint", map[string]interface{}{"orderId": order.OrderId, "error": err.Error(), "order": order})
			svc.CreateQueueEntryFromExistingQueue(ctx, order, order.Status, true, order.RetryCount+1, "", order.Type, err.Error())
			return true
		}
		txIdentifier = txId
	}
	orderDetails.Status = constants.RUNNING
	orderDetails.UpdatedAt = time.Now()
	orderDetails.RetryCount = order.RetryCount
	if err := svc.repos.OrderMetaDataRepo.UpdateOrderByOrderId(&orderDetails).Error; err != nil {
		// Need ALERTING
		logger.Error("Process-OGMINT::Could not update ordermetadata Table", map[string]interface{}{"orderId": order.OrderId, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "orderDetails": orderDetails})
	}
	if !svc.CreateQueueEntryFromExistingQueue(ctx, order, constants.RUNNING, false, 0, txIdentifier, constants.STATUSCHECK, "") {
		//TODO: ALERTING
	}
	amplitudeEvent := types.AmplitudeMintEventInterface{
		AppName:    constants.AMPLITUDE_APP_NAME,
		Network:    order.NetworkId,
		TokenId:    "",
		Type:       constants.AMPLITUDE_MINT_EVENT_TYPE,
		Status:     constants.AMPLITUDE_STATUS_MAPPING[orderDetails.Status],
		Product:    constants.AMPLITUDE_PRODUCT_TYPE,
		TokenCount: 1,
	}
	err = svc.NotifyAmplitudeEvent(order.OrderId, order.UserId, amplitudeEvent)
	if err != nil {
		logger.Error("Error Notifying Amplitude Service", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "user": order.UserId, "orderId": order.OrderId})
	}
	return true
}

func (svc *OGMintSvc) Status(ctx context.Context, order types.OGMintQueueMessage, orderDetails models.OrderMetadata) bool {
	var nftId string
	if order.VendorName == constants.NFTPORTVENDOR {
		tokenId, pending, err := svc.GetMintStatusNftPort(ctx, order)
		if err != nil {
			logger.Error("Process-OGMINT:Error Calling MINT API in NftPort", map[string]interface{}{"orderId": order.OrderId, "error": err.Error(), "order": order})
			return svc.CreateQueueEntryFromExistingQueue(ctx, order, order.Status, true, order.RetryCount+1, order.TxIdentifier, order.Type, err.Error())
		}
		if pending {
			logger.Info("Process-OGMINT:Order In Pending State", map[string]interface{}{"orderId": order.OrderId, "order": order})
			return false
		}
		nftId = tokenId
	} else {
		tokenId, txHash, pending, err := svc.GetMintStatusCrossMint(ctx, order)
		if err != nil {
			logger.Error("Process-OGMINT:Error Calling MINT API in CrossMint", map[string]interface{}{"orderId": order.OrderId, "error": err.Error(), "order": order})
			return svc.CreateQueueEntryFromExistingQueue(ctx, order, order.Status, true, order.RetryCount+1, order.TxIdentifier, order.Type, err.Error())
		}
		if pending {
			logger.Info("Process-OGMINT:Order In Pending State", map[string]interface{}{"orderId": order.OrderId, "order": order})
			return false
		}
		nftId = tokenId
		svc.createNewTransaction(order.OrderId, txHash, order.OrderType)

	}

	orderDetails.NftId = nftId
	orderDetails.Status = constants.SUCCESS
	orderDetails.UpdatedAt = time.Now()
	orderDetails.RetryCount = order.RetryCount
	if err := svc.repos.OrderMetaDataRepo.UpdateOrderByOrderId(&orderDetails).Error; err != nil {
		// Need ALERTING
		logger.Error("Process-OGMINT::Could not update ordermetadata Table", map[string]interface{}{"orderId": order.OrderId, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "orderDetails": orderDetails})
	}
	txDetail := svc.repos.TransactionDataRepo.UpdateTxStatus(order.OrderId, constants.SUCCESS)
	if txDetail.Error != nil {
		// Need ALERTING
		logger.Error("Process-OGMINT::Could not update Tx Details Table", map[string]interface{}{"orderId": order.OrderId, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "orderDetails": orderDetails})
	}
	amplitudeEvent := types.AmplitudeMintEventInterface{
		AppName:    constants.AMPLITUDE_APP_NAME,
		Network:    order.NetworkId,
		TokenId:    "",
		Type:       constants.AMPLITUDE_MINT_EVENT_TYPE,
		Status:     constants.AMPLITUDE_STATUS_MAPPING[orderDetails.Status],
		Product:    constants.AMPLITUDE_PRODUCT_TYPE,
		TokenCount: 1,
	}
	err := svc.NotifyAmplitudeEvent(order.OrderId, order.UserId, amplitudeEvent)
	if err != nil {
		logger.Error("Error Notifying Amplitude Service", map[string]interface{}{"errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "user": order.UserId, "orderId": order.OrderId})
	}
	return true
}

func (svc *OGMintSvc) CreateQueueEntryFromExistingQueue(ctx context.Context, order types.OGMintQueueMessage, status string, isRetry bool, retry_count int, txIdentifier string, processType constants.OGOPTYPE, errorMessage string) bool {

	order.Status = status
	order.RetryCount = retry_count
	order.TxIdentifier = txIdentifier
	order.Type = processType
	order.ErrorMessage = errorMessage
	delaySeconds := svc.Queues[constants.COLLECTION_VENDOR_MAPPING[order.VendorName]+order.NetworkId].QueueProps.SQSWaitTime
	if isRetry {
		delaySeconds = (retry_count + 1) * delaySeconds
	}

	if err := svc.Queues[constants.COLLECTION_VENDOR_MAPPING[order.VendorName]+order.NetworkId].WrapperObj.PublishMessage(ctx, order, order.OrderId, delaySeconds); err != nil {
		logger.Error("Process-OGMINT:Error publishing message to queue", map[string]interface{}{"orderId": order.OrderId, "error": err.Error()})
		return false
	}
	return true

}

func (svc *OGMintSvc) MintNFTNftPort(order types.OGMintQueueMessage) (string, error) {
	mintResponse, err := svc.services.NFTPORTSvc.MintNft(nil, constants.NETWORK_ID_TO_NFTPORT_CHAIN_MAPPING[order.NetworkId], order.ContractIdentifier, order.MetaDataURI, order.ToAddress)
	if err != nil {
		return "", err

	}
	if len(mintResponse.TransactionHash) <= 0 {
		return "", fmt.Errorf("No TxHash in response")
	}
	return mintResponse.TransactionHash, nil

}

func (svc *OGMintSvc) MintNFTCrossMint(order types.OGMintQueueMessage) (string, error) {
	var isBYOC bool
	if order.OrderType == constants.SS_MINT {
		isBYOC = true
	}
	mintResponse, err := svc.services.CrossMintSvc.MintNft(context.TODO(), constants.NETWORK_ID_TO_CROSSMINT_CHAIN_MAPPING[order.NetworkId], order.ContractIdentifier, order.MetaDataURI, order.ToAddress, isBYOC)
	if err != nil {
		return "", err

	}
	if len(mintResponse.ID) == 0 {
		return "", fmt.Errorf("No Id in response")
	}
	if mintResponse.OnChain.Status != "pending" {
		return "", fmt.Errorf("Unidentified status in response")
	}
	return mintResponse.ID, nil

}

func (svc *OGMintSvc) GetMintStatusNftPort(ctx context.Context, order types.OGMintQueueMessage) (string, bool, error) {
	statusResponse, err := svc.services.NFTPORTSvc.FetchMintStatus(ctx, constants.NETWORK_ID_TO_NFTPORT_CHAIN_MAPPING[order.NetworkId], order.TxIdentifier)
	if err != nil {
		return "", false, err
	}
	if statusResponse.Error != nil && statusResponse.Error.Code == "transaction_pending" {
		return "", true, nil
	}
	if statusResponse.Error != nil {
		return "", false, fmt.Errorf(statusResponse.Error.Message)
	}
	if len(statusResponse.TokenID) <= 0 {
		logger.Error("Process-OGMINT:Unknown Error", map[string]interface{}{"order": order.OrderId, "response": statusResponse})
		return "", false, fmt.Errorf("Unknown Error, please check response")
	}
	return statusResponse.TokenID, false, nil
}

func (svc *OGMintSvc) GetMintStatusCrossMint(ctx context.Context, order types.OGMintQueueMessage) (string, string, bool, error) {
	statusResponse, err := svc.services.CrossMintSvc.FetchMintStatus(ctx, constants.NETWORK_ID_TO_NFTPORT_CHAIN_MAPPING[order.NetworkId], order.TxIdentifier, order.ContractIdentifier)
	if err != nil {
		return "", "", false, err
	}
	if len(statusResponse.ID) == 0 {
		return "", "", false, fmt.Errorf("No ID found in response")
	}
	if statusResponse.OnChain.Status == "pending" {
		return "", "", true, nil
	}
	if len(statusResponse.OnChain.TxID) == 0 || len(statusResponse.OnChain.TokenID) == 0 {
		logger.Error("Process-CrossMINT:Unknown Error", map[string]interface{}{"order": order.OrderId, "response": statusResponse})
		return "", "", false, fmt.Errorf("Unknown Error, please check response")
	}
	return statusResponse.OnChain.TokenID, statusResponse.OnChain.TxID, false, nil
}

func (svc *OGMintSvc) GetOrderDetails(orderId string) (models.OrderMetadata, error) {
	response := svc.repos.OrderMetaDataRepo.GetOrdersById(orderId)
	if response.Error != nil {
		return models.OrderMetadata{}, response.Error
	}
	result, ok := response.Result.(models.OrderMetadata)
	if !ok {
		logger.Error("Process-OGMINT: unexpected result type: ", map[string]interface{}{"errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "DBResponse": response.Result})
		return models.OrderMetadata{}, fmt.Errorf("unexpected result type: %T", response.Result)
	}
	return result, nil
}

func (svc *OGMintSvc) GetLockKeyAndTTL(networkId string, reqType constants.OGOPTYPE, vendorName string) (string, int) {
	if reqType == constants.MINT {
		return svc.MintConfigs[constants.COLLECTION_VENDOR_MAPPING[vendorName]+networkId].MintApiCacheLockKey, svc.MintConfigs[constants.COLLECTION_VENDOR_MAPPING[vendorName]+networkId].MintApiCacheLockTTL
	} else {
		return svc.MintConfigs[constants.COLLECTION_VENDOR_MAPPING[vendorName]+networkId].StatusApiCacheLockKey, svc.MintConfigs[constants.COLLECTION_VENDOR_MAPPING[vendorName]+networkId].StatusApiCacheLockTTL
	}
}

func (svc *OGMintSvc) CheckLockAllowed(key string, ttl int) (bool, error) {
	lockDuration := time.Duration(ttl) * time.Millisecond
	// NOTE: Below is a temporary code to avoid breaking/queue clogs when deploying changes from seconds to ms. SHOULD be removed later
	if ttl < 10 {
		lockDuration = time.Duration(ttl) * time.Second
	}
	lockAquired, err := svc.CacheW.Driver.Mutex(key, constants.NFT_MINT_PURPOSE, "1", lockDuration)
	if err != nil {
		return false, err
	}
	if !lockAquired {
		return false, fmt.Errorf("Could not acquire lock")
	}
	return true, nil
}

func (svc *OGMintSvc) createNewTransaction(orderId, txHash, orderType string) {
	txData := models.TransactionData{
		OrderId:        orderId,
		TxHash:         txHash,
		Status:         constants.RUNNING,
		OrderTxType:    orderType,
		PayloadType:    "MINT",
		GasUsed:        "NA",
		GasPrice:       "NA",
		TokenTransfers: nil,
	}
	if err := svc.repos.TransactionDataRepo.InsertTransactionData(&txData).Error; err != nil {
		logger.Error("Process-OGMINT:Could not insert into transactionData Table", map[string]interface{}{"orderId": orderId, "txHash": txHash, "transaction": txData})
		//TODO: Need ALERTING
	}
}
