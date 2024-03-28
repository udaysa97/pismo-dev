package DQL

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"pismo-dev/constants"
	downstreamerror "pismo-dev/error/downstreamError"
	"pismo-dev/internal/types"
	"pismo-dev/pkg/httpclient"
	httpTypes "pismo-dev/pkg/httpclient/types"
	"pismo-dev/pkg/logger"
)

type DQLSvc struct {
	serviceName string
	dqlBaseURL  string
	httpClient  *httpclient.HttpClientWrapper
	vpcSecret   string
}

func NewDQLSvc(client *httpclient.HttpClientWrapper, dqlBaseURL string, vpcSecret string) *DQLSvc {
	return &DQLSvc{

		serviceName: "DQLSvc",
		httpClient:  client,
		dqlBaseURL:  dqlBaseURL,
		vpcSecret:   vpcSecret,
	}
}

func (svc *DQLSvc) GetEntityById(ctx context.Context, entityId string, isNft bool) (types.DQLEntityResponseInterface, error) {
	var dqlResponse types.DQLEntityResponseInterface
	var options httpTypes.RequestOptions
	options.RetryAttempt = constants.DEFAULT_RETRY_ATTEMPT
	options.RetryFixedInternal = constants.DEFAULT_RETRY_FIXED_INTERVAL

	dqlEndpoint := fmt.Sprintf(constants.ENDPOINTS["DQL"]["GetEntityById"], entityId)
	options.Url = fmt.Sprintf("%s%s", svc.dqlBaseURL, dqlEndpoint)

	options.QueryParams = map[string]string{}
	var responseData *http.Response
	var err error
	options.Headers = map[string]string{
		"Accept":                 "application/json",
		"Content-Type":           "application/json",
		"x-authorization-secret": svc.vpcSecret,
	}
	if isNft {
		options.QueryParams["entityType"] = "nfts"

	}
	responseData, err = svc.httpClient.Driver.Get(options)

	if err != nil {
		logger.Error("[GET_ENTITY_BY_ID] Failed to fetch data", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.serviceName, "GET_ENTITY_BY_ID", "[GET_ENTITY_BY_ID] Failed to fetch data")
		return dqlResponse, err
	}
	if responseData.StatusCode == 503 {
		logger.Error("[GET_ENTITY_BY_ID] Failed to fetch data", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "response": responseData})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, constants.ERROR_TYPES[constants.PROCESS_ERROR].HttpStatus, svc.serviceName, "GET_ENTITY_BY_ID", "[GET_ENTITY_BY_ID] Downstream service unavailable")
		return dqlResponse, err
	}

	data, err := io.ReadAll(responseData.Body)

	if err != nil {
		logger.Error("[GET_ENTITY_BY_ID] Failed to read the response", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err.Error()})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.serviceName, "GET_ENTITY_BY_ID", "[GET_ENTITY_BY_ID] Failed to read the response")
		return dqlResponse, err
	}
	if err := json.Unmarshal(data, &dqlResponse); err != nil {
		logger.Error("[GET_ENTITY_BY_ID] Failed to unmarshal the response", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err.Error()})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.serviceName, "GET_ENTITY_BY_ID", "[GET_ENTITY_BY_ID] Failed to unmarshal the response")
		return dqlResponse, err
	}
	if responseData.StatusCode >= 400 {
		logger.Error("[GET_ENTITY_BY_ID] Request Failed", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode})
		if responseData.StatusCode == 500 {
			err = downstreamerror.New(constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, constants.ERROR_TYPES[constants.PROCESS_ERROR].HttpStatus, svc.serviceName, "GET_ENTITY_BY_ID", "[GET_ENTITY_BY_ID] Request Failed")
		} else {
			err = downstreamerror.New(dqlResponse.Error.ErrorCode, responseData.StatusCode, svc.serviceName, "GET_ENTITY_BY_ID", dqlResponse.Error.Message)
		}
		return dqlResponse, err
	}
	logger.Debug("DQLResponse", map[string]interface{}{"context": ctx, "response": dqlResponse})
	return dqlResponse, nil
}

