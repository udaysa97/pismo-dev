package signing

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

type SigningSvc struct {
	ServiceName    string
	httpClient     *httpclient.HttpClientWrapper
	signingBaseUrl string
	vpcSecret      string
}

func NewSigningSvc(client *httpclient.HttpClientWrapper, signingBaseUrl string, vpcSecret string) *SigningSvc {
	return &SigningSvc{
		ServiceName:    "SigningSvc",
		httpClient:     client,
		signingBaseUrl: signingBaseUrl,
		vpcSecret:      vpcSecret,
	}
}

func (svc *SigningSvc) GetUserWalletAddress(ctx context.Context, userId, networkId string) (types.SigningSvcResponse, error) {
	var signingServiceResponse types.SigningSvcResponse
	var options httpTypes.RequestOptions
	options.RetryAttempt = constants.DEFAULT_RETRY_ATTEMPT
	options.RetryFixedInternal = constants.DEFAULT_RETRY_FIXED_INTERVAL
	options.Url = fmt.Sprintf("%s%s", svc.signingBaseUrl, constants.ENDPOINTS["SIGNING"]["GET_WALLET_ADDRESS"])

	options.QueryParams = map[string]string{}
	var responseData *http.Response
	var err error
	options.Headers = map[string]string{
		"Accept":                 "application/json",
		"Content-Type":           "application/json",
		"x-authorization-secret": svc.vpcSecret,
	}

	options.QueryParams = map[string]string{
		"user_id":               userId,
		"blockchain_network_id": networkId,
	}

	responseData, err = svc.httpClient.Driver.Get(options)

	if err != nil {
		logger.Error("[GET_USER_ADDRESS] Failed to fetch data", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.ServiceName, "GET_USER_ADDRESS", "[GET_USER_ADDRESS] Failed to fetch data")
		return types.SigningSvcResponse{}, err
	}
	if responseData.StatusCode == 503 {
		logger.Error("[GET_USER_ADDRESS] Failed to fetch data", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "response": responseData})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, constants.ERROR_TYPES[constants.PROCESS_ERROR].HttpStatus, svc.ServiceName, "GET_USER_ADDRESS", "[GET_USER_ADDRESS] Downstream service unavailable")
		return types.SigningSvcResponse{}, err
	}

	data, err := io.ReadAll(responseData.Body)

	if err != nil {
		logger.Error("[GET_USER_ADDRESS] Failed to read the response", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.ServiceName, "GET_USER_ADDRESS", "[GET_USER_ADDRESS] Failed to read the response")
		return types.SigningSvcResponse{}, err
	}
	if err := json.Unmarshal(data, &signingServiceResponse); err != nil {
		logger.Error("[GET_USER_ADDRESS] Failed to unmarshal the response", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err, "response": data})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.ServiceName, "GET_USER_ADDRESS", "[GET_USER_ADDRESS] Failed to unmarshal the response")
		return types.SigningSvcResponse{}, err
	}
	if responseData.StatusCode != 200 {
		logger.Error("[GET_USER_ADDRESS] Request Failed", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err, "response": signingServiceResponse})
		if responseData.StatusCode == 500 {
			err = downstreamerror.New(constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, constants.ERROR_TYPES[constants.PROCESS_ERROR].HttpStatus, svc.ServiceName, "GET_USER_ADDRESS", "[GET_USER_ADDRESS] Request Failed")
		} else {
			err = downstreamerror.New(strconv.Itoa(400), responseData.StatusCode, svc.ServiceName, "GET_USER_ADDRESS", "User address not found") //TODO: Add error as per response struct
		}
		return signingServiceResponse, err
	}
	if len(signingServiceResponse.Address) < 1 {
		logger.Error("[GET_USER_ADDRESS] Request Failed", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "response": data})
		return signingServiceResponse, err
	}
	logger.Debug("SigningService address", map[string]interface{}{"context": ctx, "response": signingServiceResponse})

	return signingServiceResponse, nil

}
