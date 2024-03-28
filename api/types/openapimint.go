package types

import (
	regexUtil "pismo-dev/commonpkg/regex"
	"pismo-dev/constants"
	validationerror "pismo-dev/error/validation"
	extTypes "pismo-dev/internal/types"
)

type OpenAPINFTMintOrder struct {
	OperationType     string     `json:"operation_type"`
	NetworkName       string     `json:"network_name"`
	NetworkId         string     `json:"network_id"`
	CollectionAddress string     `json:"collection_address"`
	CollectionId      string     `json:"collection_id"`
	CollectionName    string     `json:"collection_name"`
	SenderAddress     string     `json:"sender_address"`
	Quantity          string     `json:"quantity"`
	RecipientAddress  string     `json:"recipient_address"`
	UserId            string     `json:"userId"`
	ErcType           string     `json:"ercType"`
	VendorId          string     `json:"vendorId"`
	MetaData          Mintmetada `json:"metadata"`
	IsSponsored       bool       `json:"isSponsored"`
}

type Mintmetada struct {
	CollectionName string                      `json:"collectionName"`
	Uri            string                      `json:"uri"`
	NftName        string                      `json:"nftName"`
	Description    string                      `json:"description"`
	Properties     *[]extTypes.ExtraProperties `json:"properties,omitempty"`
}

func (req *OpenAPINFTMintOrder) Validate() error {
	if len(req.OperationType) <= 0 {
		return validationerror.New("operation_type value is invalid")
	}
	if _, ok := constants.ORDER_TYPES[req.OperationType]; !ok {
		return validationerror.New("operation_type value is invalid")
	}
	if len(req.VendorId) <= 0 {
		return validationerror.New("`Vendor` is a required field")
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
	if len(req.RecipientAddress) <= 0 {
		return validationerror.New("`ToAddress` is a required field")
	}
	if len(req.SenderAddress) <= 0 {
		return validationerror.New("`SenderAddress` is a required field")
	}
	//TODO based address validation based on network type
	if !regexUtil.IsValidBlockchainAddress(req.RecipientAddress) && !regexUtil.IsValidAptosBlockchainAddress(req.RecipientAddress) {
		return validationerror.NewCustomError("Invalid Address", constants.INVALID_ADDRESS_ERROR)
	}
	if len(req.NetworkId) <= 0 {
		return validationerror.New("`NetworkId` is a required field")
	}

	return nil
}
