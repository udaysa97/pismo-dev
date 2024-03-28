package util

import (
	downstreamerror "pismo-dev/error/downstreamError"
	serviceerror "pismo-dev/error/service"

	"github.com/pkg/errors"
	"github.com/sony/gobreaker"
)

func IsCircuitBreakerErr(err error) bool {
	return err == gobreaker.ErrTooManyRequests || err == gobreaker.ErrOpenState
}

func IsServiceError(err error) bool {
	var errInstance *serviceerror.ServiceError
	isServiceError := errors.As(err, &errInstance)
	return isServiceError
}

func GetServiceError(err error) (*serviceerror.ServiceError, bool) {
	var errInstance *serviceerror.ServiceError
	isServiceError := errors.As(err, &errInstance)
	return errInstance, isServiceError
}

func IsEmptyMap(obj any) bool {
	val, ok := obj.(map[string]any)
	return ok && len(val) == 0
}

func GetDownstreamError(err error) (*downstreamerror.DownstreamError, bool) {
	var errInstance *downstreamerror.DownstreamError
	isDownstreamError := errors.As(err, &errInstance)
	return errInstance, isDownstreamError
}
