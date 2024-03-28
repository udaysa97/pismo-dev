package portfolio

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pismo-dev/constants"
	downstreamerror "pismo-dev/error/downstreamError"
	"pismo-dev/external/types"
	"pismo-dev/pkg/httpclient"
	httpTypes "pismo-dev/pkg/httpclient/types"
	"pismo-dev/pkg/logger"

	"strconv"
)

type PortfolioSvc struct {
	ServiceName      string
	httpClient       *httpclient.HttpClientWrapper
	portfolioBaseUrl string
	vpcSecret        string
}

func NewPortfolioSvc(client *httpclient.HttpClientWrapper, PortfolioBaseUrl string, vpcSecret string) *PortfolioSvc {
	return &PortfolioSvc{
		ServiceName:      "PortfolioSvc",
		httpClient:       client,
		portfolioBaseUrl: PortfolioBaseUrl,
		vpcSecret:        vpcSecret,
	}
}

func (svc *PortfolioSvc) GetUserBalance(ctx context.Context, nftId string, networkId string, userId string) (types.TokenBalanceResponse, error) {
	var tokenBalanceResponse types.TokenBalanceResponse
	var options httpTypes.RequestOptions
	options.RetryAttempt = constants.DEFAULT_RETRY_ATTEMPT
	options.RetryFixedInternal = constants.DEFAULT_RETRY_FIXED_INTERVAL
	portfolioEndpoint := fmt.Sprintf(constants.ENDPOINTS["PORTFOLIO"]["GetUserBalance"], userId)
	options.Url = fmt.Sprintf("%s%s", svc.portfolioBaseUrl, portfolioEndpoint)

	options.QueryParams = map[string]string{}
	var responseData *http.Response
	var err error
	options.Headers = map[string]string{
		"Accept":                 "application/json",
		"Content-Type":           "application/json",
		"x-authorization-secret": svc.vpcSecret,
	}

	options.QueryParams = map[string]string{
		"entityId":  nftId,
		"networkId": networkId,
	}

	responseData, err = svc.httpClient.Driver.Get(options)

	if err != nil {
		logger.Error("[GET_USER_BALANCE] Failed to fetch data", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.ServiceName, "GET_USER_BALANCE", "[GET_USER_BALANCE] Failed to fetch data")
		return tokenBalanceResponse, err
	}
	if responseData.StatusCode == 503 {
		logger.Error("[GET_USER_BALANCE] Failed to fetch data", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "response": responseData})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, constants.ERROR_TYPES[constants.PROCESS_ERROR].HttpStatus, svc.ServiceName, "GET_USER_BALANCE", "[GET_USER_BALANCE] Downstream service unavailable")
		return tokenBalanceResponse, err
	}

	data, err := io.ReadAll(responseData.Body)

	if err != nil {
		logger.Error("[GET_USER_BALANCE] Failed to read the response", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.ServiceName, "GET_USER_BALANCE", "[GET_USER_BALANCE] Failed to read the response")
		return tokenBalanceResponse, err
	}
	if err := json.Unmarshal(data, &tokenBalanceResponse); err != nil {
		logger.Error("[GET_USER_BALANCE] Failed to unmarshal the response", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err, "response": data})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.ServiceName, "GET_USER_BALANCE", "[GET_USER_BALANCE] Failed to unmarshal the response")
		return tokenBalanceResponse, err
	}
	if responseData.StatusCode >= 400 {
		logger.Error("[GET_USER_BALANCE] Request Failed", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err, "response": tokenBalanceResponse})
		if responseData.StatusCode == 500 {
			err = downstreamerror.New(constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, constants.ERROR_TYPES[constants.PROCESS_ERROR].HttpStatus, svc.ServiceName, "GET_USER_BALANCE", "[GET_USER_BALANCE] Request Failed")
		} else {
			err = downstreamerror.New(strconv.Itoa(400), responseData.StatusCode, svc.ServiceName, "GET_USER_BALANCE", tokenBalanceResponse.Error.Message) //TODO: Add error as per response struct
		}
		return tokenBalanceResponse, err
	}
	if !tokenBalanceResponse.Status {
		logger.Error("[GET_USER_BALANCE] Request Failed", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "response": tokenBalanceResponse})
		err = downstreamerror.New(strconv.Itoa(tokenBalanceResponse.Error.Code), responseData.StatusCode, svc.ServiceName, "GET_USER_BALANCE", tokenBalanceResponse.Error.Message) //TODO: Add error as per response struct
		return tokenBalanceResponse, err
	}
	if len(tokenBalanceResponse.Result.Rows) < 1 {
		logger.Error("[GET_USER_BALANCE] Request Failed", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "response": tokenBalanceResponse})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, responseData.StatusCode, svc.ServiceName, options.Url, "No User balance") //TODO: Add error as per response struct
		return tokenBalanceResponse, err
	}
	logger.Debug("PSUserBalanceResponse", map[string]interface{}{"context": ctx, "response": tokenBalanceResponse})

	return tokenBalanceResponse, nil

}
