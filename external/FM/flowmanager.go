package FM

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pismo-dev/constants"
	downstreamerror "pismo-dev/error/downstreamError"
	"pismo-dev/internal/types"
	"pismo-dev/pkg/httpclient"
	httpTypes "pismo-dev/pkg/httpclient/types"
	"pismo-dev/pkg/logger"
)

type FMSvc struct {
	serviceName string
	fmBaseURL   string
	dappBaseUrl string
	httpClient  *httpclient.HttpClientWrapper
	vpcSecret   string
}

func NewFMSvc(client *httpclient.HttpClientWrapper, FMBaseURL string, dappUrl, vpcSecret string) *FMSvc {
	return &FMSvc{

		serviceName: "FMSvc",
		httpClient:  client,
		fmBaseURL:   FMBaseURL,
		vpcSecret:   vpcSecret,
		dappBaseUrl: dappUrl,
	}
}

func (svc *FMSvc) GetEstimate(ctx context.Context, payload types.FMEstimateRequest) (types.FMEstimateResult, error) {
	var response types.FMEstimateResult
	var fmResponse types.FMEstimateResponse
	var options httpTypes.RequestOptions
	options.RetryAttempt = constants.DEFAULT_RETRY_ATTEMPT
	options.RetryFixedInternal = constants.DEFAULT_RETRY_FIXED_INTERVAL

	options.Url = fmt.Sprintf("%s%s", svc.fmBaseURL, constants.ENDPOINTS["FM"]["Estimate"])

	var responseData *http.Response
	var err error
	options.Headers = map[string]string{
		"Accept":                 "application/json",
		"Content-Type":           "application/json",
		"x-authorization-secret": svc.vpcSecret,
	}
	options.Body, _ = json.Marshal(payload)

	responseData, err = svc.httpClient.Driver.Post(options)

	if err != nil {
		logger.Error("[GET_NFT_TRANSFER_ESTIMATE] Failed to fetch data", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.serviceName, "GET_NFT_TRANSFER_ESTIMATE", "[GET_NFT_TRANSFER_ESTIMATE] Failed to fetch data")
		return response, err
	}
	if responseData.StatusCode == 503 {
		logger.Error("[GET_NFT_TRANSFER_ESTIMATE] Failed to fetch data", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "response": responseData})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, constants.ERROR_TYPES[constants.PROCESS_ERROR].HttpStatus, svc.serviceName, "GET_NFT_TRANSFER_ESTIMATE", "[GET_NFT_TRANSFER_ESTIMATE] Downstream service unavailable")
		return response, err
	}

	data, err := io.ReadAll(responseData.Body)

	if err != nil {
		logger.Error("[GET_NFT_TRANSFER_ESTIMATE] Failed to read the response", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.serviceName, "GET_NFT_TRANSFER_ESTIMATE", "[GET_NFT_TRANSFER_ESTIMATE] Failed to read the response")
		return response, err
	}
	if err := json.Unmarshal(data, &fmResponse); err != nil {
		logger.Error("[GET_NFT_TRANSFER_ESTIMATE] Failed to unmarshal the response", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err.Error(), "response": data})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.serviceName, "GET_NFT_TRANSFER_ESTIMATE", "[GET_NFT_TRANSFER_ESTIMATE] Failed to unmarshal the response")
		return response, err
	}
	if responseData.StatusCode >= 400 {
		logger.Error("[GET_NFT_TRANSFER_ESTIMATE] Request Failed", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "response": fmResponse})
		if responseData.StatusCode == 500 {
			err = downstreamerror.New(constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, constants.ERROR_TYPES[constants.PROCESS_ERROR].HttpStatus, svc.serviceName, "GET_NFT_TRANSFER_ESTIMATE", "[GET_NFT_TRANSFER_ESTIMATE] Request Failed")
		} else {
			err = downstreamerror.New(fmResponse.Result.Error.ErrorCode, responseData.StatusCode, svc.serviceName, "GET_NFT_TRANSFER_ESTIMATE", fmResponse.Result.Error.Message)
		}
		return fmResponse.Result, err
	}
	if fmResponse.Status == "FAILED" || !fmResponse.Result.Output.Success {
		logger.Error("[GET_NFT_TRANSFER_ESTIMATE] Failed Status from FM", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "response": fmResponse})
		return fmResponse.Result, fmt.Errorf("Fm response failed status")
	}
	logger.Debug("FMEstimateResponse", map[string]interface{}{"context": ctx, "response": fmResponse})

	return fmResponse.Result, nil
}

