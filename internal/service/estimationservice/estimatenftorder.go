package estimationservice

import (
	"context"
	"fmt"
	apiLogger "pismo-dev/api/logger"
	"pismo-dev/commonpkg/gascalculation"
	"pismo-dev/constants"
	extTypes "pismo-dev/external/types"
	"pismo-dev/internal/types"
	kafkaclient "pismo-dev/pkg/kafka/client"
	"pismo-dev/pkg/logger"
	"strconv"
	"strings"
)

type EstimationSvc struct {
	ServiceName         string
	KafkaProducerClient *kafkaclient.ProducerClient
	services            RequiredServices
	repos               RequiredRepos
}

func NewEstimationSvc(kafkaProducerClient *kafkaclient.ProducerClient) *EstimationSvc {
	return &EstimationSvc{
		ServiceName:         "NFTEstimation",
		KafkaProducerClient: kafkaProducerClient,
	}
}

func (svc *EstimationSvc) SetRequiredServices(services RequiredServices) {
	svc.services = services
}

func (svc *EstimationSvc) SetRequiredRepos(repos RequiredRepos) {
	svc.repos = repos
}

func (svc *EstimationSvc) CheckUserEligible(ctx context.Context, userDetails types.UserDetailsInterface, nftTokenDetails types.NFTTransferDetailsInterface, floatAmount float64) (bool, error) {
	count := svc.repos.OrderMetaDataRepo.GetPendingNftOrdersForUser(userDetails.Id, nftTokenDetails.NftId)
	if count < 0 {
		logger.Error("Error fetching pending orders", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "userId": userDetails.Id})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "DB Error")
		return false, fmt.Errorf("Pending Order Error")
	}
	if count > 0 {
		logger.Error("Pending User Order Detected", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.PENDING_ORDER_ERROR].ErrorCode, "userId": userDetails.Id})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.PENDING_ORDER_ERROR].HttpStatus, constants.ERROR_TYPES[constants.PENDING_ORDER_ERROR].ErrorCode, "Order Already in Progress")
		return false, fmt.Errorf("Pending Order Error")
	}
	tokenBalance, err := svc.services.PortfolioSvc.GetUserBalance(ctx, nftTokenDetails.NftId, nftTokenDetails.NetworkId, userDetails.Id)
	if err != nil {
		logger.Error("Unable to get user balance", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "userId": userDetails.Id})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "Insufficient User balance")
		return false, err
	}
	floatBalance, err := strconv.ParseFloat(tokenBalance.Result.Rows[0].Quantity, 64)
	if err != nil {
		logger.Error("Unable to convert user balance to float", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "userId": userDetails.Id, "error": err.Error()})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, err.Error())
		return false, err
	}
	if floatAmount <= 0 {
		logger.Error("User entered amount less than or equal to 0", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.QUANTITY_ERROR].ErrorCode, "userId": userDetails.Id})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.QUANTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.QUANTITY_ERROR].ErrorCode, "Amount cannot be zero")
		return false, fmt.Errorf("0 Amount error")
	}
	if tokenBalance.Result.Rows[0].EntityId == nftTokenDetails.NftId && floatBalance >= floatAmount {
		return true, nil //TODO: Check locked balance logic with portfolio svc
	} else {
		logger.Error("User Not eligible for transfer", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.QUANTITY_ERROR].ErrorCode, "portfolioResponse": tokenBalance, "jobId": nftTokenDetails.OrderId})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.QUANTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.QUANTITY_ERROR].ErrorCode, "User Not eligible for transfer")
		return false, fmt.Errorf("user Not eligible for transfer")
	}
}

