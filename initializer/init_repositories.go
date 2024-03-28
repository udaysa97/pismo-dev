package initializer

import (
	"pismo-dev/internal/repository"
	"pismo-dev/pkg/storage"
)

func InitRepositories(store *storage.Store) *repository.Repositories {
	repositoryOptions := []repository.RepositoriesOption{
		repository.WithOrderMetadataRepo(store),
		repository.WithTransactionDataRepo(store),
	}
	return repository.New(repositoryOptions...)
}
