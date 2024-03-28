package v1

import (
	v1 "pismo-dev/api/handler/v1"
	"pismo-dev/api/handler/v1/estimate"
	"pismo-dev/api/handler/v1/execute"
	"pismo-dev/api/handler/v1/mint"
	orderservice "pismo-dev/api/handler/v1/orders"
	"pismo-dev/api/middleware/log"

	"pismo-dev/internal/service"

	"github.com/gin-gonic/gin"
)

// Auth error code will be handled separately under identifier "DEFAULT"
func Register(routerGroup *gin.RouterGroup, services *service.Service) {
	routerGroup.GET("/sample", log.LogMiddleware(), v1.SampleHandler())
	routerGroup.POST("/nft/estimate", log.LogMiddleware(), estimate.EstimateNFT(services))
	routerGroup.POST("/nft/initiate-execute", log.LogMiddleware(), execute.ValidateExecuteNFT(services))
	routerGroup.POST("/nft/execute", log.LogMiddleware(), execute.ExecuteNFT(services))
	routerGroup.POST("/orderDetails", log.LogMiddleware(), orderservice.GetOrderDetails(services))
	routerGroup.POST("/nft/mint", log.LogMiddleware(), mint.Mint(services))
	routerGroup.GET("/nft/user-collection-mint-count", log.LogMiddleware(), orderservice.GetUserCollectionMintCounts(services))

}
