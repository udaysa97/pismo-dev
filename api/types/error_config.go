package types

import (
	"fmt"
	"time"
)

type ErrorConfigRequest struct {
	Identifier   string `json:"identifier"`
	ErrorCode    string `json:"errorCode"`
	HttpCode     int32  `json:"httpCode"`
	ErrorMessage string `json:"errorMessage"`
}

type ErrorConfigResponse struct {
	Identifier   string    `json:"identifier"`
	ErrorCode    string    `json:"errorCode"`
	HttpCode     int32     `json:"httpCode"`
	ErrorMessage string    `json:"errorMessage"`
	CreatedAt    time.Time `json:"createdAt,omitempty"`
	UpdatedAt    time.Time `json:"updatedAt,omitempty"`
	Message      string    `json:"message,omitempty"`
}

func (req ErrorConfigRequest) Validate() error {
	if len(req.Identifier) == 0 {
		return fmt.Errorf("Invalid Identifier")
	} else if len(req.ErrorCode) == 0 {
		return fmt.Errorf("Invalid Error Code")
	} else if req.HttpCode == 0 {
		return fmt.Errorf("Invalid Http Code")
	} else if len(req.ErrorMessage) == 0 {
		return fmt.Errorf("Invalid Error Message")
	}
	return nil
}
