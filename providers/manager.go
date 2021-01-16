package providers

import (
	"errors"
	"fmt"
	"github.com/shekhirin/bionic-cli/database"
	"github.com/shekhirin/bionic-cli/providers/netflix"
	"github.com/shekhirin/bionic-cli/providers/provider"
	"github.com/shekhirin/bionic-cli/providers/twitter"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

var ErrProviderNotFound = errors.New("provider not found")

type Manager struct {
	db        *gorm.DB
	providers map[string]provider.Provider
}

func DefaultProviders(db *gorm.DB) []provider.Provider {
	return []provider.Provider{
		twitter.New(db),
		netflix.New(db),
	}
}

func NewManager(db *gorm.DB, providers []provider.Provider) (*Manager, error) {
	manager := &Manager{
		db:        db,
		providers: map[string]provider.Provider{},
	}

	for _, p := range providers {
		manager.providers[p.Name()] = p
	}

	if err := db.AutoMigrate(&database.Import{}); err != nil {
		return nil, err
	}

	for _, p := range manager.providers {
		if err := migrate(db, p); err != nil {
			return nil, err
		}
	}

	return manager, nil
}

func (m Manager) GetByName(name string) (provider.Provider, error) {
	p, ok := m.providers[name]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrProviderNotFound, name)
	}

	return p, nil
}

func (m Manager) Reset(p provider.Provider) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		err := tx.
			Where("provider = ?", p.Name()).
			Delete(&database.Import{}).
			Error
		if err != nil {
			return err
		}

		rows, err := tx.
			Table("sqlite_master").
			Select("name").
			Where("type = 'table' AND name LIKE ?", p.TablePrefix()+"%").
			Rows()
		if err != nil {
			return err
		}

		var tables []string

		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				return err
			}

			tables = append(tables, name)
		}

		for _, table := range tables {
			if err := tx.Exec("DROP TABLE IF EXISTS ?", clause.Table{Name: table}).Error; err != nil {
				return err
			}
		}

		if err := migrate(tx, p); err != nil {
			return err
		}

		return nil
	})
}

func migrate(db *gorm.DB, p provider.Provider) error {
	return db.AutoMigrate(tablersToInterfaces(p.Models())...)
}

func tablersToInterfaces(tablers []schema.Tabler) []interface{} {
	interfaces := make([]interface{}, len(tablers))
	for i, tabler := range tablers {
		interfaces[i] = tabler
	}
	return interfaces
}