func (svc *FMSvc) ExecuteOrder(ctx context.Context, payload types.FMExecuteRequest) (types.FMExecuteResult, error) {
	var response types.FMExecuteResult
	var fmResponse types.FMExecuteResponse
	var options httpTypes.RequestOptions
	options.RetryAttempt = constants.DEFAULT_RETRY_ATTEMPT
	options.RetryFixedInternal = constants.DEFAULT_RETRY_FIXED_INTERVAL

	options.Url = fmt.Sprintf("%s%s", svc.fmBaseURL, constants.ENDPOINTS["FM"]["Execute"])

	var responseData *http.Response
	var err error
	options.Headers = map[string]string{
		"Accept":                 "application/json",
		"Content-Type":           "application/json",
		"x-authorization-secret": svc.vpcSecret,
	}
	options.Body, _ = json.Marshal(payload)

	responseData, err = svc.httpClient.Driver.Post(options)

	if err != nil {
		logger.Error("[EXECUTE_NFT_TRANSFER] Failed to fetch data", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err.Error(), "request": options.Body})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.serviceName, "EXECUTE_NFT_TRANSFER", "[EXECUTE_NFT_TRANSFER] Failed to fetch data")
		return response, err
	}
	if responseData.StatusCode == 503 {
		logger.Error("[EXECUTE_NFT_TRANSFER] Failed to fetch data", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "response": responseData, "request": options.Body})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, constants.ERROR_TYPES[constants.PROCESS_ERROR].HttpStatus, svc.serviceName, "EXECUTE_NFT_TRANSFER", "[EXECUTE_NFT_TRANSFER] Downstream service unavailable")
		return response, err
	}

	data, err := io.ReadAll(responseData.Body)

	if err != nil {
		logger.Error("[EXECUTE_NFT_TRANSFER] Failed to read the response", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err.Error(), "request": options.Body})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.serviceName, "EXECUTE_NFT_TRANSFER", "[EXECUTE_NFT_TRANSFER] Failed to read the response")
		return response, err
	}
	if err := json.Unmarshal(data, &fmResponse); err != nil {
		logger.Error("[EXECUTE_NFT_TRANSFER] Failed to unmarshal the response", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err.Error(), "response": data})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.serviceName, "EXECUTE_NFT_TRANSFER", "[EXECUTE_NFT_TRANSFER] Failed to unmarshal the response")
		return response, err
	}
	if (fmResponse.Result != types.FMExecuteResult{}) {
		response = fmResponse.Result
	}
	if responseData.StatusCode >= 400 {
		logger.Error("[EXECUTE_NFT_TRANSFER] Request Failed", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "response": fmResponse, "request": options.Body})
		if responseData.StatusCode == 500 {
			err = downstreamerror.New(constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, constants.ERROR_TYPES[constants.PROCESS_ERROR].HttpStatus, svc.serviceName, "EXECUTE_NFT_TRANSFER", "[EXECUTE_NFT_TRANSFER] Request Failed")
		} else {
			err = downstreamerror.New(fmResponse.Result.Error.ErrorCode, responseData.StatusCode, svc.serviceName, "EXECUTE_NFT_TRANSFER", fmResponse.Result.Error.Message)
		}
		return response, err
	}
	if fmResponse.Status == "FAILED" {
		logger.Error("[EXECUTE_NFT_TRANSFER] Failed Status from FM", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "response": fmResponse})
		return response, fmt.Errorf("Fm response failed status")
	}
	logger.Debug("FMExecuteResponse", map[string]interface{}{"context": ctx, "response": fmResponse})

	return fmResponse.Result, nil

}

