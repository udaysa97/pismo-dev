package crossmint

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

type CrossMintSvc struct {
	ServiceName      string
	CrossMintBaseURL string
	HttpClient       *httpclient.HttpClientWrapper
	ClientKey        string
	ProjectId        string
}

func NewCrossMintSvc(client *httpclient.HttpClientWrapper, crossMintBaseURL, clientKey, projectId string) *CrossMintSvc {
	return &CrossMintSvc{

		ServiceName:      "CrossMintSvc",
		HttpClient:       client,
		CrossMintBaseURL: crossMintBaseURL,
		ClientKey:        clientKey,
		ProjectId:        projectId,
	}
}

func (svc *CrossMintSvc) MintNft(ctx context.Context, chain, contractIdentifier, metadataUri, toAddress string, isBYOC bool) (types.CrossMintMintResponse, error) {
	var response types.CrossMintMintResponse
	var options httpTypes.RequestOptions
	options.RetryAttempt = 0
	options.RetryFixedInternal = 0

	options.Url = fmt.Sprintf("%s%s", svc.CrossMintBaseURL, fmt.Sprintf(constants.ENDPOINTS["CROSSMINT"]["MINT"], contractIdentifier))

	var responseData *http.Response
	var err error
	options.Headers = map[string]string{
		"Accept":          "application/json",
		"Content-Type":    "application/json",
		"x-client-secret": svc.ClientKey,
		"x-project-id":    svc.ProjectId,
	}
	var payload interface{}
	if isBYOC {
		payload = types.CrossMintBYOCMintRequest{
			Recipient: fmt.Sprintf("%s:%s", chain, toAddress),
			ContractArguments: types.ContractArguments{
				URI: metadataUri,
			},
		}
	} else {
		payload = types.CrossMintMintRequest{
			Recipient: fmt.Sprintf("%s:%s", chain, toAddress),
			MetaData:  metadataUri,
		}
	}
	options.Body, _ = json.Marshal(payload)

	responseData, err = svc.HttpClient.Driver.Post(options)

	if err != nil {
		logger.Error("[MINT_NFT_USING_CROSSMINT] Failed to fetch data", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.ServiceName, "MINT_NFT_USING_CROSSMINT", "[MINT_NFT_USING_CROSSMINT] Failed to fetch data")
		return response, err
	}
	if responseData.StatusCode >= 500 {
		logger.Error("[MINT_NFT_USING_CROSSMINT] Failed to fetch data", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "response": responseData})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, constants.ERROR_TYPES[constants.PROCESS_ERROR].HttpStatus, svc.ServiceName, "MINT_NFT_USING_CROSSMINT", "[MINT_NFT_USING_CROSSMINT] Downstream service unavailable")
		return response, err
	}

	data, err := io.ReadAll(responseData.Body)

	if err != nil {
		logger.Error("[MINT_NFT_USING_CROSSMINT] Failed to read the response", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.ServiceName, "MINT_NFT_USING_CROSSMINT", "[MINT_NFT_USING_CROSSMINT] Failed to read the response")
		return response, err
	}
	if err := json.Unmarshal(data, &response); err != nil {
		logger.Error("[MINT_NFT_USING_CROSSMINT] Failed to unmarshal the response", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err.Error(), "response": data})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.ServiceName, "MINT_NFT_USING_CROSSMINT", "[MINT_NFT_USING_CROSSMINT] Failed to unmarshal the response")
		return response, err
	}
	if responseData.StatusCode != 200 || response.Error {
		logger.Error("[MINT_NFT_USING_CROSSMINT] Request Failed", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "response": response})

		err = downstreamerror.New("401", responseData.StatusCode, svc.ServiceName, "MINT_NFT_USING_CROSSMINT", response.Message)

		return response, err
	}
	//logger.Debug("Mint Response", map[string]interface{}{"context": ctx, "response": response})

	return response, nil
}

func (svc *CrossMintSvc) FetchMintStatus(ctx context.Context, chain, txHash, contractIdentifier string) (types.CrossMintStatusResponse, error) {
	var response types.CrossMintStatusResponse
	var options httpTypes.RequestOptions
	options.RetryAttempt = 0
	options.RetryFixedInternal = 0

	options.Url = fmt.Sprintf("%s%s", svc.CrossMintBaseURL, fmt.Sprintf(constants.ENDPOINTS["CROSSMINT"]["STATUS"], contractIdentifier, txHash))

	var responseData *http.Response
	var err error
	options.Headers = map[string]string{
		"Accept":          "application/json",
		"Content-Type":    "application/json",
		"x-client-secret": svc.ClientKey,
		"x-project-id":    svc.ProjectId,
	}

	responseData, err = svc.HttpClient.Driver.Get(options)

	if err != nil {
		logger.Error("[STATUS_NFT_USING_CROSSMINT] Failed to fetch data", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.ServiceName, "STATUS_NFT_USING_CROSSMINT", "[STATUS_NFT_USING_CROSSMINT] Failed to fetch data")
		return response, err
	}
	if responseData.StatusCode >= 500 {
		logger.Error("[STATUS_NFT_USING_CROSSMINT] Failed to fetch data", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "response": responseData})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, constants.ERROR_TYPES[constants.PROCESS_ERROR].HttpStatus, svc.ServiceName, "STATUS_NFT_USING_CROSSMINT", "[STATUS_NFT_USING_CROSSMINT] Downstream service unavailable")
		return response, err
	}

	data, err := io.ReadAll(responseData.Body)

	if err != nil {
		logger.Error("[STATUS_NFT_USING_CROSSMINT] Failed to read the response", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.ServiceName, "STATUS_NFT_USING_CROSSMINT", "[STATUS_NFT_USING_CROSSMINT] Failed to read the response")
		return response, err
	}
	if err := json.Unmarshal(data, &response); err != nil {
		logger.Error("[STATUS_NFT_USING_CROSSMINT] Failed to unmarshal the response", map[string]interface{}{"context": ctx, "errorCode": constants.ERROR_TYPES[constants.DATA_NOT_FOUND_ERROR].ErrorCode, "error": err.Error(), "response": data})
		err = downstreamerror.New(constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, svc.ServiceName, "STATUS_NFT_USING_CROSSMINT", "[STATUS_NFT_USING_CROSSMINT] Failed to unmarshal the response")
		return response, err
	}
	//logger.Debug("Status Response", map[string]interface{}{"context": ctx, "response": response})

	return response, nil

}
