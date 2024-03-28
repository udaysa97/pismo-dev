package transactiondataservice

import (
	"context"
	"pismo-dev/constants"
	"pismo-dev/error/service"
	"pismo-dev/internal/models"
	"pismo-dev/internal/repository/transactiondata"
)

type TransactionDataService struct {
	transactiondata transactiondata.TransactionDataRepositoryInterface
}

func NewTransactionDataService(transactionInterface transactiondata.TransactionDataRepositoryInterface) *TransactionDataService {
	return &TransactionDataService{transactiondata: transactionInterface}

}

func (transactionDataSvc *TransactionDataService) GetTransactionDetails(ctx context.Context, filters map[string][]string, pagination *models.Pagination) (models.TransactionDetails, *service.ServiceError) {
	var result = transactionDataSvc.transactiondata.GetAllTransactions(filters, pagination)
	if result.Error != nil {
		return models.TransactionDetails{}, service.New(result.Error.Error(), constants.ERROR_TYPES[constants.BAD_REQUEST_ERROR].ErrorCode)
	}
	metadata := result.Result.(models.TransactionDetails)
	return metadata, nil
}