func (svc *FMSvc) GetStatus(ctx context.Context, jobId string) (types.FMGetStatusResult, error) {
	var response types.FMGetStatusResult
	var fmResponse types.FMGetStatusResponse
	var options httpTypes.RequestOptions
	options.RetryAttempt = constants.DEFAULT_RETRY_ATTEMPT
	options.RetryFixedInternal = constants.DEFAULT_RETRY_FIXED_INTERVAL

	getStatusEndpoint := fmt.Sprintf(constants.ENDPOINTS["FM"]["GetStatus"], jobId)
	options.Url = fmt.Sprintf("%s%s", svc.fmBaseURL, getStatusEndpoint)
	var responseData *http.Response
	var err error

	options.Headers = map[string]string{
		"Accept":                 "application/json",
		"Content-Type":           "application/json",
		"x-authorization-secret": svc.vpcSecret,
	}

	responseData, err = svc.httpClient.Driver.Get(options)

	if err != nil {
		logger.Error("[GET_NFT_TRANSFER_STATUS] Failed to fetch data", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err.Error()})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.serviceName, "GET_NFT_TRANSFER_STATUS", "[GET_NFT_TRANSFER_STATUS] Failed to fetch data")
		return response, err
	}
	if responseData.StatusCode == 503 {
		logger.Error("[GET_NFT_TRANSFER_STATUS] Failed to fetch data", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "response": responseData})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, constants.ERROR_TYPES[constants.PROCESS_ERROR].HttpStatus, svc.serviceName, "GET_NFT_TRANSFER_STATUS", "[GET_NFT_TRANSFER_STATUS] Downstream service unavailable")
		return response, err
	}

	data, err := io.ReadAll(responseData.Body)

	if err != nil {
		logger.Error("[GET_NFT_TRANSFER_STATUS] Failed to read the response", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err.Error()})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.serviceName, "GET_NFT_TRANSFER_STATUS", "[GET_NFT_TRANSFER_STATUS] Failed to read the response")
		return response, err
	}
	if err := json.Unmarshal(data, &fmResponse); err != nil {
		logger.Error("[GET_NFT_TRANSFER_STATUS] Failed to unmarshal the response", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err.Error(), "response": data})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.serviceName, "GET_NFT_TRANSFER_STATUS", "[GET_NFT_TRANSFER_STATUS] Failed to unmarshal the response")
		return response, err
	}
	if responseData.StatusCode >= 400 {
		logger.Error("[GET_NFT_TRANSFER_STATUS] Request Failed", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "response": fmResponse})
		if responseData.StatusCode == 500 {
			err = downstreamerror.New(constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, constants.ERROR_TYPES[constants.PROCESS_ERROR].HttpStatus, svc.serviceName, "GET_NFT_TRANSFER_STATUS", "[GET_NFT_TRANSFER_STATUS] Request Failed")
		} else {
			err = downstreamerror.New(fmResponse.Result.Error.ErrorCode, responseData.StatusCode, svc.serviceName, "GET_NFT_TRANSFER_STATUS", fmResponse.Result.Error.Message)
		}
		return response, err
	}
	logger.Debug("FMStatusResponse", map[string]interface{}{"context": ctx, "response": fmResponse})

	return fmResponse.Result, nil

}

