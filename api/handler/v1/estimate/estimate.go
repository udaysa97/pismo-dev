package estimate

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
	"github.com/google/uuid"
)

func EstimateNFT(services *service.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestBody := &types.EstimateRequest{}
		response := types.ResponseDTO[types.EstimationResponseWrapperInterface]{}
		responseBody := &types.EstimateResponse{}
		var estimateResponse resultType.FMEstimateResultOutput
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

		orderId := uuid.New().String()
		requestBody.OrderId = orderId // This orderId will only be used to debug
		tokenTransferObject := resultType.NFTTransferDetailsInterface{
			NftId:                    requestBody.NftID,
			NetworkId:                requestBody.NetworkId,
			DestinationWalletAddress: requestBody.RecipientWalletAddress,
			Amount:                   requestBody.Quantity,
			TransferType:             requestBody.OrderType,
			ErcType:                  requestBody.ErcType,
			OrderId:                  orderId,
		}

		user := resultType.UserDetailsInterface{
			Id:     requestBody.CurrentUser.Id,
			Source: requestBody.CurrentUser.Source,
		}
		if estimateResponse, err = services.EstimationSvc.CalculateEstimate(ctx, user, tokenTransferObject); err != nil {
			return
		}
		responseBody.TransactionFee = estimateResponse.TransactionFee
		responseBody.GsnWithdrawTokens = estimateResponse.GsnWithdrawTokens
		responseBody.IsGsnPossible = estimateResponse.IsGsnPossible
		responseBody.IsGsnRequired = estimateResponse.IsGsnRequired
		responseBody.OrderId = requestBody.OrderId
		var estimationWrapper types.EstimationResponseWrapperInterface
		estimationWrapper.Estimation = *responseBody
		response.Result = &estimationWrapper
		response.Status = types.StatusSuccess
		response.Success = true
		ctx.JSON(http.StatusOK, response)
	}
}
