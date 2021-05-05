package chrome

import (
	"github.com/bionic-dev/bionic/imports/provider"
	"github.com/bionic-dev/bionic/internal/provider/database"
	"gorm.io/gorm"
)

const name = "chrome"
const tablePrefix = "chrome_"

type chrome struct {
	database.Database
}

func New(db *gorm.DB) provider.Provider {
	return &chrome{
		Database: database.New(db),
	}
}

func (chrome) Name() string {
	return name
}

func (chrome) TablePrefix() string {
	return tablePrefix
}

func (chrome) ImportDescription() string {
	return "OS X: ~/Library/Application\\ Support/Google/Chrome/Default/History\n" +
		"Windows: C:\\Users\\%USERNAME%\\AppData\\Local\\Google\\Chrome\\User Data\\Default\\History\n" +
		"Linux: ~/.config/google-chrome/Default/databases"
}

func (p *chrome) Migrate() error {
	err := p.DB().AutoMigrate(
		&URL{},
		&Segment{},
		&Visit{},
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *chrome) ImportFns(inputPath string) ([]provider.ImportFn, error) {
	if provider.IsPathDir(inputPath) {
		return nil, provider.ErrInputPathShouldBeFile
	}

	return []provider.ImportFn{
		provider.NewImportFn(
			"History",
			p.importDB,
			inputPath,
		),
	}, nil
}