func (svc *FMSvc) ExecuteOpenAPIOrder(ctx context.Context, payload types.OpenAPINFTOrder, orderType string, vendorId string) (types.FMExecuteResult, error) {
	var response types.FMExecuteResult
	var fmResponse types.FMExecuteResponse
	var options httpTypes.RequestOptions
	options.RetryAttempt = constants.DEFAULT_RETRY_ATTEMPT
	options.RetryFixedInternal = constants.DEFAULT_RETRY_FIXED_INTERVAL

	options.Url = fmt.Sprintf("%s%s", svc.dappBaseUrl, fmt.Sprintf(constants.ENDPOINTS["DAPP"]["Execute"], orderType))

	var responseData *http.Response
	var err error
	options.Headers = map[string]string{
		"Accept":                 "application/json",
		"Content-Type":           "application/json",
		"x-authorization-secret": svc.vpcSecret,
		"vendor-id":              vendorId,
	}
	options.Body, _ = json.Marshal(payload)
	logger.Info("Dapp execute url", options.Url)
	logger.Info("Dapp execute body", string(options.Body))
	responseData, err = svc.httpClient.Driver.Post(options)
	logger.Info("Dapp execute Response", responseData)
	if err != nil {
		logger.Error("[EXECUTE_OPENAPI_NFT_TRANSFER] Failed to fetch data", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err.Error(), "request": options.Body})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.serviceName, "EXECUTE_NFT_TRANSFER", "[EXECUTE_NFT_TRANSFER] Failed to fetch data")
		return response, err
	}
	if responseData.StatusCode == 503 {
		logger.Error("[EXECUTE_OPENAPI_NFT_TRANSFER] Failed to fetch data", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "response": responseData, "request": options.Body})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, constants.ERROR_TYPES[constants.PROCESS_ERROR].HttpStatus, svc.serviceName, "EXECUTE_OPENAPI_NFT_TRANSFER", "[EXECUTE_OPENAPI_NFT_TRANSFER] Downstream service unavailable")
		return response, err
	}

	data, err := io.ReadAll(responseData.Body)

	if err != nil {
		logger.Error("[EXECUTE_OPENAPI_NFT_TRANSFER] Failed to read the response", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err.Error(), "request": options.Body})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.serviceName, "EXECUTE_OPENAPI_NFT_TRANSFER", "[EXECUTE_OPENAPI_NFT_TRANSFER] Failed to read the response")
		return response, err
	}
	if err := json.Unmarshal(data, &fmResponse); err != nil {
		logger.Error("[EXECUTE_OPENAPI_NFT_TRANSFER] Failed to unmarshal the response", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err.Error(), "response": data})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.serviceName, "EXECUTE_OPENAPI_NFT_TRANSFER", "[EXECUTE_OPENAPI_NFT_TRANSFER] Failed to unmarshal the response")
		return response, err
	}
	if (fmResponse.Result != types.FMExecuteResult{}) {
		response = fmResponse.Result
	}
	if responseData.StatusCode >= 400 {
		logger.Error("[EXECUTE_OPENAPI_NFT_TRANSFER] Request Failed", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "response": fmResponse, "request": options.Body})
		if responseData.StatusCode == 500 {
			err = downstreamerror.New(constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, constants.ERROR_TYPES[constants.PROCESS_ERROR].HttpStatus, svc.serviceName, "EXECUTE_OPENAPI_NFT_TRANSFER", "[EXECUTE_OPENAPI_NFT_TRANSFER] Request Failed")
		} else {
			err = downstreamerror.New(fmResponse.Result.Error.ErrorCode, responseData.StatusCode, svc.serviceName, "EXECUTE_NFT_TRANSFER", fmResponse.Result.Error.Message)
		}
		return response, err
	}
	if fmResponse.Status == "FAILED" {
		logger.Error("[EXECUTE_OPENAPI_NFT_TRANSFER] Failed Status from FM", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "response": fmResponse})
		return response, fmt.Errorf("Fm response failed status")
	}
	logger.Debug("FMExecuteResponse", map[string]interface{}{"context": ctx, "response": fmResponse})

	return fmResponse.Result, nil

}
