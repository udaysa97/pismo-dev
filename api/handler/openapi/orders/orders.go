package orders

import (
	"net/http"
	"pismo-dev/api/logger"
	"pismo-dev/api/types"
	"pismo-dev/internal/repository/utils"
	"pismo-dev/internal/service"
	"pismo-dev/internal/util"

	"github.com/gin-gonic/gin"
)

func GetOrderDetails(services *service.Service) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		var queryParams map[string][]string
		queryParams = ctx.Request.URL.Query()
		response := types.ResponseDTO[types.OrderDetailsResponse]{}
		responseBody := &types.OrderDetailsResponse{}

		pagination := utils.GeneratePaginationFromRequest(ctx)
		queryParams["user_id"] = []string{ctx.Param("userid")}
		if result, err := services.OrderSvc.GetOrderDetails(ctx, queryParams, pagination); err != nil {
			if errInstance, isServiceError := util.GetServiceError(err); isServiceError {
				logger.ServiceErrorResponse(ctx, errInstance)
				return
			}
			logger.ErrorResponse(ctx, err.ErrorCode, err.Message)
			logger.Error(ctx, "error :", err)
			return
		} else {
			responseBody.OrderDetails = result
			response.Result = responseBody
			response.Status = types.StatusSuccess
			ctx.JSON(http.StatusOK, response)
		}
	}
}
