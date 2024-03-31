package initializer

import (
	"pismo-dev/internal/repository"
	"pismo-dev/internal/service"
)

func InitServices(repositories *repository.Repositories) (services *service.Service) {
	options := []service.ServiceOption{
		//service.WithTransactionDataSvc(repositories.TransactionDataRepo),
	}
	services = service.NewService(options...)

	return
}
