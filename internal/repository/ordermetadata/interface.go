package ordermetadata

import (
	"pismo-dev/internal/models"
	"pismo-dev/internal/repository/result"

	"gorm.io/gorm"
)

type OrderMetadataRepositoryInterface interface {
	InsertOrderMetadata(orderMetadata *models.OrderMetadata) result.RepositoryResult
	GetAllOrders(filters map[string][]string, pagination *models.Pagination) result.RepositoryResult
	GetOrdersById(orderId string) result.RepositoryResult
	GetPendingNftOrdersForUser(userId string, nftId string) int
	UpdateOrderByOrderId(orderMetadata *models.OrderMetadata) result.RepositoryResult
	GetTotalOGMintCount(contractAddress string) int
	UpdateOrderStatusForFailedOrder(orderId string, status string, executionResponse string) result.RepositoryResult
	GetUserNftCollectionOrderCountByOrderTypeAndStatus(userId string, contractAddress string, status []string, orderType string) int
	UpdateOrderStatusByOrderId(orderId string, status string) result.RepositoryResult
	GetDbTransaction() (OrderMetadataRepository, *gorm.DB, bool)
	Commit() error
	RollBack()
}
