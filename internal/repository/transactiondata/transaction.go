package transactiondata

import (
	"pismo-dev/api/types"
	"pismo-dev/internal/models"
	"pismo-dev/internal/repository/result"

	"gorm.io/gorm"
)

type TransactionsRepository struct {
	db        *gorm.DB
	TableName string
}

func NewTransactionRepository(db *gorm.DB) *TransactionsRepository {
	return &TransactionsRepository{db: db, TableName: "transactions"}
}

func (er *TransactionsRepository) GetAllTransactions(filter map[string][]string, pagination *models.Pagination) result.RepositoryResult {
	var transactions []models.Transaction
	var query = &models.Transaction{}
	offset := (pagination.Page - 1) * pagination.Limit
	queryBuilder := er.db.Limit(pagination.Limit).Offset(offset).Order(pagination.Sort + " " + pagination.Direction)
	if filter != nil {

	}
	var count int64
	err := queryBuilder.Where(query).Find(&transactions).Count(&count).Error
	if err != nil {
		return result.RepositoryResult{Error: err}
	}
	return result.RepositoryResult{Result: transactions}
}

func (er *TransactionsRepository) InsertTransaction(req types.TransactionRequest) models.Transaction {

	td := models.Transaction{AccountID: uint(req.AccountId), OperationTypeID: uint(req.OperationId), Amount: float64(*req.Amount)}
	err := er.db.Create(&td).Error
	if err != nil {
		return models.Transaction{}
	}
	return td

}
