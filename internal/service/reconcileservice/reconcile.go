package reconcileservice

import (
	"context"
	"encoding/json"
	"fmt"
	"pismo-dev/commonpkg/utils"
	"pismo-dev/constants"
	"pismo-dev/internal/models"
	"pismo-dev/internal/types"
	"pismo-dev/pkg/logger"
	"pismo-dev/pkg/queue"
	"strconv"
	"time"
)

type ReconcileSvc struct {
	ServiceName          string
	services             RequiredServices
	repos                RequiredRepos
	SqsQueueWrapper      *queue.QueueWrapper
	SQSWaitTime          int
	SQSVisibilityTimeout int
	MaxSqsMessage        int
}

func NewReconcileSvc(sqsqw queue.QueueWrapper, sqsWaitTime int, sqsVisibilityTimeout int, maxSqsMessage int) *ReconcileSvc {
	return &ReconcileSvc{
		ServiceName:          "ReconcileService",
		SqsQueueWrapper:      &sqsqw,
		SQSWaitTime:          sqsWaitTime,
		SQSVisibilityTimeout: sqsVisibilityTimeout,
		MaxSqsMessage:        maxSqsMessage,
	}
}

func (svc *ReconcileSvc) SetRequiredServices(services RequiredServices) {
	svc.services = services
}

func (svc *ReconcileSvc) SetRequiredRepos(repos RequiredRepos) {
	svc.repos = repos
}

func (svc *ReconcileSvc) GetJobStatusFromFM(ctx context.Context, jobId string) (types.FMGetStatusResult, error) {
	getStatusResult, err := svc.services.FlowManagerSvc.GetStatus(ctx, jobId)
	if err != nil {
		logger.Error("Reconcile: failed to fetch GetStatus response from FM", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "error": err})
		return types.FMGetStatusResult{}, err
	}

	// Inserting the Transaction data into Transaction table
	if len(getStatusResult.Metadata.Transactions) > 0 {
		transactionRequestDataArray := getStatusResult.Metadata.Transactions
		for _, transaction := range transactionRequestDataArray {
			transactionDataExists := svc.repos.TransactionDataRepo.FindTransactionDataByOrderIdAndTxHash(jobId, transaction.TransactionHash)

			// Initialize TokenTransfers to an empty slice if it's nil (in case of Non-EVM transactions)
			tokenTransfers := transaction.TokenTransfers
			if tokenTransfers == nil {
				tokenTransfers = []map[string]interface{}{}
			}

			if !transactionDataExists {

				response := svc.repos.TransactionDataRepo.InsertTransactionData(&models.TransactionData{OrderId: transaction.JobId, TxHash: transaction.TransactionHash, Status: transaction.TransactionState, OrderTxType: transaction.OrderTxType, PayloadType: transaction.PayloadType, GasUsed: strconv.Itoa(transaction.GasUsed), GasPrice: utils.ToString(transaction.GasPrice), TokenTransfers: tokenTransfers})
				if response.Error != nil {
					logger.Error("Reconcile: error in inserting transaction data: ", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "Error": err})
				}
			}
		}
	}
	return getStatusResult, nil
}

func (svc *ReconcileSvc) GetJobStatus(ctx context.Context, jobId string) (string, error) {
	response := svc.repos.OrderMetaDataRepo.GetOrdersById(jobId)
	if response.Error != nil {
		return "", response.Error
	}
	result, ok := response.Result.(models.OrderMetadata)
	if !ok {
		logger.Error("Reconcile: unexpected result type: ", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "DBResponse": response.Result})
		return "", fmt.Errorf("unexpected result type: %T", response.Result)
	}
	status := result.Status
	return status, nil
}

