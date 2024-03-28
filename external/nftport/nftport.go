package nftport

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
)

type NftPortSvc struct {
	ServiceName    string
	NftPortBaseURL string
	HttpClient     *httpclient.HttpClientWrapper
	ApiKey         string
}

func NewNftPortSvc(client *httpclient.HttpClientWrapper, nftPortBaseURL string, apiKey string) *NftPortSvc {
	return &NftPortSvc{

		ServiceName:    "NFTPortSvc",
		HttpClient:     client,
		NftPortBaseURL: nftPortBaseURL,
		ApiKey:         apiKey,
	}
}

func (svc *NftPortSvc) MintNft(ctx context.Context, chain, contract_address, metadata_uri, to_address string) (types.NftPortMintResponse, error) {
	var response types.NftPortMintResponse
	var options httpTypes.RequestOptions
	options.RetryAttempt = 0
	options.RetryFixedInternal = 0

	options.Url = fmt.Sprintf("%s%s", svc.NftPortBaseURL, constants.ENDPOINTS["NFTPORT"]["MINT"])

	var responseData *http.Response
	var err error
	options.Headers = map[string]string{
		"Accept":        "application/json",
		"Content-Type":  "application/json",
		"Authorization": svc.ApiKey,
	}

	payload := types.NftPortMintRequest{
		Chain:           chain,
		ContractAddress: contract_address,
		MetadataURI:     metadata_uri,
		MintToAddress:   to_address,
	}
	options.Body, _ = json.Marshal(payload)

	responseData, err = svc.HttpClient.Driver.Post(options)

	if err != nil {
		logger.Error("[MINT_NFT_USING_PORT] Failed to fetch data", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.ServiceName, "MINT_NFT_USING_PORT", "[MINT_NFT_USING_PORT] Failed to fetch data")
		return response, err
	}
	if responseData.StatusCode >= 500 {
		logger.Error("[MINT_NFT_USING_PORT] Failed to fetch data", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "response": responseData})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, constants.ERROR_TYPES[constants.PROCESS_ERROR].HttpStatus, svc.ServiceName, "MINT_NFT_USING_PORT", "[MINT_NFT_USING_PORT] Downstream service unavailable")
		return response, err
	}

	data, err := io.ReadAll(responseData.Body)

	if err != nil {
		logger.Error("[MINT_NFT_USING_PORT] Failed to read the response", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.ServiceName, "MINT_NFT_USING_PORT", "[MINT_NFT_USING_PORT] Failed to read the response")
		return response, err
	}
	if err := json.Unmarshal(data, &response); err != nil {
		logger.Error("[MINT_NFT_USING_PORT] Failed to unmarshal the response", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err.Error(), "response": data})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.ServiceName, "MINT_NFT_USING_PORT", "[MINT_NFT_USING_PORT] Failed to unmarshal the response")
		return response, err
	}
	if responseData.StatusCode != 200 || response.Response != "OK" {
		logger.Error("[MINT_NFT_USING_PORT] Request Failed", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "response": response})

		err = downstreamerror.New(response.Error.Code, responseData.StatusCode, svc.ServiceName, "MINT_NFT_USING_PORT", response.Error.Message)

		return response, err
	}
	logger.Debug("Mint Response", map[string]interface{}{"context": ctx, "response": response})

	return response, nil
}

func (svc *NftPortSvc) FetchMintStatus(ctx context.Context, chain, txHash string) (types.NftPortStatusResponse, error) {
	var response types.NftPortStatusResponse
	var options httpTypes.RequestOptions
	options.RetryAttempt = 0
	options.RetryFixedInternal = 0

	options.Url = fmt.Sprintf("%s%s", svc.NftPortBaseURL, fmt.Sprintf(constants.ENDPOINTS["NFTPORT"]["STATUS"], txHash, chain))

	var responseData *http.Response
	var err error
	options.Headers = map[string]string{
		"Accept":        "application/json",
		"Content-Type":  "application/json",
		"Authorization": svc.ApiKey,
	}

	responseData, err = svc.HttpClient.Driver.Get(options)

	if err != nil {
		logger.Error("[STATUS_NFT_USING_PORT] Failed to fetch data", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.ServiceName, "STATUS_NFT_USING_PORT", "[STATUS_NFT_USING_PORT] Failed to fetch data")
		return response, err
	}
	if responseData.StatusCode >= 500 {
		logger.Error("[STATUS_NFT_USING_PORT] Failed to fetch data", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "response": responseData})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, constants.ERROR_TYPES[constants.PROCESS_ERROR].HttpStatus, svc.ServiceName, "STATUS_NFT_USING_PORT", "[STATUS_NFT_USING_PORT] Downstream service unavailable")
		return response, err
	}

	data, err := io.ReadAll(responseData.Body)

	if err != nil {
		logger.Error("[STATUS_NFT_USING_PORT] Failed to read the response", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.ServiceName, "STATUS_NFT_USING_PORT", "[STATUS_NFT_USING_PORT] Failed to read the response")
		return response, err
	}
	if err := json.Unmarshal(data, &response); err != nil {
		logger.Error("[STATUS_NFT_USING_PORT] Failed to unmarshal the response", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err.Error(), "response": data})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.ServiceName, "STATUS_NFT_USING_PORT", "[STATUS_NFT_USING_PORT] Failed to unmarshal the response")
		return response, err
	}
	//logger.Debug("Status Response", map[string]interface{}{"context": ctx, "response": response})

	return response, nil

}
