package database

import "gorm.io/gorm"

type MockModel struct {
	gorm.Model
}

func (MockModel) TableName() string {
	return "mock_model"
}
