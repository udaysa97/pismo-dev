package types

import (
	regexUtil "pismo-dev/commonpkg/regex"
	"pismo-dev/constants"
	validationerror "pismo-dev/error/validation"
	"pismo-dev/internal/types"
)

type EstimateRequest struct {
	OrderType              string                     `json:"operation_type"`
	NftID                  string                     `json:"nft_id"`
	ErcType                string                     `json:"erc_type"`
	Quantity               string                     `json:"amount"`
	RecipientWalletAddress string                     `json:"recipient_wallet_address"`
	NetworkId              string                     `json:"network_id"`
	CurrentUser            types.UserDetailsInterface `json:"user_details"`
	OrderId                string                     `json:"order_id,emitonempty"`
}

func (req *EstimateRequest) Validate() error {
	if len(req.OrderType) <= 0 || req.OrderType != constants.OPERATION_TYPE {
		return validationerror.New("operation_type value is invalid")
	}
	if len(req.NftID) <= 0 {
		return validationerror.New("`NftId` is a required field")
	}
	if !regexUtil.IsValidUUID(req.NftID) {
		return validationerror.New("Invalid nft_id")
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
	if (req.CurrentUser == types.UserDetailsInterface{}) {
		return validationerror.New("User details not provided")
	}
	if len(req.CurrentUser.Id) <= 0 {
		return validationerror.New("`userId` is a required field")
	}

	return nil
}