func (svc *DQLSvc) GetTokenByAddress(ctx context.Context, address string, networkId string) (types.DQLByAddressResponseInterface, error) {
	var dqlResponse types.DQLByAddressResponseInterface
	var options httpTypes.RequestOptions
	options.RetryAttempt = constants.DEFAULT_RETRY_ATTEMPT
	options.RetryFixedInternal = constants.DEFAULT_RETRY_FIXED_INTERVAL

	options.Url = fmt.Sprintf("%s%s", svc.dqlBaseURL, constants.ENDPOINTS["DQL"]["GetEntity"])

	filters := []map[string]interface{}{
		{
			"attribute": "address",
			"operation": "eq",
			"values":    []string{address},
		},
		{
			"attribute": "network_id",
			"operation": "eq",
			"values":    []string{networkId},
		},
	}
	filtersJson, _ := json.Marshal(filters)
	filtersEncoded := url.QueryEscape(string(filtersJson))

	options.QueryParams = map[string]string{
		"entityType": "tokens",
		"filters":    filtersEncoded,
	}
	var responseData *http.Response
	var err error
	options.Headers = map[string]string{
		"Accept":                 "application/json",
		"Content-Type":           "application/json",
		"x-authorization-secret": svc.vpcSecret,
	}
	responseData, err = svc.httpClient.Driver.Get(options)

	if err != nil {
		logger.Error("[GET_ENTITY_BY_ADDRESS] Failed to fetch data", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err.Error()})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.serviceName, "GET_ENTITY_BY_ADDRESS", "[GET_ENTITY_BY_ADDRESS] Failed to fetch data")
		return dqlResponse, err
	}
	if responseData.StatusCode == 503 {
		logger.Error("[GET_ENTITY_BY_ADDRESS] Failed to fetch data", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "response": responseData})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, constants.ERROR_TYPES[constants.PROCESS_ERROR].HttpStatus, svc.serviceName, "GET_ENTITY_BY_ADDRESS", "[GET_ENTITY_BY_ADDRESS] Downstream service unavailable")
		return dqlResponse, err
	}

	data, err := io.ReadAll(responseData.Body)

	if err != nil {
		logger.Error("[GET_ENTITY_BY_ADDRESS] Failed to read the response", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err.Error(), "response": data})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.serviceName, "GET_ENTITY_BY_ADDRESS", "[GET_ENTITY_BY_ADDRESS] Failed to read the response")
		return dqlResponse, err
	}
	if err := json.Unmarshal(data, &dqlResponse); err != nil {
		logger.Error("[GET_ENTITY_BY_ADDRESS] Failed to unmarshal the response", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err.Error()})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.serviceName, "GET_ENTITY_BY_ADDRESS", "[GET_ENTITY_BY_ADDRESS] Failed to unmarshal the response")
		return dqlResponse, err
	}
	if responseData.StatusCode >= 400 {
		logger.Error("[GET_ENTITY_BY_ADDRESS] Request Failed", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "response": dqlResponse})
		if responseData.StatusCode == 500 {
			err = downstreamerror.New(constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, constants.ERROR_TYPES[constants.PROCESS_ERROR].HttpStatus, svc.serviceName, "GET_ENTITY_BY_ADDRESS", "[GET_ENTITY_BY_ADDRESS] Request Failed")
		} else {
			err = downstreamerror.New(dqlResponse.Error.ErrorCode, responseData.StatusCode, svc.serviceName, "GET_ENTITY_BY_ADDRESS", dqlResponse.Error.Message)
		}
		return dqlResponse, err
	}
	logger.Debug("DQLResponseForAddress", map[string]interface{}{"context": ctx, "response": dqlResponse})
	return dqlResponse, nil
}

