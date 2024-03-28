package repository

import (
	"pismo-dev/internal/repository/ordermetadata"
	trasactiondata "pismo-dev/internal/repository/transactiondata"
	"pismo-dev/pkg/logger"
	"pismo-dev/pkg/storage"
)

type Repositories struct {
	OrderMetadataRepo   ordermetadata.OrderMetadataRepositoryInterface
	TransactionDataRepo trasactiondata.TransactionDataRepositoryInterface
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

func WithOrderMetadataRepo(store *storage.Store) RepositoriesOption {
	return func(repositories *Repositories) error {
		repositories.OrderMetadataRepo = ordermetadata.NewOrderMetadataRepository(store.GormPsql)
		return nil
	}
}

func WithTransactionDataRepo(store *storage.Store) RepositoriesOption {
	return func(repositories *Repositories) error {
		repositories.TransactionDataRepo = trasactiondata.NewTransactionDataRepository(store.GormPsql)
		return nil
	}
}
