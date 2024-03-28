package vpc

import (
	"pismo-dev/api/handler/v1/mint"
	orderservice "pismo-dev/api/handler/v1/orders"
	transactionservice "pismo-dev/api/handler/v1/transactions"

	"pismo-dev/internal/service"

	"github.com/gin-gonic/gin"
)

func Register(routerGroup *gin.RouterGroup, services *service.Service) {
	v1ApiGroup := routerGroup.Group("/api/v1")
	v1ApiGroup.GET("/orderdetails", orderservice.GetOrderDetails(services))
	v1ApiGroup.GET("/transactiondetails", transactionservice.GetTransactionDetails(services))
	v1ApiGroup.POST("/nft/mint", mint.Mint(services))
}
