package orderservice

import (
	"context"
	"pismo-dev/api/types"
	"pismo-dev/error/service"
	"pismo-dev/internal/models"
)

type OrderServiceInterface interface {
	GetOrderDetails(ctx context.Context, filters map[string][]string, pagination *models.Pagination) (models.OrderDetails, *service.ServiceError)
	GetUserCollectionMintCounts(ctx context.Context, userId string, collectionAddress string) (types.GetUserCollectionMintCountsResponse, *service.ServiceError)
}
