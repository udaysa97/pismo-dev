package types

import (
	validationerror "pismo-dev/error/validation"
)

type TransactionRequest struct {
	AccountId   int      `json:"account_id"`
	OperationId int      `json:"operation_id"`
	Amount      *float32 `json:"amount"`
}

func (req *TransactionRequest) Validate() error {
	if req.AccountId <= 0 {
		return validationerror.New("account_id value is not sent")
	}
	if req.OperationId <= 0 {
		return validationerror.New("operation_id value is not sent")
	}

	if req.Amount == nil {
		return validationerror.New("amount value is not sent")
	}

	return nil
}
