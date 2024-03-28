package orderservice

import (
	"errors"
	"fmt"
	"net/http"
	"pismo-dev/api/logger"
	"pismo-dev/api/types"
	validationerror "pismo-dev/error/validation"
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

func GetUserCollectionMintCounts(services *service.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := types.GetUserCollectionMintCountsRequest{}
		request.UserId = ctx.Query("user_id")
		request.CollectionAddress = ctx.Query("collection_address")
		response := types.ResponseDTO[types.GetUserCollectionMintCountsResponse]{}

		if err := request.Validate(); err != nil {
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

		if result, err := services.OrderSvc.GetUserCollectionMintCounts(ctx, request.UserId, request.CollectionAddress); err != nil {
			if errInstance, isServiceError := util.GetServiceError(err); isServiceError {
				logger.ServiceErrorResponse(ctx, errInstance)
				return
			}
			logger.ErrorResponse(ctx, err.ErrorCode, err.Message)
			logger.Error(ctx, "error :", err)
			return
		} else {
			response.Status = types.StatusSuccess
			response.Success = true
			response.Result = &result
			ctx.JSON(http.StatusOK, response)
		}

	}
}

func GetUserOrderDetails(services *service.Service) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		var queryParams map[string][]string
		userId := ctx.Param("userid")
		queryParams = ctx.Request.URL.Query()
		queryParams["userId"] = []string{userId}
		response := types.ResponseDTO[types.OrderDetailsResponse]{}
		responseBody := &types.OrderDetailsResponse{}

		pagination := utils.GeneratePaginationFromRequest(ctx)

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
