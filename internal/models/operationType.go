package models

import "github.com/jinzhu/gorm"

type OperationType struct {
	gorm.Model
	Description string
}
