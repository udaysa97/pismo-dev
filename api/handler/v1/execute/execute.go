package execute

import (
	"errors"
	"fmt"
	"net/http"
	"pismo-dev/api/common"
	"pismo-dev/api/logger"
	"pismo-dev/api/types"
	device "pismo-dev/commonpkg/devicedetectionutils"
	validationerror "pismo-dev/error/validation"
	"pismo-dev/internal/service"
	resultType "pismo-dev/internal/types"

	"github.com/gin-gonic/gin"
)

func ValidateExecuteNFT(services *service.Service) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		requestBody := &types.OTPRequest{}
		response := types.ResponseDTO[string]{}
		var err error
		deviceDetails, err := device.SetDevice(ctx)
		if err != nil {
			errorMsg := fmt.Sprintf("validation: %s", err.Error())
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
		// Extract userId from request and add to above variable
		if err = common.ReadAndValidateRequestBody(ctx.Request, requestBody); err != nil {
			var validationError *validationerror.ValidationError
			if errors.As(err, &validationError) {
				logger.VpcErrorResponse(ctx, validationError.HttpCode, validationError.ErrorCode, validationError.Message)
				return
			}
			errorMsg := fmt.Sprintf("validation: %s", err)
			response.Status = types.StatusError
			response.Success = false
			errorResponse := types.ErrorResponse{}
			errorResponse.Code = http.StatusBadRequest
			errorResponse.ErrorCode = types.BADREQUEST
			errorResponse.Message = errorMsg
			response.Error = &errorResponse
			ctx.JSON(http.StatusBadRequest, response)
			return
		}
		tokenTransferObject := resultType.NFTTransferDetailsInterface{
			NftId:                    requestBody.NftID,
			NetworkId:                requestBody.NetworkId,
			DestinationWalletAddress: requestBody.RecipientWalletAddress,
			Amount:                   requestBody.Quantity,
			TransferType:             requestBody.OrderType,
		}

		user := resultType.UserDetailsInterface{
			Id:            requestBody.CurrentUser.Id,
			Source:        requestBody.CurrentUser.Source,
			ReloginPin:    requestBody.CurrentUser.ReloginPin,
			AuthToken:     requestBody.CurrentUser.AuthToken,
			DeviceDetails: deviceDetails,
		}

		// user := resultType.UserDetailsInterface {

		// };
		if _, err := services.OtpSvc.GenerateOTP(ctx, user, tokenTransferObject); err != nil {
			return
		}
		response.Status = types.StatusSuccess
		response.Success = true
		ctx.JSON(http.StatusOK, response)

	}
}

func ExecuteNFT(services *service.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestBody := &types.ExecuteRequest{}
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
		tokenTransferObject := resultType.NFTTransferDetailsInterface{
			NftId:                    requestBody.NftID,
			NetworkId:                requestBody.NetworkId,
			DestinationWalletAddress: requestBody.RecipientWalletAddress,
			Amount:                   requestBody.Quantity,
			TransferType:             requestBody.OrderType,
			GsnIncludeToken:          requestBody.GsnIncludeToken,
			IsGsnRequired:            requestBody.IsGsnRequired,
			GsnIncludeNetworkId:      requestBody.GsnIncludeNetworkId,
			GsnIncludeMaxAmount:      requestBody.GsnIncludeMaxAmount,
			ErcType:                  requestBody.ErcType,
		}

		user := resultType.UserDetailsInterface{
			Id:         requestBody.CurrentUser.Id,
			Source:     requestBody.CurrentUser.Source,
			ReloginPin: requestBody.CurrentUser.ReloginPin,
			AuthToken:  requestBody.CurrentUser.AuthToken,
			UserOTP:    requestBody.CurrentUser.UserOTP,
			DeviceId:   requestBody.CurrentUser.DeviceId,
		}
		if executeResponse, err = services.ExecutionSvc.ExecuteOrder(ctx, user, tokenTransferObject); err != nil {
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
