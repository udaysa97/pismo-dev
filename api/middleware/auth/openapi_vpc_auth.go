package auth

import (
	"pismo-dev/api/logger"

	"os"
	"pismo-dev/constants"

	"github.com/gin-gonic/gin"
)

func OpenAPIVPCAuthorizationMiddleware(skip ...bool) gin.HandlerFunc {
	return func(c *gin.Context) {

		authtoken := c.GetHeader("x-authorization-secret")
		if len(skip) > 0 && skip[0] {
			c.Next()
			return
		} else if authtoken == "" {
			errorMessage := "unauthorized: empty `Authorization` header"
			logger.ErrorResponse(c, constants.ERROR_TYPES[constants.UNAUTHORISED_ERROR].ErrorCode, errorMessage)
			return

		} else if success := isOpenApiWhitelisted(authtoken); !success {
			errorMessage := "x-authorization-secret invalid"
			logger.ErrorResponse(c, constants.ERROR_TYPES[constants.UNAUTHORISED_ERROR].ErrorCode, errorMessage)
			return
		} else {
			c.Next()
		}

	}
}

func isOpenApiWhitelisted(token string) bool {
	whitelist := os.Getenv("OPENAPI_VPC_SECRET")
	return token == whitelist
}
