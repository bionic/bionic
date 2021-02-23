package database

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"path"
)

func New(dbPath string) (*gorm.DB, error) {
	if err := os.MkdirAll(path.Dir(dbPath), 0755); err != nil && !os.IsExist(err) {
		return nil, err
	}

	if _, err := os.Stat(dbPath); err != nil {
		if os.IsNotExist(err) {
			if _, err := os.Create(dbPath); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	loggerConfig := logger.Config{LogLevel: logger.Silent}

	if viper.GetBool("verbose") {
		loggerConfig.LogLevel = logger.Warn
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.New(logrus.StandardLogger(), loggerConfig),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetTables(db *gorm.DB) ([]string, error) {
	rows, err := db.
		Table("sqlite_master").
		Select("name").
		Where("type = 'table' AND name NOT LIKE 'sqlite_%'").
		Rows()
	if err != nil {
		return nil, err
	}

	var tables []string

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}

		tables = append(tables, name)
	}

	return tables, nil
}
