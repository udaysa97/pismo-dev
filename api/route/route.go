package route

import (
	"net/http"
	"pismo-dev/api/handler"

	"pismo-dev/internal/repository"
	"pismo-dev/internal/service"

	v1 "pismo-dev/api/route/v1"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.Engine, services *service.Service, repos *repository.Repositories) {

	// Base API routes
	router.GET("/", func(ctx *gin.Context) {
		response := gin.H{
			"message": "Served from Prismo",
			"status":  http.StatusOK,
			"code":    "200",
		}
		ctx.JSON(http.StatusOK, response)
	})

	router.GET("/health", handler.GetHealth)

	v1ApiGroup := router.Group("/api/v1")
	v1.Register(v1ApiGroup, services, repos)

}
