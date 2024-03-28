package route

import (
	"net/http"
	"pismo-dev/api/handler"
	"pismo-dev/api/middleware/auth"
	"pismo-dev/api/middleware/log"
	openApi "pismo-dev/api/route/openapi/v1"
	"pismo-dev/api/route/vpc"

	"pismo-dev/internal/service"

	v1 "pismo-dev/api/route/v1"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.Engine, services *service.Service) {

	// Base API routes
	router.GET("/", func(ctx *gin.Context) {
		response := gin.H{
			"message": "Served from okto NFT MS",
			"status":  http.StatusOK,
			"code":    "200",
		}
		ctx.JSON(http.StatusOK, response)
	})

	router.GET("/health", handler.GetHealth)

	v1ApiGroup := router.Group("/api/v1")
	v1.Register(v1ApiGroup, services)

	vpcApiGroup := router.Group("/vpc", log.LogMiddleware(), auth.VPCAuthorizationMiddleware())
	vpc.Register(vpcApiGroup, services)

	openApiGroup := router.Group("/openapi", log.LogMiddleware(), auth.OpenAPIVPCAuthorizationMiddleware())
	openApi.Register(openApiGroup, services)

}
