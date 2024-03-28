package mint

import (
	"errors"
	"fmt"
	"net/http"
	"pismo-dev/api/common"
	"pismo-dev/api/logger"
	"pismo-dev/api/types"
	validationerror "pismo-dev/error/validation"
	"pismo-dev/internal/service"

	"github.com/gin-gonic/gin"
)

func Mint(services *service.Service) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		requestBody := &types.MintOGRequest{}
		response := types.ResponseDTO[types.MintOGResponse]{}
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
		orderId, err := services.OgMintSvc.PushOrderForProcessing(ctx, requestBody.UserId, requestBody.ContractId, requestBody.NetworkId, requestBody.OperationType, requestBody.CustomDataUri)
		if err != nil {
			return
		}

		resultantJobId := types.MintOGResponse{}
		resultantJobId.OrderId = orderId
		response.Result = &resultantJobId
		response.Status = types.StatusSuccess
		response.Success = true
		ctx.JSON(http.StatusOK, response)

	}
}
