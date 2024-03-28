package transactiondata

import (
	"encoding/json"
	"pismo-dev/internal/models"
	"pismo-dev/internal/query"
	"pismo-dev/internal/repository/result"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type TransactionDataRepository struct {
	db        *gorm.DB
	TableName string
}

func NewTransactionDataRepository(db *gorm.DB) *TransactionDataRepository {
	return &TransactionDataRepository{db: db, TableName: "transaction_data"}
}

func (er *TransactionDataRepository) GetAllTransactions(filter map[string][]string, pagination *models.Pagination) result.RepositoryResult {
	var orders []models.TransactionData
	var query = &models.TransactionData{}
	offset := (pagination.Page - 1) * pagination.Limit
	queryBuilder := er.db.Limit(pagination.Limit).Offset(offset).Order(pagination.Sort + " " + pagination.Direction)
	if filter != nil {
		if filter["order_id"] != nil && len(filter["order_id"]) > 0 {
			query.OrderId = filter["order_id"][0]
		}

		if filter["status"] != nil && len(filter["status"]) > 0 {
			query.Status = filter["status"][0]
		}

		if filter["order_type"] != nil && len(filter["order_type"]) > 0 {
			query.OrderTxType = filter["order_type"][0]
		}
		if filter["created_at_from"] != nil && len(filter["created_at_from"]) > 0 {
			unixTimestamp, err := strconv.ParseInt(filter["created_at_from"][0], 10, 64)
			var createdAtFrom time.Time
			if err == nil {
				createdAtFrom = time.Unix(unixTimestamp, 0).UTC()
			}
			if !createdAtFrom.IsZero() {
				queryBuilder = queryBuilder.Where("created_at >= ?", createdAtFrom)
			}
		}
		if filter["created_at_to"] != nil && len(filter["created_at_to"]) > 0 {
			unixTimestamp, err := strconv.ParseInt(filter["created_at_to"][0], 10, 64)
			var createdAtTo time.Time
			if err == nil {
				createdAtTo = time.Unix(unixTimestamp, 0).UTC()
			}
			if !createdAtTo.IsZero() {
				queryBuilder = queryBuilder.Where("created_at <= ?", createdAtTo)
			}
		}
	}
	var count int64
	err := queryBuilder.Where(query).Omit("token_transfers").Find(&orders).Count(&count).Error
	if err != nil {
		return result.RepositoryResult{Error: err}
	}
	return result.RepositoryResult{Result: models.TransactionDetails{Count: count, TransactionDetails: orders}}
}

func (er *TransactionDataRepository) InsertTransactionData(transactionData *models.TransactionData) result.RepositoryResult {
	var err error
	var td query.TransactionData
	if transactionData.TokenTransfers != nil {
		tokenTransferMarshalObj, _ := json.Marshal(transactionData.TokenTransfers)
		tokenTransfers := string(tokenTransferMarshalObj)
		td = query.TransactionData{OrderId: transactionData.OrderId, TxHash: transactionData.TxHash, Status: transactionData.Status, OrderTxType: transactionData.OrderTxType, PayloadType: transactionData.PayloadType, GasUsed: transactionData.GasUsed, GasPrice: transactionData.GasPrice, TokenTransfers: tokenTransfers}
		err = er.db.Create(&td).Error
	} else {
		err = er.db.Create(&transactionData).Error
	}
	if err != nil {
		return result.RepositoryResult{Error: err}
	}
	return result.RepositoryResult{Result: transactionData}
}

func (er *TransactionDataRepository) UpdateTransactionDataByOrderIdAndTxHash(orderId string, txHash string, transactionData *models.TransactionData) result.RepositoryResult {
	err := er.db.Omit("created_at").Where(models.TransactionData{OrderId: orderId, TxHash: txHash}).Save(&transactionData).Error
	if err != nil {
		return result.RepositoryResult{Error: err}
	}
	return result.RepositoryResult{Result: transactionData}
}

func (er *TransactionDataRepository) FindTransactionDataByOrderIdAndTxHash(orderId string, txHash string) bool {
	var transactionData models.TransactionData
	err := er.db.Where(&models.TransactionData{OrderId: orderId, TxHash: txHash}).First(&transactionData).Error
	return err == nil
}

func (er *TransactionDataRepository) UpdateTxStatus(orderId string, status string) result.RepositoryResult {
	err := er.db.Model(&models.TransactionData{}).Omit("created_at").Where(models.TransactionData{OrderId: orderId}).Update("status", status).Error
	if err != nil {
		return result.RepositoryResult{Error: err}
	}
	return result.RepositoryResult{Result: true}
}

func (er *TransactionDataRepository) SetDBTransaction(transactionObj *gorm.DB) TransactionDataRepository {
	return TransactionDataRepository{db: transactionObj, TableName: "Temp_transaction_data"}
}
