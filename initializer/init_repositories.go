package initializer

import (
	"pismo-dev/internal/repository"
	"pismo-dev/pkg/storage"
)

func InitRepositories(store *storage.Store) *repository.Repositories {
	repositoryOptions := []repository.RepositoriesOption{
		repository.WithTransactionDataRepo(store),
		repository.WithAccountRepo(store),
	}
	return repository.New(repositoryOptions...)
}
