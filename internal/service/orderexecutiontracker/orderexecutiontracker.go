package orderexecutiontracker

import (
	"encoding/json"
	"fmt"
	"pismo-dev/commonpkg/utils"
	"pismo-dev/constants"
	"pismo-dev/internal/appconfig"
	"pismo-dev/internal/models"
	"pismo-dev/internal/types"
	kafkaclient "pismo-dev/pkg/kafka/client"
	"pismo-dev/pkg/logger"
	"strconv"
	"time"
)

type OrderExecutionTrackerSvc struct {
	ServiceName         string
	repos               RequiredRepos
	KafkaProducerClient *kafkaclient.ProducerClient
	KafkaConsumerClient *kafkaclient.ConsumerClient
	FMConsumerTopic     string
	AmplitudeEventTopic string
}

func NewOrderExecutionTrackerSvc(kafkaProducerClient *kafkaclient.ProducerClient, kafkaConsumerClient *kafkaclient.ConsumerClient, amplitudeEventTopic string, fmConsumerTopic string) *OrderExecutionTrackerSvc {
	return &OrderExecutionTrackerSvc{
		ServiceName:         "NFTOrderExecutionTrackerSvc",
		KafkaProducerClient: kafkaProducerClient,
		KafkaConsumerClient: kafkaConsumerClient,
		AmplitudeEventTopic: amplitudeEventTopic,
		FMConsumerTopic:     fmConsumerTopic,
	}
}

func (svc *OrderExecutionTrackerSvc) SetRequiredRepos(repos RequiredRepos) {
	svc.repos = repos
}

func (svc *OrderExecutionTrackerSvc) InitJobTrackerConsumer() {
	logger.Info("Consumer Initiated!")
	err := svc.KafkaConsumerClient.Consume(svc.FMConsumerTopic, svc.ProcessJobEvent)
	if err != nil {
		logger.Error("Urgent! Kafka Consumer Stopped!", map[string]interface{}{"Error": err, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "topic": svc.FMConsumerTopic})
		time.Sleep(5 * time.Second) // Giving breathing time before retrying connection
		svc.KafkaConsumerClient, _ = kafkaclient.NewConsumerClient(appconfig.KAFKA_HOST, appconfig.KAFKA_GROUP_ID, appconfig.FM_CONSUMER_POLL_INTERVAL)
		svc.InitJobTrackerConsumer()
	}
}

