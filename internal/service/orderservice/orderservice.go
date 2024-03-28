package orderservice

import (
	"context"
	"pismo-dev/api/types"
	"pismo-dev/constants"
	"pismo-dev/error/service"
	"pismo-dev/internal/models"
	"pismo-dev/internal/repository/ordermetadata"
)

type OrderService struct {
	ordermetadata ordermetadata.OrderMetadataRepositoryInterface
}

func NewOrderService(ordermetadata ordermetadata.OrderMetadataRepositoryInterface) *OrderService {
	return &OrderService{ordermetadata: ordermetadata}

}

func (orderService *OrderService) GetOrderDetails(ctx context.Context, filters map[string][]string, pagination *models.Pagination) (models.OrderDetails, *service.ServiceError) {
	var result = orderService.ordermetadata.GetAllOrders(filters, pagination)
	if result.Error != nil {
		return models.OrderDetails{}, service.New(result.Error.Error(), constants.ERROR_TYPES[constants.BAD_REQUEST_ERROR].ErrorCode)
	}
	metadata := result.Result.(models.OrderDetails)
	return metadata, nil
}

func (orderService *OrderService) GetUserCollectionMintCounts(ctx context.Context, userId string, collectionAddress string) (types.GetUserCollectionMintCountsResponse, *service.ServiceError) {
	totalOGmintCount := orderService.ordermetadata.GetTotalOGMintCount(collectionAddress)
	if totalOGmintCount < 0 {
		return types.GetUserCollectionMintCountsResponse{}, service.New("Could not retrieve Mint Count", constants.ERROR_TYPES[constants.DB_ERROR].ErrorCode, constants.ERROR_TYPES[constants.DB_ERROR].HttpStatus)
	}
	userMintsCount := orderService.ordermetadata.GetUserNftCollectionOrderCountByOrderTypeAndStatus(userId, collectionAddress, constants.CHECK_ALLOW_OG_MINT_STATUS, constants.OG_MINT)
	if userMintsCount < 0 {
		return types.GetUserCollectionMintCountsResponse{}, service.New("Could not retrieve User Mint Count", constants.ERROR_TYPES[constants.DB_ERROR].ErrorCode, constants.ERROR_TYPES[constants.DB_ERROR].HttpStatus)
	}
	return types.GetUserCollectionMintCountsResponse{
		UserMints:  userMintsCount,
		TotalMints: totalOGmintCount,
	}, nil
}
