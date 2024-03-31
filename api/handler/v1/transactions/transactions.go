package transaction

import (
	"pismo-dev/api/common"
	"pismo-dev/api/types"

	"pismo-dev/internal/repository"
	"pismo-dev/internal/service"

	"github.com/gin-gonic/gin"
)

func InsertTransaction(services *service.Service, repos *repository.Repositories) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		request := types.TransactionRequest{}

		if err := common.ReadAndValidateRequestBody(ctx.Request, &request); err != nil {
			ctx.AbortWithStatus(400)
			return
		}
		result := repos.TransactionRepo.InsertTransaction(request)
		ctx.AbortWithStatusJSON(200, result)
	}
}

func InsertAccount(services *service.Service, repos *repository.Repositories) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		request := types.CreateAccountRequest{}

		if err := common.ReadAndValidateRequestBody(ctx.Request, &request); err != nil {
			ctx.AbortWithStatus(400)
			return
		}
		result := repos.AccountRepo.InsertAccount(request)
		ctx.AbortWithStatusJSON(200, result)
	}
}

func GetAccount(services *service.Service, repos *repository.Repositories) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		accountId := ctx.Param("accountId")

		result := repos.AccountRepo.GetAccountData(accountId)
		ctx.AbortWithStatusJSON(200, result)
	}
}
