package types

import (
	validationerror "pismo-dev/error/validation"
)

type CreateAccountRequest struct {
	DocumentId string `json:"document_number"`
}

func (req *CreateAccountRequest) Validate() error {
	if len(req.DocumentId) <= 0 {
		return validationerror.New("document_id value is not sent")
	}

	return nil
}
