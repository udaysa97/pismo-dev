package types

import (
	"fmt"
	regexUtil "pismo-dev/commonpkg/regex"
	validationerror "pismo-dev/error/validation"
)

type GetUserCollectionMintCountsRequest struct {
	UserId            string `json:"user_id"`
	CollectionAddress string `josn:"collection_address"`
}

func (req GetUserCollectionMintCountsRequest) Validate() error {
	if len(req.UserId) == 0 {
		return fmt.Errorf("user_id is necessary")
	}
	if !regexUtil.IsValidUUID(req.UserId) {
		return validationerror.New("Invalid user_id")
	}
	if len(req.CollectionAddress) == 0 {
		return fmt.Errorf("collection_address is necessary")
	}
	return nil
}