func (svc *DQLSvc) GetNftByCollectionAndTokenId(ctx context.Context, collectionAddress string, nftTokenId string, networkId string) (types.DQLByAddressResponseInterface, error) {
	var dqlResponse types.DQLByAddressResponseInterface
	var options httpTypes.RequestOptions
	options.RetryAttempt = constants.DEFAULT_RETRY_ATTEMPT
	options.RetryFixedInternal = constants.DEFAULT_RETRY_FIXED_INTERVAL

	options.Url = fmt.Sprintf("%s%s", svc.dqlBaseURL, constants.ENDPOINTS["DQL"]["GetEntity"])

	//TODO: Remove this after Testing
	logger.Error("DQL :: GetNftByCollectionAndTokenId", map[string]interface{}{"URL": options.Url, "collectionAddress": collectionAddress, "nftTokenId": nftTokenId, "networkId": networkId})

	filters := []map[string]interface{}{
		{
			"attribute": "address",
			"operation": "eq",
			"values":    []string{collectionAddress},
		},
		{
			"attribute": "nft_token_id",
			"operation": "eq",
			"values":    []string{nftTokenId},
		},
		{
			"attribute": "network_id",
			"operation": "eq",
			"values":    []string{networkId},
		},
	}
	filtersJson, _ := json.Marshal(filters)
	filtersEncoded := url.QueryEscape(string(filtersJson))

	options.QueryParams = map[string]string{
		"entityType": "nfts",
		"filters":    filtersEncoded,
	}
	var responseData *http.Response
	var err error
	options.Headers = map[string]string{
		"Accept":                 "application/json",
		"Content-Type":           "application/json",
		"x-authorization-secret": svc.vpcSecret,
	}
	responseData, err = svc.httpClient.Driver.Get(options)

	if err != nil {
		logger.Error("[GET_NFT_ENTITY_BY_ADDRESS] Failed to fetch data", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err.Error()})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.serviceName, "GET_NFT_ENTITY_BY_ADDRESS", "[GET_NFT_ENTITY_BY_ADDRESS] Failed to fetch data")
		return dqlResponse, err
	}
	if responseData.StatusCode == 503 {
		logger.Error("[GET_NFT_ENTITY_BY_ADDRESS] Failed to fetch data", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "response": responseData})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, constants.ERROR_TYPES[constants.PROCESS_ERROR].HttpStatus, svc.serviceName, "GET_NFT_ENTITY_BY_ADDRESS", "[GET_NFT_ENTITY_BY_ADDRESS] Downstream service unavailable")
		return dqlResponse, err
	}

	data, err := io.ReadAll(responseData.Body)

	if err != nil {
		logger.Error("[GET_NFT_ENTITY_BY_ADDRESS] Failed to read the response", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err.Error(), "response": data})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.serviceName, "GET_NFT_ENTITY_BY_ADDRESS", "[GET_NFT_ENTITY_BY_ADDRESS] Failed to read the response")
		return dqlResponse, err
	}
	if err := json.Unmarshal(data, &dqlResponse); err != nil {
		logger.Error("[GET_NFT_ENTITY_BY_ADDRESS] Failed to unmarshal the response", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err.Error()})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.serviceName, "GET_NFT_ENTITY_BY_ADDRESS", "[GET_NFT_ENTITY_BY_ADDRESS] Failed to unmarshal the response")
		return dqlResponse, err
	}
	if responseData.StatusCode >= 400 {
		logger.Error("[GET_NFT_ENTITY_BY_ADDRESS] Request Failed", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "response": dqlResponse})
		if responseData.StatusCode == 500 {
			err = downstreamerror.New(constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, constants.ERROR_TYPES[constants.PROCESS_ERROR].HttpStatus, svc.serviceName, "GET_NFT_ENTITY_BY_ADDRESS", "[GET_NFT_ENTITY_BY_ADDRESS] Request Failed")
		} else {
			err = downstreamerror.New(dqlResponse.Error.ErrorCode, responseData.StatusCode, svc.serviceName, "GET_NFT_ENTITY_BY_ADDRESS", dqlResponse.Error.Message)
		}
		return dqlResponse, err
	}
	logger.Debug("DQLResponseForAddress", map[string]interface{}{"context": ctx, "response": dqlResponse})
	return dqlResponse, nil
}

