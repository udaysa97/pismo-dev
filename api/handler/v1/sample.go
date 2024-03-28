package v1

import (
	"net/http"
	"pismo-dev/api/types"

	"github.com/gin-gonic/gin"
)

func SampleHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		response := types.SampleResponse{
			Status:  "success",
			Message: "Working as expected",
		}
		ctx.JSON(http.StatusOK, response)
	}
}
