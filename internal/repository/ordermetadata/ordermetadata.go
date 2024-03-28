package ordermetadata

import (
	"pismo-dev/constants"
	"pismo-dev/internal/models"
	"pismo-dev/internal/repository/result"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type OrderMetadataRepository struct {
	db        *gorm.DB
	TableName string
}

func NewOrderMetadataRepository(db *gorm.DB) *OrderMetadataRepository {
	return &OrderMetadataRepository{db: db, TableName: "order_metadata"}
}

func (er *OrderMetadataRepository) InsertOrderMetadata(orderMetadata *models.OrderMetadata) result.RepositoryResult {
	err := er.db.Create(&orderMetadata).Error
	if err != nil {
		return result.RepositoryResult{Error: err}
	}
	return result.RepositoryResult{Result: orderMetadata}
}

func (er *OrderMetadataRepository) GetAllOrders(filter map[string][]string, pagination *models.Pagination) result.RepositoryResult {
	var orders []models.OrderMetadataWithTx
	var query = &models.OrderMetadata{}
	offset := (pagination.Page - 1) * pagination.Limit
	queryBuilder := er.db.Limit(pagination.Limit).Offset(offset).Order(pagination.Sort + " " + pagination.Direction)
	if filter != nil {
		if filter["order_id"] != nil && len(filter["order_id"]) > 0 {
			query.OrderId = filter["order_id"][0]
		}
		if filter["vendor_id"] != nil && len(filter["vendor_id"]) > 0 {
			query.VendorId = filter["vendor_id"][0]
		}
		if filter["network_id"] != nil && len(filter["network_id"]) > 0 {
			networkIdList := strings.Split(filter["network_id"][0], ",")
			if len(networkIdList) > 1 {
				queryBuilder = queryBuilder.Where("network_id IN (?)", networkIdList)
			} else {
				query.NetworkId = filter["network_id"][0]
			}

		}
		if filter["user_id"] != nil && len(filter["user_id"]) > 0 {
			userIdList := strings.Split(filter["user_id"][0], ",")
			if len(userIdList) > 1 {
				queryBuilder = queryBuilder.Where("user_id IN (?)", userIdList)
			} else {
				query.UserId = filter["user_id"][0]
			}

		}
		if filter["nft_id"] != nil && len(filter["nft_id"]) > 0 {
			query.NftId = filter["nft_id"][0]
		}
		if filter["entity_address"] != nil && len(filter["entity_address"]) > 0 {
			entityAddressList := strings.Split(filter["entity_address"][0], ",")
			if len(entityAddressList) > 1 {
				queryBuilder = queryBuilder.Where("entity_address IN (?)", entityAddressList)
			} else {
				query.EntityAddress = filter["entity_address"][0]
			}
		}
		if filter["status"] != nil && len(filter["status"]) > 0 {
			if filter["status"][0] == constants.CHECK_ALLOW_OG_MINT {
				queryBuilder = queryBuilder.Where("status IN (?)", constants.CHECK_ALLOW_OG_MINT_STATUS)
			} else {
				query.Status = filter["status"][0]
			}
		}
		if filter["entity_type"] != nil && len(filter["entity_type"]) > 0 {
			query.EntityType = filter["entity_type"][0]
		}
		if filter["order_type"] != nil && len(filter["order_type"]) > 0 {
			if filter["status"] != nil && len(filter["status"]) > 0 && filter["status"][0] == constants.CHECK_ALLOW_OG_MINT {
				query.OrderType = constants.OG_MINT
			} else {
				query.OrderType = filter["order_type"][0]
			}
		}
		if filter["created_at_from"] != nil && len(filter["created_at_from"]) > 0 {
			unixTimestamp, err := strconv.ParseInt(filter["created_at_from"][0], 10, 64)
			var createdAtFrom time.Time
			if err == nil {
				createdAtFrom = time.Unix(unixTimestamp, 0).UTC()
			}
			if !createdAtFrom.IsZero() {
				queryBuilder = queryBuilder.Where("created_at >= ?", createdAtFrom)
			}
		}
		if filter["created_at_to"] != nil && len(filter["created_at_to"]) > 0 {
			unixTimestamp, err := strconv.ParseInt(filter["created_at_to"][0], 10, 64)
			var createdAtTo time.Time
			if err == nil {
				createdAtTo = time.Unix(unixTimestamp, 0).UTC()
			}
			if !createdAtTo.IsZero() {
				queryBuilder = queryBuilder.Where("created_at <= ?", createdAtTo)
			}
		}
	}
	// Change below condition once provided by FM
	queryBuilder = queryBuilder.Table("order_metadata").
		Joins("LEFT OUTER JOIN transaction_data ON order_metadata.order_id = transaction_data.order_id AND transaction_data.order_tx_type = 'ORDER'")
	queryBuilder = queryBuilder.Select("order_metadata.*, transaction_data.tx_hash")
	var count int64
	//	queryBuilder = queryBuilder.Debug()
	err := queryBuilder.Where(query).Find(&orders).Error

	if err != nil {
		return result.RepositoryResult{Error: err}
	}

	err = queryBuilder.Where(query).Limit(-1).Offset(-1).Count(&count).Error
	if err != nil {
		return result.RepositoryResult{Error: err}
	}
	return result.RepositoryResult{Result: models.OrderDetails{Count: count, OrderDetails: orders}}
}

func (er *OrderMetadataRepository) GetOrdersById(orderId string) result.RepositoryResult {
	var orderMetadata models.OrderMetadata
	err := er.db.Where(map[string]interface{}{"order_id": orderId}).First(&orderMetadata).Error
	if err != nil {
		return result.RepositoryResult{Error: err}
	}
	return result.RepositoryResult{Result: orderMetadata}
}

func (er *OrderMetadataRepository) GetPendingNftOrdersForUser(userId string, nftId string) int {
	var query = &models.OrderMetadata{}
	var count int64
	filterString := []string{constants.CREATED, constants.RUNNING, constants.WAITING_FOR_SIGNATURE}
	err := er.db.Model(query).Where("status IN (?) AND nft_id = ? AND user_id = ?", filterString, nftId, userId).Count(&count).Error
	if err != nil {
		return -1
	}
	return int(count)
}

func (er *OrderMetadataRepository) GetUserNftCollectionOrderCountByOrderTypeAndStatus(userId string, contractAddress string, status []string, orderType string) int {
	var query = &models.OrderMetadata{}
	var count int64
	err := er.db.Model(query).Where("user_id = ? AND status IN (?) AND entity_address = ? AND order_type = ?", userId, status, contractAddress, orderType).Count(&count).Error
	if err != nil {
		return -1
	}
	return int(count)
}

func (er *OrderMetadataRepository) GetTotalOGMintCount(contractAddress string) int {
	var query = &models.OrderMetadata{}
	var count int64
	err := er.db.Model(query).Where("status IN (?) AND entity_address = ? AND order_type = ?", constants.CHECK_ALLOW_OG_MINT_STATUS, contractAddress, constants.OG_MINT).Count(&count).Error
	if err != nil {
		return -1
	}
	return int(count)
}

func (er *OrderMetadataRepository) UpdateOrderByOrderId(orderMetadata *models.OrderMetadata) result.RepositoryResult {
	err := er.db.Omit("created_at").Where(models.OrderMetadata{OrderId: orderMetadata.OrderId}).Save(&orderMetadata).Error
	if err != nil {
		return result.RepositoryResult{Error: err}
	}
	return result.RepositoryResult{Result: orderMetadata}
}

func (er *OrderMetadataRepository) UpdateOrderStatusByOrderId(orderId string, status string) result.RepositoryResult {
	err := er.db.Model(&models.OrderMetadata{}).Omit("created_at").Where(models.OrderMetadata{OrderId: orderId}).Update("status", status).Error
	if err != nil {
		return result.RepositoryResult{Error: err}
	}
	return result.RepositoryResult{Result: true}
}

func (er *OrderMetadataRepository) UpdateOrderStatusForFailedOrder(orderId string, status string, executionResponse string) result.RepositoryResult {
	data := map[string]interface{}{
		"status":             status,
		"execution_response": executionResponse,
	}
	err := er.db.Model(&models.OrderMetadata{}).Omit("created_at").Where(models.OrderMetadata{OrderId: orderId}).Updates(data).Error
	if err != nil {
		return result.RepositoryResult{Error: err}
	}
	return result.RepositoryResult{Result: true}
}

func (er *OrderMetadataRepository) GetDbTransaction() (OrderMetadataRepository, *gorm.DB, bool) {
	tx := er.db.Begin()
	return OrderMetadataRepository{TableName: "order_metadata", db: tx}, tx, tx.Error == nil
}

func (er *OrderMetadataRepository) Commit() error {
	return er.db.Commit().Error
}

func (er *OrderMetadataRepository) RollBack() {
	er.db.Rollback()
}
