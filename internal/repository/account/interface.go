package account

import (
	"pismo-dev/api/types"
	"pismo-dev/internal/models"
)

type AccountRepositoryInterface interface {
	InsertAccount(req types.CreateAccountRequest) models.Account
	GetAccountData(accountId string) models.Account
}