func (svc *EstimationSvc) CalculateEstimate(ctx context.Context, userDetails types.UserDetailsInterface, nftTokenDetails types.NFTTransferDetailsInterface) (types.FMEstimateResultOutput, error) {
	floatAmount, err := strconv.ParseFloat(nftTokenDetails.Amount, 64)
	if err != nil {
		logger.Error("EstimateFlow:Unable to convert requested amount to float", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.QUANTITY_ERROR].ErrorCode, "userId": userDetails.Id, "error": err.Error()})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.QUANTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.QUANTITY_ERROR].ErrorCode, "Invalid Quantity To Transfer")
		return types.FMEstimateResultOutput{}, err
	}
	integerAmount := int(floatAmount)
	if floatAmount != float64(int(floatAmount)) {
		logger.Error("EstimateFlow:User entered amount in decimals", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.QUANTITY_ERROR].ErrorCode, "userId": userDetails.Id})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.QUANTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.QUANTITY_ERROR].ErrorCode, "Invalid Quantity To Transfer")
		return types.FMEstimateResultOutput{}, fmt.Errorf("decimal Amount error")
	}
	_, err = svc.CheckUserEligible(ctx, userDetails, nftTokenDetails, floatAmount)
	if err != nil {
		return types.FMEstimateResultOutput{}, err
	}

	// Can make parallel calls for DQL and Signing
	var signingSvcResponse extTypes.SigningSvcResponse
	if signingSvcResponse, err = svc.services.SigningSvc.GetUserWalletAddress(ctx, userDetails.Id, nftTokenDetails.NetworkId); err != nil {
		logger.Error("EstimationFlow: Error Fetching user address", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "request": nftTokenDetails, "jobId": nftTokenDetails.OrderId, "error": err.Error()})
		if strings.Contains(signingSvcResponse.Error.Message, constants.WALLET_NOT_BACKEDUP_SS_ERROR) {
			apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.WALLET_NOT_BACKED_UP_ERROR].HttpStatus, constants.ERROR_TYPES[constants.WALLET_NOT_BACKED_UP_ERROR].ErrorCode, "Wallet not backed up by user")
		} else {
			apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "Cannot Fetch User Wallet Address")
		}
		return types.FMEstimateResultOutput{}, err
	} else {
		userDetails.UserWalletAddress = signingSvcResponse.Address
	}
	if userDetails.UserWalletAddress == nftTokenDetails.DestinationWalletAddress {
		logger.Error("EstimationFlow: Self Transfer detected", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.SELF_TRANSFER_ERROR].ErrorCode, "request": nftTokenDetails, "jobId": nftTokenDetails.OrderId, "userId": userDetails.Id})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.SELF_TRANSFER_ERROR].HttpStatus, constants.ERROR_TYPES[constants.SELF_TRANSFER_ERROR].ErrorCode, "Cannot Transfer to self")
		return types.FMEstimateResultOutput{}, fmt.Errorf("Self Transfer detected")
	}

	nftDetails, err := svc.services.DQLSvc.GetEntityById(ctx, nftTokenDetails.NftId, true)
	if err != nil {
		logger.Error("EstimationFlow: Error while querying data from DQL", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "DQLResponse": nftDetails, "request": nftTokenDetails, "jobId": nftTokenDetails.OrderId, "error": err.Error()})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "Nft Details not found")
		return types.FMEstimateResultOutput{}, err
	}
	if nftDetails.Details.NetworkId != nftTokenDetails.NetworkId {
		logger.Error("EstimationFlow: Network ID mismatch", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "DQLResponse": nftDetails, "request": nftTokenDetails, "jobId": nftTokenDetails.OrderId})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "Network ID mismatch")
		return types.FMEstimateResultOutput{}, fmt.Errorf("network ID mismatch")
	}
	if nftDetails.Details.ErcType != nftTokenDetails.ErcType {
		logger.Error("EstimationFlow: ErcType mismatch", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "DQLResponse": nftDetails, "request": nftTokenDetails, "jobId": nftTokenDetails.OrderId})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "ErcType mismatch")
		return types.FMEstimateResultOutput{}, fmt.Errorf("ErcType mismatch")
	}

	var payloadForFM types.FMEstimateRequest
	payloadForFM.UserId = userDetails.Id
	payloadForFM.FlowType = constants.FM_FLOW_TYPE
	payloadForFM.Operation = constants.FM_ESTIMATE_OPERATION
	payloadForFM.Payload = types.EstimatePayload{
		Amount:             strconv.Itoa(integerAmount),
		NftContractAddress: nftDetails.Details.Address,
		NftId:              nftDetails.Details.TokenId,
		NftType:            nftDetails.Details.ErcType,
		NetworkId:          nftDetails.Details.NetworkId,
		SenderAddress:      userDetails.UserWalletAddress,
		RecepientAddress:   nftTokenDetails.DestinationWalletAddress,
	}
	logger.Info("Payload for Estimate", map[string]interface{}{"context": ctx, "payload": payloadForFM, "jobId": nftTokenDetails.OrderId})

	responseFromFM, err := svc.services.FlowManagerSvc.GetEstimate(ctx, payloadForFM)
	if err != nil {
		errMsg := "Error estimating transfer for NFT"
		fmErrorMsgSet := false
		if (responseFromFM.Error != types.Error{} && responseFromFM.Error.ErrorCode != "") {
			if msg, ok := constants.FE_ERROR_CODES_MAPPING[responseFromFM.Error.ErrorCode]; ok {
				errMsg = fmt.Sprintf("%s:%s", errMsg, msg.Message)
				fmErrorMsgSet = true
			}
		}
		if (!fmErrorMsgSet && responseFromFM.Error != types.Error{} && responseFromFM.Error.Message != "") {
			errMsg = fmt.Sprintf("%s:%s", errMsg, responseFromFM.Error.Message)
		}
		logger.Error("EstimationFlow: Call to FM failed", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "FMResponse": responseFromFM, "request": nftDetails, "jobId": nftTokenDetails.OrderId})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, errMsg)
		return types.FMEstimateResultOutput{}, err
	}

	logger.Info(`OrderEstimationService response from flow manager:`, map[string]interface{}{"responseFromFM": responseFromFM, "jobId": nftTokenDetails.OrderId})

	finalResponse, err := svc.TransformResponse(ctx, responseFromFM.Output, nftTokenDetails.NetworkId)
	if err != nil {
		logger.Error("OrderEstimationService: error while transforming object", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "error": err, "FMResponse": responseFromFM, "jobId": nftTokenDetails.OrderId})
		return types.FMEstimateResultOutput{}, err
	}
	logger.Info(`OrderEstimationService response from NFTMS:`, map[string]interface{}{"responseFromFM": responseFromFM, "jobId": nftTokenDetails.OrderId})
	return finalResponse, nil
}

