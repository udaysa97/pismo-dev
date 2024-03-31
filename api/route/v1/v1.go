package v1

import (
	transaction "pismo-dev/api/handler/v1/transactions"
	"pismo-dev/api/middleware/log"

	"pismo-dev/internal/repository"
	"pismo-dev/internal/service"

	"github.com/gin-gonic/gin"
)

// Auth error code will be handled separately under identifier "DEFAULT"
func Register(routerGroup *gin.RouterGroup, services *service.Service, repos *repository.Repositories) {
	routerGroup.POST("/accounts", log.LogMiddleware(), transaction.InsertTransaction(services, repos))
	routerGroup.GET("/accounts/:accountId", log.LogMiddleware(), transaction.GetAccount(services, repos))
	routerGroup.POST("/transactions", log.LogMiddleware(), transaction.InsertAccount(services, repos))
}
