package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path"
)

func New(dbPath string) (*gorm.DB, error) {
	if err := os.MkdirAll(path.Dir(dbPath), 0755); err != nil && !os.IsExist(err) {
		return nil, err
	}
	if _, err := os.Create(dbPath); err != nil && !os.IsExist(err) {
		return nil, err
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
