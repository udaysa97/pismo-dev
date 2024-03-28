package execute

import (
	"errors"
	"fmt"
	"net/http"
	"pismo-dev/api/common"
	"pismo-dev/api/logger"
	"pismo-dev/api/types"
	validationerror "pismo-dev/error/validation"
	"pismo-dev/internal/service"
	resultType "pismo-dev/internal/types"

	"github.com/gin-gonic/gin"
)

func ExecuteNFT(services *service.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestBody := &types.OpenApiExecuteRequest{}
		response := types.ResponseDTO[types.ExecuteResponse]{}
		executeResponse := resultType.FMExecuteResult{}
		var err error
		// Extract userId from request and add to above variable
		if err = common.ReadAndValidateRequestBody(ctx.Request, requestBody); err != nil {
			var validationError *validationerror.ValidationError
			if errors.As(err, &validationError) {
				logger.VpcErrorResponse(ctx, validationError.HttpCode, validationError.ErrorCode, validationError.Message)
				return
			}
			errorMsg := fmt.Sprintf("validation: %s", err)
			response.Status = types.StatusError
			errorResponse := types.ErrorResponse{}
			errorResponse.Code = http.StatusBadRequest
			errorResponse.ErrorCode = types.BADREQUEST
			errorResponse.Message = errorMsg
			response.Error = &errorResponse
			response.Success = false
			ctx.JSON(http.StatusBadRequest, response)
			return
		}

		if executeResponse, err = services.ExecutionSvc.OpenApiExecuteOrder(ctx, *requestBody); err != nil {
			return
		}
		resultantJobId := types.ExecuteResponse{}
		resultantJobId.OrderId = executeResponse.JobId
		response.Result = &resultantJobId
		response.Status = types.StatusSuccess
		response.Success = true
		ctx.JSON(http.StatusOK, response)

	}
}
