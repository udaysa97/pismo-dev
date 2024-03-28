package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"pismo-dev/constants"
	"pismo-dev/internal/appconfig"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	logger     *logrus.Logger
	loggerOnce sync.Once
	appName    string
)

func init() {
	logger = GetDefaultLogger()
}

func SetAppName(currAppName string) {
	appName = currAppName
}

// creating the single instance of the logger and setting its format
func GetDefaultLogger() *logrus.Logger {
	if logger == nil {
		loggerOnce.Do(func() {
			logger = logrus.New()
			logger.SetFormatter(&logrus.JSONFormatter{
				// PrettyPrint: true, DisableHTMLEscape: true,
			})
			logger.SetLevel(getLogLevel())
			logger.SetOutput(os.Stdout)
			logger.Info("Logger initialized")
		})
	}
	return logger
}

// read the log level from env and return the logrus log level
func getLogLevel() logrus.Level {
	if level, err := logrus.ParseLevel(appconfig.LOG_LEVEL); err == nil {
		return level
	} else {
		return logrus.DebugLevel
	}
}

// GetStackAndFunctionName return the stack and function name
func GetStackAndFunctionName(callerframeSkip ...int) (stack string, fn string) {
	defaultSkip := 2
	if len(callerframeSkip) > 0 {
		defaultSkip = callerframeSkip[0]
	}
	pc, file, line, _ := runtime.Caller(defaultSkip)
	frame, _ := runtime.CallersFrames([]uintptr{pc}).Next()
	splitpath := strings.Split(file, appName)

	if len(splitpath) == 2 {
		file = splitpath[1]
	}
	//extract function name
	split := strings.Split(frame.Function, "/")
	function := split[len(split)-1]
	return fmt.Sprintf(".%s:%d", file, line), fmt.Sprintf("%s()", function)
}

// customizing the log with the kind of data we want to return
func defaultLogEntry(argument ...interface{}) *logrus.Entry {
	stack, function := GetStackAndFunctionName(2)
	fields := logrus.Fields{
		"stack":    stack,
		"function": function,
	}

	if len(argument) > 0 {
		for index, arguments := range argument {
			if val, ok := IsInterfaceMap(arguments); ok {
				for key, value := range val {
					if key == "context" {
						if ginVal, ok := value.(*gin.Context); ok {
							fields[constants.TRACE_ID_KEY] = ginVal.GetString(constants.TRACE_ID_KEY)
							fields[constants.INTERNAL_TRACE_ID_KEY] = ginVal.GetString(constants.TRACE_ID_KEY)
						}
					} else {
						if byteval, ok := isByteSlice(value); ok {
							fields[key] = string(byteval)
							continue
						}
						marshalledVal, _ := json.Marshal(value)
						fields[key] = string(marshalledVal)
					}
				}
			} else {

				key := "argument" + strconv.Itoa(index)
				fields[key] = fmt.Sprintf("%+v", arguments)
			}
		}
	}

	return logger.WithFields(fields)
}

func Info(message string, argument ...interface{}) {
	defaultLogEntry(argument...).Info(message)
}

func IsInterfaceMap(value any) (map[string]interface{}, bool) {
	finalVal, ok := value.(map[string]any)
	return finalVal, ok
}

func isByteSlice(value any) ([]byte, bool) {
	byteSlice, ok := value.([]byte)
	return byteSlice, ok
}

func Error(message string, argument ...interface{}) {
	stack, function := GetStackAndFunctionName(2)
	fields := logrus.Fields{
		"stack":    stack,
		"function": function,
	}

	if len(argument) > 0 {
		for index, arguments := range argument {
			if val, ok := IsInterfaceMap(arguments); ok {
				for key, value := range val {
					if key == "context" {
						if ginVal, ok := value.(*gin.Context); ok {
							fields[constants.TRACE_ID_KEY] = ginVal.GetString(constants.TRACE_ID_KEY)
							fields[constants.INTERNAL_TRACE_ID_KEY] = ginVal.GetString(constants.TRACE_ID_KEY)
						}
					} else {
						if byteval, ok := isByteSlice(value); ok {
							fields[key] = string(byteval)
							continue
						}
						marshalledVal, _ := json.Marshal(value)
						fields[key] = string(marshalledVal)
					}
				}
			} else {

				key := "argument" + strconv.Itoa(index)
				fields[key] = fmt.Sprintf("%+v", arguments)
			}
		}
	}
	logger.WithFields(fields).Error(message)
}

func Debug(message string, argument ...interface{}) {
	defaultLogEntry(argument...).Debug(message)
}

func Warn(message string, argument ...interface{}) {
	defaultLogEntry(argument...).Warn(message)
}

func Log(argument ...interface{}) *logrus.Entry {

	pc, file, line, _ := runtime.Caller(1)
	frame, _ := runtime.CallersFrames([]uintptr{pc}).Next()
	splitpath := strings.Split(file, appName)

	if len(splitpath) == 2 {
		file = splitpath[1]
	}
	//extract function name
	split := strings.Split(frame.Function, "/")
	function := split[len(split)-1]

	fields := logrus.Fields{
		"stack":    fmt.Sprintf(".%s:%d", file, line),
		"function": fmt.Sprintf("%s()", function),
	}
	var arg interface{}
	if len(argument) > 0 {
		arg = argument[0]
		fields["argument"] = fmt.Sprintf("%+v", arg)
	}

	return logger.WithFields(fields)
}
