package api

import (
	"fmt"
	"net/http"
	"os"
	"pismo-dev/api/middleware/cors"
	"pismo-dev/api/route"
	"pismo-dev/constants"
	"pismo-dev/internal/appconfig"
	"pismo-dev/internal/repository"
	"pismo-dev/internal/service"
	"pismo-dev/pkg/logger"

	"github.com/gin-gonic/gin"
	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
)

func InitServer(services *service.Service, repos *repository.Repositories) {
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.HandleMethodNotAllowed = true

	router.Use(cors.CORSMiddleware())
	router.Use(gintrace.Middleware("pismo-dev"))

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"status": "error", "error": map[string]any{
			"Code":      404,
			"ErrorCode": "nometh",
			"Message":   "Method not found",
			"TraceID":   c.GetString(constants.TRACE_ID_KEY),
		}})
	})

	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"status": "error", "error": map[string]any{
			"Code":      404,
			"ErrorCode": "nometh",
			"Message":   "Method not found",
			"TraceID":   c.GetString(constants.TRACE_ID_KEY),
		}})
	})
	port := appconfig.PORT
	host := appconfig.HOST
	logger.Info(fmt.Sprintf("port: %s host: %s ", port, host))

	if len(port) == 0 {
		panic("env variable PORT is missing")
	}

	if len(host) == 0 {
		panic("env variable HOST is missing")
	}

	addr := fmt.Sprintf("%s:%s", host, port)

	logger.Info(fmt.Sprintf("Listening and serving HTTP on %s", addr))
	route.Register(router, services, repos)

	if err := router.Run(addr); err != nil {
		panic(err)
	}

}
