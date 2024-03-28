package types

import (
	regexUtil "pismo-dev/commonpkg/regex"
	"pismo-dev/constants"
	validationerror "pismo-dev/error/validation"
	"pismo-dev/internal/types"
)

type ExecuteRequest struct {
	OrderType              string                     `json:"operation_type"`
	NftID                  string                     `json:"nft_id"`
	ErcType                string                     `json:"erc_type"`
	Quantity               string                     `json:"amount"`
	RecipientWalletAddress string                     `json:"recipient_wallet_address"`
	NetworkId              string                     `json:"network_id"`
	IsGsnRequired          bool                       `json:"is_gsn_required"`
	GsnIncludeToken        string                     `json:"gsn_include_token"`
	GsnIncludeNetworkId    string                     `json:"gsn_include_network_id,omitempty"`
	GsnIncludeMaxAmount    string                     `json:"gsn_include_max_amount,omitempty"`
	CurrentUser            types.UserDetailsInterface `json:"user_details"`
}

func (req *ExecuteRequest) Validate() error {
	if len(req.OrderType) <= 0 {
		return validationerror.New("operation_type value is invalid")
	}
	if _, ok := constants.ORDER_TYPES[req.OrderType]; !ok {
		return validationerror.New("operation_type value is invalid")
	}
	if len(req.NftID) <= 0 {
		return validationerror.New("`NftId` is a required field")
	}
	if !regexUtil.IsValidUUID(req.NftID) {
		return validationerror.New("Invalid nft_id")
	}
	if (req.CurrentUser == types.UserDetailsInterface{}) {
		return validationerror.New("`user details` is a required field")
	}
	if len(req.CurrentUser.Id) <= 0 {
		return validationerror.New("`userId` is a required field")
	}
	if len(req.CurrentUser.AuthToken) <= 0 {
		return validationerror.New("`AuthToken` is a required field")
	}
	if len(req.CurrentUser.UserOTP) <= 0 {
		return validationerror.New("`OTP` is a required field")
	}
	if len(req.CurrentUser.ReloginPin) <= 0 {
		return validationerror.New("`Relogin PIN` is a required field")
	}
	if len(req.ErcType) <= 0 {
		return validationerror.New("`ercType` is a required field")
	}
	if _, ok := constants.ErcTypes[req.ErcType]; !ok {
		return validationerror.New("Erc Type value is invalid")
	}
	if len(req.Quantity) <= 0 {
		return validationerror.NewCustomError("Invalid Quantity To Transfer", constants.QUANTITY_ERROR)
	}
	if len(req.RecipientWalletAddress) <= 0 {
		return validationerror.New("`ToAddress` is a required field")
	}
	if !regexUtil.IsValidBlockchainAddress(req.RecipientWalletAddress) {
		return validationerror.NewCustomError("Invalid Address", constants.INVALID_ADDRESS_ERROR)
	}
	if len(req.NetworkId) <= 0 {
		return validationerror.New("`NetworkId` is a required field")
	}
	if !regexUtil.IsValidUUID(req.NetworkId) {
		return validationerror.New("Invalid network_id")
	}
	if req.IsGsnRequired {
		// if len(req.GsnIncludeToken) <= 0 {
		// 	return validationerror.New("`GsnIncludeToken` is a required field for GSNRequest")
		// }
		if len(req.GsnIncludeNetworkId) <= 0 {
			return validationerror.New("`GsnIncludeNetworkId` is a required field for GSNRequest")
		}
		if len(req.GsnIncludeMaxAmount) <= 0 {
			return validationerror.New("`GsnIncludeMaxAmount` is a required field for GSNRequest")
		}

	}

	return nil
}
