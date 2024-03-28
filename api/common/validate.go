package common

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	validationerror "pismo-dev/error/validation"
)

type Validator interface {
	Validate() error
}

func ReadAndValidateRequestBody(request *http.Request, dest Validator) error {
	if body, err := io.ReadAll(request.Body); err != nil {
		return validationerror.New(err.Error())
	} else if err := json.Unmarshal(body, &dest); err != nil {
		return validationerror.New(err.Error())
	} else if err := dest.Validate(); err != nil {
		var validationErr *validationerror.ValidationError
		if errors.As(err, &validationErr) {
			return err
		} else {
			validationerror.New(err.Error())
		}
	}
	return nil
}

func ReadRequestBody(request *http.Request, dest any) error {
	if body, err := io.ReadAll(request.Body); err != nil {
		return validationerror.New(err.Error())
	} else if err := json.Unmarshal(body, &dest); err != nil {
		return validationerror.New(err.Error())
	}
	return nil
}