func (svc *DQLSvc) GetNftCollectionById(ctx context.Context, entityId string) (types.DQLNftCollectionResponseInterface, error) {
	var dqlResponse types.DQLNftCollectionResponseInterface
	var options httpTypes.RequestOptions
	options.RetryAttempt = constants.DEFAULT_RETRY_ATTEMPT
	options.RetryFixedInternal = constants.DEFAULT_RETRY_FIXED_INTERVAL

	dqlEndpoint := fmt.Sprintf(constants.ENDPOINTS["DQL"]["GetEntityById"], entityId)
	options.Url = fmt.Sprintf("%s%s", svc.dqlBaseURL, dqlEndpoint)

	options.QueryParams = map[string]string{}
	var responseData *http.Response
	var err error
	options.Headers = map[string]string{
		"Accept":                 "application/json",
		"Content-Type":           "application/json",
		"x-authorization-secret": svc.vpcSecret,
	}

	options.QueryParams["entityType"] = "nft_collections"

	responseData, err = svc.httpClient.Driver.Get(options)

	if err != nil {
		logger.Error("[GET_COLLECTION_ENTITY_BY_ID] Failed to fetch data", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.serviceName, "GET_COLLECTION_ENTITY_BY_ID", "[GET_COLLECTION_ENTITY_BY_ID] Failed to fetch data")
		return dqlResponse, err
	}
	if responseData.StatusCode == 503 {
		logger.Error("[GET_COLLECTION_ENTITY_BY_ID] Failed to fetch data", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "response": responseData})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, constants.ERROR_TYPES[constants.PROCESS_ERROR].HttpStatus, svc.serviceName, "GET_COLLECTION_ENTITY_BY_ID", "[GET_COLLECTION_ENTITY_BY_ID] Downstream service unavailable")
		return dqlResponse, err
	}

	data, err := io.ReadAll(responseData.Body)

	if err != nil {
		logger.Error("[GET_COLLECTION_ENTITY_BY_ID] Failed to read the response", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err.Error()})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.serviceName, "GET_COLLECTION_ENTITY_BY_ID", "[GET_COLLECTION_ENTITY_BY_ID] Failed to read the response")
		return dqlResponse, err
	}
	if err := json.Unmarshal(data, &dqlResponse); err != nil {
		logger.Error("[GET_COLLECTION_ENTITY_BY_ID] Failed to unmarshal the response", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err.Error()})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.serviceName, "GET_COLLECTION_ENTITY_BY_ID", "[GET_COLLECTION_ENTITY_BY_ID] Failed to unmarshal the response")
		return dqlResponse, err
	}
	if responseData.StatusCode >= 400 {
		logger.Error("[GET_COLLECTION_ENTITY_BY_ID] Request Failed", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode})
		if responseData.StatusCode == 500 {
			err = downstreamerror.New(constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, constants.ERROR_TYPES[constants.PROCESS_ERROR].HttpStatus, svc.serviceName, "GET_COLLECTION_ENTITY_BY_ID", "[GET_COLLECTION_ENTITY_BY_ID] Request Failed")
		} else {
			err = downstreamerror.New(dqlResponse.Error.ErrorCode, responseData.StatusCode, svc.serviceName, "GET_COLLECTION_ENTITY_BY_ID", dqlResponse.Error.Message)
		}
		return dqlResponse, err
	}
	logger.Debug("DQLResponse For Collection", map[string]interface{}{"context": ctx, "response": dqlResponse})
	return dqlResponse, nil
}
