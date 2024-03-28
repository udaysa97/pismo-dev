package transactiondataservice

import (
	"context"
	"pismo-dev/error/service"
	"pismo-dev/internal/models"
)

type TransactionServiceInterface interface {
	GetTransactionDetails(ctx context.Context, filters map[string][]string, pagination *models.Pagination) (models.TransactionDetails, *service.ServiceError)
}
