package transactiondata

import (
	"pismo-dev/internal/models"
	"pismo-dev/internal/repository/result"

	"gorm.io/gorm"
)

type TransactionDataRepositoryInterface interface {
	InsertTransactionData(transactionData *models.TransactionData) result.RepositoryResult
	GetAllTransactions(filter map[string][]string, pagination *models.Pagination) result.RepositoryResult
	UpdateTransactionDataByOrderIdAndTxHash(orderId string, txHash string, transactionData *models.TransactionData) result.RepositoryResult
	FindTransactionDataByOrderIdAndTxHash(orderId string, txHash string) bool
	SetDBTransaction(transactionObj *gorm.DB) TransactionDataRepository
	UpdateTxStatus(orderId string, status string) result.RepositoryResult
}
