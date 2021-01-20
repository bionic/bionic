package view

import "gorm.io/gorm"

type View interface {
	TableName() string
	Update(db *gorm.DB) error
}
