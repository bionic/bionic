package exports

import (
	"errors"
	"fmt"
	"github.com/bionic-dev/bionic/exports/markdown"
	"github.com/bionic-dev/bionic/exports/provider"
	"gorm.io/gorm"
)

var ErrProviderNotFound = errors.New("provider not found")

type Manager struct {
	db        *gorm.DB
	providers map[string]provider.Provider
}

func DefaultProviders(db *gorm.DB) []provider.Provider {
	return []provider.Provider{
		markdown.New(db),
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

	return manager, nil
}

func (m Manager) GetByName(name string) (provider.Provider, error) {
	p, ok := m.providers[name]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrProviderNotFound, name)
	}

	return p, nil
}
