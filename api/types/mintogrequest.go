package types

import (
	regexUtil "pismo-dev/commonpkg/regex"
	"pismo-dev/constants"
	validationerror "pismo-dev/error/validation"
)

type MintOGRequest struct {
	ContractId    string `json:"contract_id"`
	UserId        string `json:"user_id"`
	OperationType string `json:"operation_type"`
	NetworkId     string `json:"network_id"`
	CustomDataUri string `json:"data"`
}

func (req *MintOGRequest) Validate() error {

	if len(req.ContractId) <= 0 {
		return validationerror.New("Please send contract_id")
	}
	if !regexUtil.IsValidUUID(req.ContractId) {
		return validationerror.New("Invalid contract_id")
	}
	if len(req.UserId) <= 0 {
		return validationerror.New("Please send a user_id")
	}
	if !regexUtil.IsValidUUID(req.UserId) {
		return validationerror.New("Invalid user_id")
	}
	if len(req.OperationType) <= 0 {
		return validationerror.New("Please Send operation_type")
	}
	if _, ok := constants.ORDER_TYPES[req.OperationType]; !ok {
		return validationerror.New("Invalid operation_type")
	}
	if req.OperationType == constants.SS_MINT && len(req.CustomDataUri) == 0 {
		return validationerror.New("Please send data to be sent in nft")
	}
	if len(req.NetworkId) <= 0 {
		return validationerror.New("Please send a network_id")
	}
	if !regexUtil.IsValidUUID(req.NetworkId) {
		return validationerror.New("Invalid network_id")
	}
	return nil

}
