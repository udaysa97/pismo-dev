package validation

import (
	"fmt"
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
	return &ValidationError{fmt.Sprintf("validation: %s", message), "SOME_ERROR_CODE", 401}
}
