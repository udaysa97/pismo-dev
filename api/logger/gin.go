package logger

import (
	"net/http"
	"pismo-dev/api/types"
	"pismo-dev/constants"
	serviceerror "pismo-dev/error/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GinLogErrorAndAbort populates the context with relevant keys
// `returnErr` is the error passed to user in JSON
// if `logmessage` is not supplied, it is set to `returnErr`
// `logmessage` is the actual log supplied to datadog
func GinLogErrorAndAbort(ctx interface{}, code int, errorCode string, message string) {
	stack, fn := GetStackAndFunctionName()

	c, _ := ctx.(*gin.Context)

	c.Set(constants.ERROR_KEY, errorCode)
	c.Set(constants.STACK_KEY, stack)
	c.Set(constants.FUNCTION_KEY, fn)
	c.Set(constants.LEVEL_KEY, logrus.ErrorLevel)
	c.Set(constants.LOG_MESSAGE_KEY, message)

	c.AbortWithStatusJSON(code, types.ResponseDTO[any]{
		Status: types.StatusError,
		Error: &types.ErrorResponse{
			Code:      code,
			ErrorCode: errorCode,
			Message:   message,
			TraceID:   c.GetString(constants.TRACE_ID_KEY),
		},
	})
}

func ErrorResponse(c *gin.Context, errorCode string, logMessages ...string) {
	stack, fn := GetStackAndFunctionName()

	var message string
	var logMessage string
	var code int
	if val, ok := constants.ERROR_CODE_TO_STATUS[errorCode]; ok {
		code = val
	} else {
		code = http.StatusBadRequest
	}
	if len(logMessages) > 0 {
		logMessage = logMessages[0]
	}

	if len(logMessage) == 0 {
		logMessage = message
	}

	c.Set(constants.ERROR_KEY, errorCode)
	c.Set(constants.STACK_KEY, stack)
	c.Set(constants.FUNCTION_KEY, fn)
	c.Set(constants.LEVEL_KEY, logrus.ErrorLevel)

	c.Set(constants.LOG_MESSAGE_KEY, logMessage)

	c.AbortWithStatusJSON(code, types.ResponseDTO[any]{
		Status: types.StatusError,
		Error: &types.ErrorResponse{
			Code:      code,
			ErrorCode: errorCode,
			Message:   message,
			TraceID:   c.GetString(constants.TRACE_ID_KEY),
		},
	})
}

func VpcErrorResponse(c *gin.Context, code int, errorCode string, message string) {
	stack, fn := GetStackAndFunctionName()

	c.Set(constants.ERROR_KEY, errorCode)
	c.Set(constants.STACK_KEY, stack)
	c.Set(constants.FUNCTION_KEY, fn)
	c.Set(constants.LEVEL_KEY, logrus.ErrorLevel)

	c.Set(constants.LOG_MESSAGE_KEY, message)

	c.AbortWithStatusJSON(code, types.ResponseDTO[any]{
		Status: types.StatusError,
		Error: &types.ErrorResponse{
			Code:      code,
			ErrorCode: errorCode,
			Message:   message,
			TraceID:   c.GetString(constants.TRACE_ID_KEY),
		},
	})
}

func ServiceErrorResponse(c *gin.Context, errInstance *serviceerror.ServiceError) {
	ErrorResponse(c, errInstance.ErrorCode, errInstance.Message)
}
