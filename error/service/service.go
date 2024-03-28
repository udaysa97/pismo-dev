package service

import (
	"pismo-dev/pkg/logger"
)

type ServiceError struct {
	Message   string
	ErrorCode string
	HttpCode  int
}

func (d *ServiceError) Error() string {
	return d.ErrorCode
}

func New(message string, errorCode string, code ...int) *ServiceError {
	httpCode := 400
	if len(code) > 0 {
		httpCode = code[0]
	}
	svcError := &ServiceError{message, errorCode, httpCode}
	logger.Error("Service Error", svcError)
	return svcError
}