func (svc *EstimationSvc) TransformResponse(ctx context.Context, details types.FMEstimateResultOutput, networkId string) (types.FMEstimateResultOutput, error) {
	transactionFees := details.TransactionFee
	networkTokenDecimals := make(map[string]int)
	networkTokenId := make(map[string]string)
	for i, j := range transactionFees {
		tokenDecimals, ok := networkTokenDecimals[j.NetworkId]
		tokenId, tok := networkTokenId[j.NetworkId]
		if !ok || !tok {
			networkDetails, err := svc.services.DQLSvc.GetEntityById(ctx, j.NetworkId, false)
			if err != nil {
				//TODO: Add error and another check for NativeTokenId present in response or not since we are using same stuct across
				logger.Error("OrderEstimationService: error while transforming object", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "error": err, "networkId": j.NetworkId})
				apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "Could not fetch network details")
				return types.FMEstimateResultOutput{}, err
			}
			if (networkDetails.Details == types.DQLDetails{}) {
				logger.Error("OrderEstimationService | Empty network Details", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "network": j.NetworkId})
				apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "Downstream Error: GSN Network info not found")
				return types.FMEstimateResultOutput{}, fmt.Errorf("OrderEstimationService | Network Not found")
			}

			tokenDetails, err := svc.services.DQLSvc.GetEntityById(ctx, networkDetails.Details.NativeCurrencyId, false)
			if err != nil {
				//TODO: Add error and another check for decimals present in response or not since we are using same stuct across
				logger.Error("OrderEstimationService: error while transforming object", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "error": err, "tokenDetails": networkDetails.Details.NativeCurrencyId})
				apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "Could not fetch native token details")
				return types.FMEstimateResultOutput{}, err
			}

			tokenDecimals, err = strconv.Atoi(tokenDetails.Details.Decimals)
			if err != nil {
				logger.Error("OrderEstimationService | Decimals not sent for token", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "tokenId": networkDetails.Details.NativeCurrencyId, "network": networkId, "error": err.Error()})
				apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, err.Error())
				return types.FMEstimateResultOutput{}, err

			}
			tokenId = tokenDetails.Details.Id
			networkTokenDecimals[j.NetworkId] = tokenDecimals
			networkTokenId[j.NetworkId] = tokenId
		}

		transactionFees[i].GasAmount = gascalculation.CalculateGasFees(j.GasAmount, tokenDecimals) //TODO: Research library for big decimal: = j.gasAmount/math.pow(10,decimal)
		transactionFees[i].TokenId = tokenId
	}
	if len(details.GsnWithdrawTokens) > 0 {
		for i, withdrawToken := range details.GsnWithdrawTokens {
			token, err := svc.services.DQLSvc.GetTokenByAddress(ctx, withdrawToken.TokenAddress, withdrawToken.NetworkId)
			if err != nil {
				logger.Error("OrderEstimationService | Token Not found", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "address": withdrawToken.TokenAddress, "network": networkId})
				apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "GSN Token not found at DQL")
				return types.FMEstimateResultOutput{}, fmt.Errorf("OrderEstimationService | Token Not found")
			}
			if len(token.Entities) == 0 || (token.Entities[0].Details == types.DQLDetails{}) {
				logger.Error("OrderEstimationService | No Token details from DQL", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "address": withdrawToken.TokenAddress, "network": networkId})
				apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "Downstream Error: GSN Token info not found")
				return types.FMEstimateResultOutput{}, fmt.Errorf("OrderEstimationService | Token Not found")
			}
			tokenDecimals, err := strconv.Atoi(token.Entities[0].Details.Decimals)
			if err != nil {
				logger.Error("OrderEstimationService | Decimals not sent for token", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "address": withdrawToken.TokenAddress, "network": networkId, "error": err.Error()})
				apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "Downstream Error:GSN Calculation Error: Incorrect Token detail")
				return types.FMEstimateResultOutput{}, err

			}
			details.GsnWithdrawTokens[i].TokenId = token.Entities[0].Details.Id
			details.GsnWithdrawTokens[i].TokenAmount = gascalculation.CalculateGasFees(withdrawToken.TokenAmount, tokenDecimals)
		}
	}
	return details, nil

}