func (svc *ReconcileSvc) UpdateJobStatus(ctx context.Context, jobId string, status types.FMGetStatusResult) (bool, error) {
	dbResult := svc.repos.OrderMetaDataRepo.GetOrdersById(jobId)
	if dbResult.Error != nil {
		logger.Info("Reconcile: Could not retreive data from DB", map[string]interface{}{"statusEvent": status, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "Error": dbResult.Error.Error()})
		return false, dbResult.Error
	}
	savedNftTransfer, ok := dbResult.Result.(models.OrderMetadata)
	if !ok {
		logger.Error("Reconcile: unexpected result type: ", map[string]interface{}{"errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "DBResponse": dbResult.Result})
		return false, fmt.Errorf("unexpected result type: %T", dbResult.Result)
	}

	if status.Status == constants.SUCCESS {
		txDbObj, _, txStarted := svc.repos.OrderMetaDataRepo.GetDbTransaction()
		if !txStarted {
			logger.Error("Reconcile: Could not obtain transaction", map[string]interface{}{"statusEvent": status, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode})
		}
		savedNftTransfer.Status = constants.SUCCESS
		savedNftTransfer.UpdatedAt = time.Now()

		if status.Metadata.Transactions != nil && len(status.Metadata.Transactions) > 0 {
			slice := status.Metadata.Transactions
			for _, transaction := range slice {
				// Processing NFT Mint transactions for APTOS to save newly minted NFT's id
				if transaction.NetworkType == constants.NETWORK_TYPE_APTOS && transaction.Receipt != nil {
					// this is done to cast transaction.Receipt to types.IAptosTransactionReceipt
					marshalReceipt, err := json.Marshal(transaction.Receipt)
					if err != nil {
						logger.Error("Reconcile: Unable to marshal APTOS transaction receipt", map[string]interface{}{"receipt": transaction.Receipt, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "Error": err})
						return false, err
					}
					var receipt types.IAptosTransactionReceipt
					er := json.Unmarshal(marshalReceipt, &receipt)
					if er != nil {
						logger.Error("Reconcile: Unable to unmarshal APTOS transaction receipt", map[string]interface{}{"marshalReceipt": marshalReceipt, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "Error": er})
						return false, er
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
			}
		}

		//TODO: Remove this after testing
		logger.Error("Reconcile: Updated NFT Order Details", map[string]interface{}{"order": savedNftTransfer})

		if err := txDbObj.UpdateOrderByOrderId(&savedNftTransfer).Error; err != nil {
			logger.Error("Reconcile: Could not update ordermetadata Table", map[string]interface{}{"statusEvent": status, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "dbResult": savedNftTransfer})
			txDbObj.RollBack()
			return false, err
		}

		if err := txDbObj.Commit(); err != nil {
			logger.Info("Reconcile: Unable to commit transaction", map[string]interface{}{"txObj": savedNftTransfer, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "error": err, "jobId": jobId})
			txDbObj.RollBack()
			return false, err
		}
	} else {
		// This is additional log to monitor nonce issue coming at signign service for stan mints
		if len(savedNftTransfer.VendorId) > 0 && savedNftTransfer.OrderType == constants.NFT_MINT {
			logger.Error("NFT_VENDOR_MINT_ERROR: Please check reason for failure", map[string]interface{}{"orderid": jobId, "userid": savedNftTransfer.UserId, "vendorid": savedNftTransfer.VendorId})
		}
		logger.Error("NFT Order Failed or Rejected", map[string]interface{}{"jobId": jobId, "status": status})
		errorInResponse, err := json.Marshal(status.Error)
		failureReason := "something went wrong"
		if err != nil {
			logger.Error("Reconcile: failed to marshal Error response for failed order", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "error": err.Error(), "statusError": status.Error})
		} else {
			failureReason = fmt.Sprintf("%s:%s", failureReason, string(errorInResponse))
		}

		if status.Output != nil && status.Output.Error != nil && status.Output.Error.ErrorCode != "" {
			failureReason = constants.FE_ERROR_CODES_MAPPING[status.Output.Error.ErrorCode].Message
			logger.Debug(fmt.Sprintf("JobStatusTracking: failureReason %s", failureReason))
		} else {
			logger.Info("Reconcile: Could not detect error in failure", map[string]interface{}{"respone": status})
		}
		response := svc.repos.OrderMetaDataRepo.UpdateOrderStatusForFailedOrder(jobId, status.Status, failureReason)
		if response.Error != nil {
			logger.Error("Reconcile: failed to update orderMetadata status for failed order", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "Response": response.Error.Error()})
			return false, response.Error
		}
	}

	return true, nil
}

func (svc *ReconcileSvc) InitiateSqsConsumer() {
	svc.SqsQueueWrapper.InitiateQueueConsumption(svc.SQSVisibilityTimeout, svc.SQSWaitTime, svc.MaxSqsMessage, svc.CheckJobStatus, false)
}

func (svc *ReconcileSvc) CheckJobStatus(message string) bool {
	var sqsMessage types.SQSJobInterface
	err := json.Unmarshal([]byte(message), &sqsMessage)
	if err != nil {
		logger.Error("Reconcile: Could not Read SQS message", map[string]interface{}{"message": message, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode})
		return false
	}
	jobId := sqsMessage.JobId
	if len(jobId) <= 0 {
		logger.Error("Reconcile: Could not Read SQS message: No jobId present", map[string]interface{}{"message": message, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode})
		return true
	}
	logger.Info("Reconcile: Started reconcile for jobId", map[string]interface{}{"jobId": jobId})
	return svc.ProcessJobStatus(context.TODO(), jobId)
}

func (svc *ReconcileSvc) ProcessJobStatus(ctx context.Context, jobId string) bool {
	statusResponse, err := svc.GetJobStatusFromFM(ctx, jobId)
	if err != nil {
		logger.Error("Reconcile: Error in fetching JobStatus from FM", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "error": err})
		return false
	}

	dbStatus, err := svc.GetJobStatus(ctx, jobId)
	if err != nil {
		logger.Error("Reconcile: error in fetching Status from OrderMetadata Repo", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "error": err})
		return false
	}
	//terminal state check
	if dbStatus == constants.SUCCESS || dbStatus == constants.FAILED || dbStatus == constants.REJECTED {
		if statusResponse.Status != constants.SUCCESS && statusResponse.Status != constants.FAILED && statusResponse.Status != constants.REJECTED {
			logger.Error("Reconcile: jobStatus for terminal state not updated", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "JobId": jobId})
		}
		return true
	}

	if statusResponse.Status == constants.SUCCESS || statusResponse.Status == constants.FAILED || statusResponse.Status == constants.REJECTED {
		_, err := svc.UpdateJobStatus(ctx, jobId, statusResponse)
		return err == nil
	}
	return false
}
