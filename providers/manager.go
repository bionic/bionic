package providers

import (
	"errors"
	"fmt"
	"github.com/shekhirin/bionic-cli/database"
	"github.com/shekhirin/bionic-cli/providers/google"
	"github.com/shekhirin/bionic-cli/providers/netflix"
	"github.com/shekhirin/bionic-cli/providers/provider"
	"github.com/shekhirin/bionic-cli/providers/twitter"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var ErrProviderNotFound = errors.New("provider not found")

type Manager struct {
	db        *gorm.DB
	providers map[string]provider.Provider
}

func NewManager(dbPath string) (*Manager, error) {
	db, err := database.New(dbPath)
	if err != nil {
		return nil, err
	}

	manager := &Manager{
		db: db,
		providers: map[string]provider.Provider{
			"twitter": twitter.New(db),
			"netflix": netflix.New(db),
			"google":  google.New(db),
		},
	}

	for _, p := range manager.providers {
		if err := manager.migrate(manager.db, p); err != nil {
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
		if err := tx.Migrator().DropTable(tablersToInterfaces(p.Models())...); err != nil {
			return err
		}

		if err := m.migrate(tx, p); err != nil {
			return err
		}

		return nil
	})
}

func (m Manager) migrate(db *gorm.DB, p provider.Provider) error {
	return db.AutoMigrate(tablersToInterfaces(p.Models())...)
}

func tablersToInterfaces(tablers []schema.Tabler) []interface{} {
	interfaces := make([]interface{}, len(tablers))
	for i, tabler := range tablers {
		interfaces[i] = tabler
	}
	return interfaces
}
