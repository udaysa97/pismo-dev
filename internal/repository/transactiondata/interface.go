package transactiondata

import (
	"pismo-dev/api/types"
	"pismo-dev/internal/models"
)

type TransactionDataRepositoryInterface interface {
	InsertTransaction(req types.TransactionRequest) models.Transaction
}
