package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetHealth(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "health ok from pismo"})
}
