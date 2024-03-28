package validation

import (
	"fmt"
	"pismo-dev/constants"
)

type ValidationError struct {
	Message   string
	ErrorCode string
	HttpCode  int
}

func (d ValidationError) Error() string {
	return d.ErrorCode
}

func New(message string) error {
	return &ValidationError{fmt.Sprintf("validation: %s", message), constants.ERROR_TYPES[constants.BAD_REQUEST_ERROR].ErrorCode, constants.ERROR_TYPES[constants.BAD_REQUEST_ERROR].HttpStatus}
}

func NewCustomError(message string, errorType string) error {
	return &ValidationError{fmt.Sprintf("validation: %s", message), constants.ERROR_TYPES[errorType].ErrorCode, constants.ERROR_TYPES[errorType].HttpStatus}
}