func (svc *OrderExecutionTrackerSvc) ProcessJobEvent(message []byte) bool {
	var event types.JobExecutionKafka
	//logger.Info("JobStatusTracking: processJobEvent started")
	err := json.Unmarshal(message, &event)
	if err != nil {
		logger.Error("JobStatusTracking: Unable to unmarshal message", map[string]interface{}{"kafkamessage": message, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "Error": err})
		return true
	}

	//TODO: Check behaviour when no matching entry in DB. As that is a valid case and needs to be skipped
	dbResult := svc.repos.OrderMetaDataRepo.GetOrdersById(event.JobId)
	if dbResult.Error != nil {
		logger.Info("JobStatusTracking: Could not retreive data from DB", map[string]interface{}{"kafkamessage": message, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "Error": dbResult.Error.Error()})
		return true
	}
	logger.Info("JobStatusTracking: event received", map[string]interface{}{"kafkamessage": event, "jobId": event.JobId})

	savedNftTransfer, ok := dbResult.Result.(models.OrderMetadata)
	if !ok {
		logger.Error("JobStatusTracking: unexpected result type: ", map[string]interface{}{"errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "DBResponse": dbResult.Result})
		return false
	}
	if savedNftTransfer.Status == constants.SUCCESS || savedNftTransfer.Status == constants.FAILED || savedNftTransfer.Status == constants.REJECTED {
		logger.Info("JobStatusTracking: Terminal state found", map[string]interface{}{"jobId": event.JobId, "status": savedNftTransfer.Status})
		return true
	}
	if event.Status == constants.SUCCESS {
		txDbObj, dbObj, txStarted := svc.repos.OrderMetaDataRepo.GetDbTransaction()
		if !txStarted {
			logger.Error("JobStatusTracking: Could not obtain transaction", map[string]interface{}{"kafkamessage": event, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "dbResult": savedNftTransfer})
		}
		savedNftTransfer.Status = constants.SUCCESS
		savedNftTransfer.UpdatedAt = time.Now()

		if event.Metadata != nil && event.Metadata.Transactions != nil && len(*event.Metadata.Transactions) > 0 {
			txDbObjTransaction := svc.repos.TransactionDataRepo.SetDBTransaction(dbObj) // Sets transaction(atomic) for transaction repo using atomic obj created in ordermetadata
			slice := *event.Metadata.Transactions
			for _, transaction := range slice {
				// Processing NFT Mint transactions for APTOS to save newly minted NFT's id
				if transaction.NetworkType == constants.NETWORK_TYPE_APTOS && transaction.Receipt != nil {
					// this is done to cast transaction.Receipt to types.IAptosTransactionReceipt
					marshalReceipt, err := json.Marshal(transaction.Receipt)
					if err != nil {
						logger.Error("JobStatusTracking: Unable to marshal APTOS transaction receipt", map[string]interface{}{"receipt": transaction.Receipt, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "Error": err})
						return false
					}
					var receipt types.IAptosTransactionReceipt
					er := json.Unmarshal(marshalReceipt, &receipt)
					if er != nil {
						logger.Error("JobStatusTracking: Unable to unmarshal APTOS transaction receipt", map[string]interface{}{"marshalReceipt": marshalReceipt, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "Error": er})
						return false
					}
					events := receipt.Events
					if len(events) > 0 {
						for _, event := range events {
							if event.Type == constants.NFT_MINT_APTOS_HASH_TYPE {
								savedNftTransfer.NftId = event.Data["token"].(string)
								break
							}
						}
					}
				}

				// Initialize TokenTransfers to an empty slice if it's nil (in case of Non-EVM transactions)
				tokenTransfers := transaction.TokenTransfers
				if tokenTransfers == nil {
					tokenTransfers = []map[string]interface{}{}
				}

				txData := models.TransactionData{
					OrderId:        event.JobId,
					TxHash:         transaction.TransactionHash,
					Status:         strconv.FormatBool(transaction.TransactionStatus),
					OrderTxType:    transaction.OrderTxType,
					PayloadType:    transaction.PayloadType,
					GasUsed:        strconv.Itoa(transaction.GasUsed),
					GasPrice:       utils.ToString(transaction.GasPrice),
					TokenTransfers: tokenTransfers,
				}

				if err = txDbObjTransaction.InsertTransactionData(&txData).Error; err != nil {
					logger.Error("JobStatusTracking:Could not insert into transactionData Table", map[string]interface{}{"kafkamessage": event, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "dbResult": savedNftTransfer, "transaction": txData})
					txDbObj.RollBack()
					return false
				}
			}
		}

		//TODO: Remove this after testing
		logger.Error("JobStatusTracking: Updated NFT Order Details", map[string]interface{}{"order": savedNftTransfer})

		if err = txDbObj.UpdateOrderByOrderId(&savedNftTransfer).Error; err != nil {
			logger.Error("JobStatusTracking:Could not update ordermetadata Table", map[string]interface{}{"kafkamessage": event, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "dbResult": savedNftTransfer})
			txDbObj.RollBack()
			return false
		}

		if err := txDbObj.Commit(); err != nil {
			logger.Info("JobStatusTracking:Unable to commit transaction", map[string]interface{}{"txObj": savedNftTransfer, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "error": err, "jobId": event.JobId})
			txDbObj.RollBack()
			return false
		}
		//If any error txDbObj.Rollback()

	} else if event.Status == constants.FAILED || event.Status == constants.REJECTED {
		// This is additional log to monitor nonce issue coming at signign service for stan mints
		if len(savedNftTransfer.VendorId) > 0 && savedNftTransfer.OrderType == constants.NFT_MINT {
			logger.Error("NFT_VENDOR_MINT_ERROR: Please check reason for failure", map[string]interface{}{"orderid": event.JobId, "userid": savedNftTransfer.UserId, "vendorid": savedNftTransfer.VendorId, "event": event})
		}
		logger.Error("NFT Order Failed or Rejected", map[string]interface{}{"event": event, "jobId": event.JobId, "status": event.Status})
		failureReason := "Something went wrong. Try again"
		txDbObj, dbObj, txStarted := svc.repos.OrderMetaDataRepo.GetDbTransaction()
		if !txStarted {
			logger.Error("JobStatusTracking: Could not obtain transaction", map[string]interface{}{"event": event, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "dbResult": savedNftTransfer})
		}
		savedNftTransfer.Status = event.Status
		if event.Error != nil {
			failureReason = constants.FE_ERROR_CODES_MAPPING[event.Error.ErrorCode].Message
			logger.Debug(fmt.Sprintf("JobStatusTracking: failureReason %s", failureReason))
			if event.Error.ErrorCode == constants.FE_ERROR_CODES_MAPPING["FE0004"].ErrorCode {
				failureReason = constants.FE_ERROR_CODES_MAPPING["FE0004"].SubErrorName[event.Error.Name]
			}
			if event.Error.Details != nil {
				details, err := json.Marshal(event.Error.Details)
				if err != nil {
					logger.Warn("Could not marshal error details", map[string]interface{}{"event": event, "error": err, "jobId": event.JobId})
				} else {
					failureReason += string(details)
				}
			}
		}
		savedNftTransfer.ExecutionResponse = failureReason
		if err = txDbObj.UpdateOrderByOrderId(&savedNftTransfer).Error; err != nil {
			logger.Error("JobStatusTracking:Could not update ordermetadata Table", map[string]interface{}{"event": event, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "dbResult": savedNftTransfer})
			txDbObj.RollBack()
			return false
		}
		if event.Metadata != nil && event.Metadata.Transactions != nil && len(*event.Metadata.Transactions) > 0 {
			txDbObjTransaction := svc.repos.TransactionDataRepo.SetDBTransaction(dbObj) // Sets transaction(atomic) for transacction repo using atomic obj created in ordermetadata
			slice := *event.Metadata.Transactions
			for _, transaction := range slice {
				// Initialize TokenTransfers to an empty slice if it's nil (in case of Non-EVM transactions)
				tokenTransfers := transaction.TokenTransfers
				if tokenTransfers == nil {
					tokenTransfers = []map[string]interface{}{}
				}

				txData := models.TransactionData{
					OrderId:        event.JobId,
					TxHash:         transaction.TransactionHash,
					Status:         strconv.FormatBool(transaction.TransactionStatus),
					OrderTxType:    transaction.OrderTxType,
					PayloadType:    transaction.PayloadType,
					GasUsed:        strconv.Itoa(transaction.GasUsed),
					GasPrice:       utils.ToString(transaction.GasPrice),
					TokenTransfers: tokenTransfers,
				}

				if err = txDbObjTransaction.InsertTransactionData(&txData).Error; err != nil {
					logger.Error("JobStatusTracking:Could not insert into transactionData Table", map[string]interface{}{"event": event, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "dbResult": savedNftTransfer, "transaction": txData})
					txDbObj.RollBack()
					return false
				}
			}
		}
		if err := txDbObj.Commit(); err != nil {
			logger.Info("JobStatusTracking:Unable to commit transaction", map[string]interface{}{"txObj": savedNftTransfer, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "error": err, "jobId": event.JobId})
			txDbObj.RollBack()
			return false
		}
		//TODO: PUSH AMPLITUDE EVENTY

	} else if event.Status == constants.RUNNING {
		savedNftTransfer.Status = constants.RUNNING
		savedNftTransfer.UpdatedAt = time.Now()
		if err = svc.repos.OrderMetaDataRepo.UpdateOrderByOrderId(&savedNftTransfer).Error; err != nil {
			logger.Error("JobStatusTracking:Could not update ordermetadata Table", map[string]interface{}{"event": event, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "dbResult": savedNftTransfer})
			return false
		}
		//TODO: PUSH AMPLITUDE EVENT
	} else if event.Status == constants.WAITING_FOR_SIGNATURE {
		savedNftTransfer.Status = constants.WAITING_FOR_SIGNATURE
		savedNftTransfer.UpdatedAt = time.Now()
		if err = svc.repos.OrderMetaDataRepo.UpdateOrderByOrderId(&savedNftTransfer).Error; err != nil {
			logger.Error("JobStatusTracking:Could not update ordermetadata Table", map[string]interface{}{"event": event, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "dbResult": savedNftTransfer})
			return false
		}
		//TODO: PUSH AMPLITUDE EVENT
	}
	err = svc.notifyAmplitudeEvent("", savedNftTransfer.NetworkId, savedNftTransfer.OrderId, savedNftTransfer.Status, "", "", savedNftTransfer.UserId)
	if err != nil {
		logger.Error("Error Notifying Amplitude Service in order tracker", map[string]interface{}{"errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "user": savedNftTransfer.UserId, "orderId": savedNftTransfer.OrderId})
	}
	return true

}

func (svc *OrderExecutionTrackerSvc) notifyAmplitudeEvent(source, network, orderId, status, nftName, ercType, userId string) error {
	amplitudeEvent := types.AmplitudeTransferEventInterface{
		AppName:    constants.AMPLITUDE_APP_NAME,
		Network:    network,
		Token:      nftName, //TODO: Retrieve name
		NftType:    "",      //TODO: Retrieve Type, image,GIF etc
		ErcType:    ercType,
		Type:       constants.AMPLITUDE_TRANSFER_EVENT_TYPE,
		OrderId:    orderId,
		Status:     constants.AMPLITUDE_STATUS_MAPPING[status],
		DeviceId:   "",
		UserId:     userId,
		DeviceType: source,
	}
	amplitudeNotification := types.AmplitudeEventInterface{
		UserId:          userId,
		Eventtype:       constants.AMPLITUDE_EVENT_NAME,
		EventProperties: amplitudeEvent,
	}
	return svc.KafkaProducerClient.Produce(svc.AmplitudeEventTopic, orderId, amplitudeNotification)
}
