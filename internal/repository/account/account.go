package account

import (
	"pismo-dev/api/types"
	"pismo-dev/internal/models"
	"strconv"

	"gorm.io/gorm"
)

type AccountRepository struct {
	db        *gorm.DB
	TableName string
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db, TableName: "account"}
}

func (er *AccountRepository) GetAccountData(accountId string) models.Account {
	var transactions models.Account
	var query = &models.Account{}
	id, _ := strconv.ParseUint(accountId, 10, 32)
	query.ID = uint(id)
	var count int64
	err := er.db.Where(query).Find(&transactions).Count(&count).Error
	if err != nil {
		return models.Account{}
	}
	return transactions
}

func (er *AccountRepository) InsertAccount(req types.CreateAccountRequest) models.Account {

	td := models.Account{DocumentNumber: req.DocumentId}
	err := er.db.Create(&td).Error
	if err != nil {
		return models.Account{}
	}
	return td

}
