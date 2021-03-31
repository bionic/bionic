package rescuetime

import (
	"github.com/bionic-dev/bionic/imports/provider"
	"github.com/bionic-dev/bionic/internal/provider/database"
	"gorm.io/gorm"
)

const name = "rescuetime"
const tablePrefix = "rescuetime_"

type rescuetime struct {
	database.Database
}

func New(db *gorm.DB) provider.Provider {
	return &rescuetime{
		Database: database.New(db),
	}
}

func (rescuetime) Name() string {
	return name
}

func (rescuetime) TablePrefix() string {
	return tablePrefix
}

func (rescuetime) ImportDescription() string {
	return "https://www.rescuetime.com/accounts/your-data => \"Your Logged Time\" => \"Activity report history\""
}

func (p *rescuetime) Migrate() error {
	err := p.DB().AutoMigrate(
		&ActivityHistoryItem{},
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *rescuetime) ImportFns(inputPath string) ([]provider.ImportFn, error) {
	if provider.IsPathDir(inputPath) {
		return nil, provider.ErrInputPathShouldBeFile
	}

	return []provider.ImportFn{
		provider.NewImportFn(
			"Activity History",
			p.importActivityHistory,
			inputPath,
		),
	}, nil
}
