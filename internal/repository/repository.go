package repository

import (
	"pismo-dev/internal/repository/account"
	trasactiondata "pismo-dev/internal/repository/transactiondata"
	"pismo-dev/pkg/logger"
	"pismo-dev/pkg/storage"
)

type Repositories struct {
	AccountRepo     account.AccountRepositoryInterface
	TransactionRepo trasactiondata.TransactionDataRepositoryInterface
}

type RepositoriesOption func(*Repositories) error

func New(opts ...RepositoriesOption) *Repositories {
	repositories := &Repositories{}

	// Loop through each option
	for _, opt := range opts {
		// Call the option giving the instantiated
		// *service as the argument
		err := opt(repositories)
		if err != nil {
			logger.Error("Error: ", err)
		}
	}

	return repositories
}

func WithAccountRepo(store *storage.Store) RepositoriesOption {
	return func(repositories *Repositories) error {
		repositories.AccountRepo = account.NewAccountRepository(store.GormPsql)
		return nil
	}
}

func WithTransactionDataRepo(store *storage.Store) RepositoriesOption {
	return func(repositories *Repositories) error {
		repositories.TransactionRepo = trasactiondata.NewTransactionRepository(store.GormPsql)
		return nil
	}
}
