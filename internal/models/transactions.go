// models/transaction.go
package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	AccountID       uint
	OperationTypeID uint
	Amount          float64
	EventDate       string `gorm:"type:timestamp"`
}
