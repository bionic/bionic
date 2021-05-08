package imports

import (
	"errors"
	"fmt"
	"github.com/bionic-dev/bionic/imports/chrome"
	"github.com/bionic-dev/bionic/imports/google"
	"github.com/bionic-dev/bionic/imports/health"
	"github.com/bionic-dev/bionic/imports/instagram"
	"github.com/bionic-dev/bionic/imports/netflix"
	"github.com/bionic-dev/bionic/imports/ofx"
	"github.com/bionic-dev/bionic/imports/provider"
	"github.com/bionic-dev/bionic/imports/rescuetime"
	"github.com/bionic-dev/bionic/imports/spotify"
	"github.com/bionic-dev/bionic/imports/telegram"
	"github.com/bionic-dev/bionic/imports/twitter"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
		google.New(db),
		telegram.New(db),
		health.New(db),
		spotify.New(db),
		instagram.New(db),
		rescuetime.New(db),
		ofx.New(db),
		chrome.New(db),
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

func (m Manager) Migrate() error {
	if err := m.db.AutoMigrate(&Import{}); err != nil {
		return err
	}

	for _, p := range m.providers {
		if err := p.Migrate(); err != nil {
			return err
		}
	}

	return nil
}

func (m Manager) GetByName(name string) (provider.Provider, error) {
	p, ok := m.providers[name]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrProviderNotFound, name)
	}

	return p, nil
}

func (m Manager) Reset(p provider.Provider) error {
	err := m.db.Transaction(func(tx *gorm.DB) error {
		err := tx.
			Where("provider = ?", p.Name()).
			Delete(&Import{}).
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

		return nil
	})
	if err != nil {
		return err
	}

	return p.Migrate()
}
