package log

import (
	"bytes"
	"encoding/json"
	"io"
	"pismo-dev/constants"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func LogMiddleware(skipPath ...string) gin.HandlerFunc {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		// PrettyPrint: true, DisableHTMLEscape: true
	})

	return func(c *gin.Context) {
		ctx := c.Request.Context()

		span, _ := tracer.SpanFromContext(ctx)
		trace_id := strconv.FormatUint(span.Context().TraceID(), 10)

		c.Set(constants.TRACE_ID_KEY, trace_id)

		if shouldSkip(c.Request.RequestURI, skipPath...) {
			c.Next()
			return
		}

		requestLogEntry(c, logger)

		start := time.Now()
		c.Next()
		duration := time.Since(start)

		responseLogEntry(c, logger, duration)

	}
}

// all the fields required to be logged regarding a request
func requestLogEntry(c *gin.Context, logger *logrus.Logger) *logrus.Entry {

	requestLogEntry := logger.WithFields(logrus.Fields{
		"method":        c.Request.Method,
		"endpoint":      c.Request.RequestURI,
		"handler_chain": handlerChainString(c),
		"dd.trace_id":   c.GetString(constants.TRACE_ID_KEY),
	})

	switch c.Request.Method {
	case "GET":
		requestLogEntry = requestLogEntry.WithField("query", urlQueryParamMap(c))
	case "POST":
		body, err := jsonRequestBody(c)
		if err != nil {
			requestLogEntry.WithField("request", c.Request).Errorf("[logging] failed to print request body: %s ", err)
		} else {
			if _, exists := body["user_details"]; exists {
				delete(body["user_details"].(map[string]interface{}), "reloginPin")
				delete(body["user_details"].(map[string]interface{}), "authToken")
				delete(body["user_details"].(map[string]interface{}), "userWalletAddress")
				delete(body["user_details"].(map[string]interface{}), "userOTP")

			}
			requestLogEntry = requestLogEntry.WithField("body", body)
		}
	}

	requestLogEntry.Infof("[GIN] %s %s %s", c.Request.Method, c.Request.RequestURI, c.Keys[constants.TRACE_ID_KEY])

	return requestLogEntry
}

// Fields required to be logged regarding a response
func responseLogEntry(c *gin.Context, logger *logrus.Logger, timetaken ...time.Duration) *logrus.Entry {
	responseLogEntry := logger.WithFields(logrus.Fields{
		"code":          c.Writer.Status(),
		"method":        c.Request.Method,
		"endpoint":      c.Request.RequestURI,
		"handler_chain": handlerChainString(c),
		"userId":        c.Request.Header.Get("user-id"),
		"dd.trace_id":   c.GetString(constants.TRACE_ID_KEY),
	})

	if len(timetaken) > 0 {
		responseLogEntry = responseLogEntry.WithFields(logrus.Fields{
			"duration": timetaken[0].String(),
		})
	}

	if c.Writer.Status() != 200 {
		requestbody, _ := c.Get(constants.REQUEST_BODY_KEY)
		responseLogEntry = responseLogEntry.WithFields(logrus.Fields{
			"error":       c.GetString(constants.ERROR_KEY),
			"logMessage":  c.GetString(constants.LOG_MESSAGE_KEY),
			"stack":       c.GetString(constants.STACK_KEY),
			"function":    c.GetString(constants.FUNCTION_KEY),
			"requestBody": requestbody,
		})
		responseLogEntry.Errorf("[%d][GIN] request fail for trace_id %s: %s", c.Writer.Status(), c.Keys["dd.trace_id"], c.GetString(constants.LOG_MESSAGE_KEY))

	} else {
		responseLogEntry.Infof("[%d][GIN] request OK for trace_id: %s", c.Writer.Status(), c.Keys["dd.trace_id"])
	}

	return responseLogEntry
}

func handlerChainString(c *gin.Context) string {
	chain := []string{}

	for _, handler := range c.HandlerNames() {
		handlerPathTokens := strings.Split(handler, "/")
		if len(handlerPathTokens) != 0 {
			functionName := handlerPathTokens[len(handlerPathTokens)-1]
			functionName = strings.TrimSuffix(functionName, ".func1")
			chain = append(chain, functionName)
		}
	}

	return strings.Join(chain, " -> ")
}

func urlQueryParamMap(c *gin.Context) map[string]any {
	querymap := map[string]any{}
	for k, v := range c.Request.URL.Query() {
		if len(v) == 1 {
			querymap[k] = v[0]
		} else if len(v) > 1 {
			querymap[k] = v
		} else if len(v) == 0 {
			querymap[k] = "empty"
		}
	}

	return querymap
}

func jsonRequestBody(c *gin.Context) (map[string]any, error) {
	var kvpair map[string]any

	if byteBody, err := io.ReadAll(c.Request.Body); err != nil {
		return kvpair, err
	} else if err = json.Unmarshal(byteBody, &kvpair); err != nil {
		return kvpair, err
	} else {
		c.Request.Body = io.NopCloser(bytes.NewBuffer(byteBody))
		return kvpair, nil
	}

}

func shouldSkip(path string, skiplist ...string) bool {
	for i := range skiplist {
		if strings.Contains(path, skiplist[i]) {
			return true
		}
	}
	return false
}
