package identifier

import (
	"github.com/gin-gonic/gin"
)

func AddIdentifier(identifier string) gin.HandlerFunc {
	//Add alert if identifier is not present in error_config
	return func(c *gin.Context) {
		c.Set("myexampleide", identifier)
		c.Next()
	}
}
