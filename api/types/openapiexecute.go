package types

import (
	regexUtil "pismo-dev/commonpkg/regex"
	"pismo-dev/constants"
	validationerror "pismo-dev/error/validation"
)

type OpenApiExecuteRequest struct {
	UserId                 string `json:"userId"`
	VendorId               string `json:"vendorId"`
	NetworkId              string `json:"networkId"`
	NFTId                  string `json:"nftId"`
	NFTTokenId             string `json:"nftTokenId"`
	CollectionAddress      string `json:"collectionAddress"`
	CollectionId           string `json:"collectionId"`
	RecipientWalletAddress string `json:"recipientWalletAddress"`
	SenderAddress          string `json:"SenderWalletAddress"`
	ErcType                string `json:"ercType"`
	Quantity               string `json:"quantity"` // in gwei
	OperationType          string `json:"operationType"`
	IsGsnRequired          bool   `json:"isGsnRequired"`
	GsnIncludeToken        string `json:"gsnIncludeToken"` // address
	GsnIncludeNetworkId    string `json:"gsnIncludeNetworkId"`
	GsnIncludeMaxAmount    string `json:"gsnIncludeMaxAmount"`
	IsSponsored            bool   `json:"isSponsored"`
}

func (req *OpenApiExecuteRequest) Validate() error {
	if len(req.OperationType) <= 0 {
		return validationerror.New("operation_type value is invalid")
	}
	if _, ok := constants.ORDER_TYPES[req.OperationType]; !ok {
		return validationerror.New("operation_type value is invalid")
	}
	if len(req.NFTId) <= 0 {
		return validationerror.New("`NftId` is a required field")
	}
	if len(req.VendorId) <= 0 {
		return validationerror.New("`Vendor` is a required field")
	}

	//if !regexUtil.IsValidUUID(req.NFTId) {
	//	return validationerror.New("Invalid nft_id")
	//}
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
	if len(req.SenderAddress) <= 0 {
		return validationerror.New("`SenderAddress` is a required field")
	}
	if !regexUtil.IsValidBlockchainAddress(req.RecipientWalletAddress) && !regexUtil.IsValidAptosBlockchainAddress(req.RecipientWalletAddress) {
		return validationerror.NewCustomError("Invalid Address", constants.INVALID_ADDRESS_ERROR)
	}
	if len(req.NetworkId) <= 0 {
		return validationerror.New("`NetworkId` is a required field")
	}

	return nil
}
