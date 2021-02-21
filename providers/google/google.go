package google

import (
	"github.com/bionic-dev/bionic/providers/provider"
	"gorm.io/gorm"
	"path"
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
		&Detail{},
	)
	if err != nil {
		return err
	}

	if err := p.DB().SetupJoinTable(&Action{}, "Products", &ActionProductAssoc{}); err != nil {
		return err
	}

	return nil
}

func (p *google) ImportFns(inputPath string) ([]provider.ImportFn, error) {
	directoryProviders := []provider.ImportFn{
		provider.NewImportFn(
			"Activity",
			p.importActivityFromDirectory,
			path.Join(inputPath, "My Activity"),
		),
	}
	archiveProviders := []provider.ImportFn{
		provider.NewImportFn(
			"Activity",
			p.importActivityFromArchive,
			inputPath,
		),
	}

	if provider.IsPathDir(inputPath) {
		return directoryProviders, nil
	} else {
		return archiveProviders, nil
	}
}
