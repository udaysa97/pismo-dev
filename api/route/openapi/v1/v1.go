package v1

import (
	"pismo-dev/api/handler/openapi/execute"
	"pismo-dev/api/handler/openapi/mint"
	"pismo-dev/api/middleware/log"

	orderservice "pismo-dev/api/handler/openapi/orders"
	"pismo-dev/internal/service"

	"github.com/gin-gonic/gin"
)

// Auth error code will be handled separately under identifier "DEFAULT"
func Register(routerGroup *gin.RouterGroup, services *service.Service) {

	routerGroup.POST("/nft/execute", log.LogMiddleware(), execute.ExecuteNFT(services))
	routerGroup.POST("/nft/mint", log.LogMiddleware(), mint.MintNFT(services))
	routerGroup.GET("/:userid/nft/orderDetails", log.LogMiddleware(), orderservice.GetOrderDetails(services))

}
