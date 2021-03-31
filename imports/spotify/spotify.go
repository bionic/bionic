package spotify

import (
	"github.com/bionic-dev/bionic/imports/provider"
	"github.com/bionic-dev/bionic/internal/provider/database"
	"gorm.io/gorm"
)

const name = "spotify"
const tablePrefix = "spotify_"

type spotify struct {
	database.Database
}

func New(db *gorm.DB) provider.Provider {
	return &spotify{
		Database: database.New(db),
	}
}

func (spotify) Name() string {
	return name
}

func (spotify) TablePrefix() string {
	return tablePrefix
}

func (p *spotify) Migrate() error {
	err := p.DB().AutoMigrate(
		&StreamingHistoryItem{},
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *spotify) ImportFns(inputPath string) ([]provider.ImportFn, error) {
	return []provider.ImportFn{
		provider.NewImportFn(
			"Streaming History",
			p.importStreamingHistory,
			inputPath,
		),
	}, nil
}
