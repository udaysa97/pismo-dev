package auth

import (
	"os"
	"pismo-dev/api/logger"
	"pismo-dev/constants"

	"github.com/gin-gonic/gin"
)

func VPCAuthorizationMiddleware(skip ...bool) gin.HandlerFunc {
	return func(c *gin.Context) {

		authtoken := c.GetHeader("x-authorization-secret")
		if len(skip) > 0 && skip[0] {
			c.Next()
			return
		} else if authtoken == "" {
			errorMessage := "unauthorized: empty `Authorization` header"
			logger.ErrorResponse(c, constants.ERROR_TYPES[constants.UNAUTHORISED_ERROR].ErrorCode, errorMessage)
			return

		} else if success := isWhitelisted(authtoken); !success {
			errorMessage := "x-authorization-secret invalid"
			logger.ErrorResponse(c, constants.ERROR_TYPES[constants.UNAUTHORISED_ERROR].ErrorCode, errorMessage)
			return
		} else {
			c.Next()
		}

	}
}

func isWhitelisted(token string) bool {
	whitelist := os.Getenv("OKTO_VPC_SECRET")
	return token == whitelist
}
