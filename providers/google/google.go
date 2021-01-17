package google

import (
	"github.com/shekhirin/bionic-cli/providers/provider"
	"gorm.io/gorm"
)

const name = "google"
const tablePrefix = "google_"

type google struct {
	provider.Database
}

func New(db *gorm.DB) provider.Provider {
	return &google{
		Database: provider.NewDatabase(db),
	}
}

func (google) Name() string {
	return name
}

func (google) TablePrefix() string {
	return tablePrefix
}

func (p *google) Migrate() error {
	err := p.DB().AutoMigrate(
		&Action{},
		&Product{},
		&LocationInfo{},
		&Subtitle{},
		&Detail{})
	if err != nil {
		return err
	}

	if err := p.DB().SetupJoinTable(&Action{}, "Products", &ActionProductAssoc{}); err != nil {
		return err
	}

	return nil
}

func (p *google) ImportFns(inputPath string) ([]provider.ImportFn, error) {
	if !provider.IsPathDir(inputPath) {
		//return nil, provider.ErrInputPathShouldBeDirectory FIXME Add conditional work on zip/folder
	}

	return []provider.ImportFn{
		provider.NewImportFn(
			"Search Activity",
			p.importActivity,
			inputPath,
		),
	}, nil
}
