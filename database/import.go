package database

import "gorm.io/gorm"

type Import struct {
	gorm.Model
	Provider string
}
