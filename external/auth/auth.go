package auth

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

type AuthSvc struct {
	serviceName string
	authBaseURL string
	httpClient  *httpclient.HttpClientWrapper
	vpcSecret   string
}

func NewAuthSvc(client *httpclient.HttpClientWrapper, authBaseURL string, vpcSecret string) *AuthSvc {
	return &AuthSvc{

		serviceName: "AuthSVC",
		httpClient:  client,
		authBaseURL: authBaseURL,
		vpcSecret:   vpcSecret,
	}
}

func (svc *AuthSvc) VerifyReloginPin(ctx context.Context, user types.UserDetailsInterface) (bool, error) {
	var options httpTypes.RequestOptions
	options.RetryAttempt = constants.DEFAULT_RETRY_ATTEMPT
	options.RetryFixedInternal = constants.DEFAULT_RETRY_FIXED_INTERVAL

	options.Url = fmt.Sprintf("%s%s", svc.authBaseURL, constants.ENDPOINTS["AUTH"]["VERIFY_RELOGIN_PIN"])

	var responseData *http.Response
	var err error
	options.Headers = map[string]string{
		"Accept":                 "application/json",
		"Content-Type":           "application/json",
		"x-authorization-secret": svc.vpcSecret,
		"Authorization":          user.AuthToken,
	}

	payload := map[string]interface{}{
		"coindcx_id":  user.Id,
		"relogin_pin": user.ReloginPin,
	}
	options.Body, _ = json.Marshal(payload)

	responseData, err = svc.httpClient.Driver.Post(options)

	if err != nil {
		logger.Error("[VERIFY_RELOGIN_PIN] Failed to fetch data", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err, "userId": user.Id})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.serviceName, "VERIFY_RELOGIN_PIN", "[VERIFY_RELOGIN_PIN] Failed to fetch data")
		return false, err
	}
	data, err := io.ReadAll(responseData.Body)
	if err != nil {
		logger.Error("[VERIFY_RELOGIN_PIN] Failed to read response", map[string]interface{}{"context": ctx, "response": responseData.Body, "userId": user.Id})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.serviceName, "VERIFY_RELOGIN_PIN", "[VERIFY_RELOGIN_PIN] Failed to fetch data")
		return false, err
	}
	var responseBody map[string]interface{}
	err = json.Unmarshal(data, &responseBody)
	if err != nil {
		logger.Error("[VERIFY_RELOGIN_PIN] Error in unmarshalling response", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "response": responseBody, "userId": user.Id, "error": err.Error()})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, constants.ERROR_TYPES[constants.PROCESS_ERROR].HttpStatus, svc.serviceName, "VERIFY_RELOGIN_PIN", "[VERIFY_RELOGIN_PIN] Cannot unmarshal data")
		return false, err
	}
	if responseBody["code"] != float64(200) {
		logger.Error("[VERIFY_RELOGIN_PIN] Invalid PIN", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "response": responseBody, "userId": user.Id})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, constants.ERROR_TYPES[constants.PROCESS_ERROR].HttpStatus, svc.serviceName, "VERIFY_RELOGIN_PIN", "[VERIFY_RELOGIN_PIN] INVALID PIN")
		return false, err
	}
	return true, nil

}

func (svc *AuthSvc) ForceLogout(ctx context.Context, user types.UserDetailsInterface) {
	var options httpTypes.RequestOptions
	options.RetryAttempt = constants.DEFAULT_RETRY_ATTEMPT
	options.RetryFixedInternal = constants.DEFAULT_RETRY_FIXED_INTERVAL

	options.Url = fmt.Sprintf("%s%s", svc.authBaseURL, constants.ENDPOINTS["AUTH"]["FORCE_LOGOUT"])

	var responseData *http.Response
	var err error
	options.Headers = map[string]string{
		"Accept":                 "application/json",
		"Content-Type":           "application/json",
		"x-authorization-secret": svc.vpcSecret,
		"Authorization":          user.AuthToken,
	}

	payload := map[string]interface{}{
		"coindcx_id": user.Id,
	}
	options.Body, _ = json.Marshal(payload)

	responseData, err = svc.httpClient.Driver.Post(options)

	if err != nil {
		logger.Error("[VERIFY_RELOGIN_PIN] Failed to fetch data", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err, "userId": user.Id})
	}
	if responseData.StatusCode != 200 {
		logger.Error("[VERIFY_RELOGIN_PIN] Failed to fetch data", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "response": responseData, "userId": user.Id})
	}

}
