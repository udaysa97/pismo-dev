package logger

import (
	"fmt"
	"pismo-dev/constants"
	"pismo-dev/pkg/logger"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetStackAndFunctionName(callerframeSkip ...int) (stack string, fn string) {
	defaultSkip := 2
	if len(callerframeSkip) > 0 {
		defaultSkip = callerframeSkip[0]
	}
	pc, file, line, _ := runtime.Caller(defaultSkip)
	frame, _ := runtime.CallersFrames([]uintptr{pc}).Next()
	splitpath := strings.Split(file, constants.PACKAGE_NAME)

	if len(splitpath) == 2 {
		file = splitpath[1]
	}
	//extract function name
	split := strings.Split(frame.Function, "/")
	function := split[len(split)-1]
	return fmt.Sprintf(".%s:%d", file, line), fmt.Sprintf("%s()", function)
}

func defaultLogEntry(c *gin.Context, argument ...interface{}) *logrus.Entry {
	stack, fn := GetStackAndFunctionName(3)

	fields := logrus.Fields{
		"stack":       stack,
		"function":    fn,
		"dd.trace_id": c.GetString(constants.TRACE_ID_KEY),
	}

	var arg interface{}
	if len(argument) > 0 {
		arg = argument[0]
		fields["metadata"] = fmt.Sprintf("%+v", arg)
	}
	return logger.GetDefaultLogger().WithFields(fields)

}

func Info(c *gin.Context, message string, argument ...interface{}) {
	defaultLogEntry(c, argument...).Info(message)
}

func Error(c *gin.Context, message string, argument ...interface{}) {
	defaultLogEntry(c, argument...).Error(message)
}

func Debug(c *gin.Context, message string, argument ...interface{}) {
	defaultLogEntry(c, argument...).Debug(message)
}

func Warn(c *gin.Context, message string, argument ...interface{}) {
	defaultLogEntry(c, argument...).Warn(message)
}
