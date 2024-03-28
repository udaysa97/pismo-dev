package email

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pismo-dev/constants"
	downstreamerror "pismo-dev/error/downstreamError"
	"pismo-dev/pkg/httpclient"
	httpTypes "pismo-dev/pkg/httpclient/types"
	"pismo-dev/pkg/logger"
)

type EmailSvc struct {
	serviceName     string
	EmailSvcBaseUrl string
	httpClient      *httpclient.HttpClientWrapper
	vpcSecret       string
}

func NewEmailSvc(client *httpclient.HttpClientWrapper, emailSvcBaseURL string, vpcSecret string) *EmailSvc {
	return &EmailSvc{

		serviceName:     "EmailSVC",
		httpClient:      client,
		EmailSvcBaseUrl: emailSvcBaseURL,
		vpcSecret:       vpcSecret,
	}
}

func (svc *EmailSvc) SendMail(ctx context.Context, payload map[string]interface{}) (bool, error) {
	var options httpTypes.RequestOptions
	options.RetryAttempt = constants.DEFAULT_RETRY_ATTEMPT
	options.RetryFixedInternal = constants.DEFAULT_RETRY_FIXED_INTERVAL

	options.Url = fmt.Sprintf("%s%s", svc.EmailSvcBaseUrl, constants.ENDPOINTS["MAIL"]["SEND_MAIL"])

	var responseData *http.Response
	var err error
	options.Headers = map[string]string{
		"Content-Type":           "application/json",
		"x-authorization-secret": svc.vpcSecret,
	}

	options.Body, _ = json.Marshal(payload)

	responseData, err = svc.httpClient.Driver.Post(options)

	if err != nil {
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.serviceName, options.Url, err.Error())
		logger.Error("[SEND_MAIL] Failed to fetch data", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err, "request": options.Body})
		return false, err
	}
	if responseData.StatusCode != 200 {
		logger.Error("[SEND_MAIL] Failed to Send Mail", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "response": responseData, "request": options.Body})
		data, err := io.ReadAll(responseData.Body)
		if err != nil {
			err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.serviceName, options.Url, "[SEND_MAIL] Failed to fetch data")
			logger.Error("[SEND_MAIL] Failed to read response", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "response": responseData.Body, "error": err})
			return false, err
		}
		var responseBody interface{}
		err = json.Unmarshal(data, &responseBody)
		err = downstreamerror.New(constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, constants.ERROR_TYPES[constants.PROCESS_ERROR].HttpStatus, svc.serviceName, options.Url, err.Error())
		logger.Error("[SEND_MAIL] Error in response is", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "response": responseBody, "error": err})
		return false, err
	}
	return true, nil

}
