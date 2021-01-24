package views

import (
	"github.com/BionicTeam/bionic/views/google"
	"github.com/BionicTeam/bionic/views/view"
	"gorm.io/gorm"
)

type Manager struct {
	db    *gorm.DB
	Views []view.View
}

func DefaultViews() []view.View {
	return google.Views
}

func NewManager(db *gorm.DB, views []view.View) (*Manager, error) {
	manager := &Manager{
		db:    db,
		Views: []view.View{},
	}

	for _, v := range views {
		manager.Views = append(manager.Views, v)
	}

	return manager, nil
}

func (m *Manager) Migrate() error {
	for _, v := range m.Views {
		err := v.Migrate(m.db)
		if err != nil {
			return err
		}
	}

	return nil
}
