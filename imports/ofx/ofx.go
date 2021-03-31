package ofx

import (
	"github.com/bionic-dev/bionic/imports/provider"
	"github.com/bionic-dev/bionic/internal/provider/database"
	"gorm.io/gorm"
)

const name = "ofx"
const tablePrefix = "ofx_"

type OFX struct {
	database.Database
}

func New(db *gorm.DB) provider.Provider {
	return &OFX{
		Database: database.New(db),
	}
}

func (OFX) Name() string {
	return name
}

func (OFX) TablePrefix() string {
	return tablePrefix
}

func (OFX) ImportDescription() string {
	return "OFX is a file format for financial data. " +
		"Export an OFX file from your bank or convert your bank export file to OFX with https://github.com/kedder/ofxstatement."
}

func (p *OFX) Migrate() error {
	err := p.DB().AutoMigrate(
		&Account{},
		&Transaction{},
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *OFX) ImportFns(inputPath string) ([]provider.ImportFn, error) {
	if provider.IsPathDir(inputPath) {
		return nil, provider.ErrInputPathShouldBeFile
	}

	return []provider.ImportFn{
		provider.NewImportFn(
			"Statement",
			p.importStatement,
			inputPath,
		),
	}, nil
}
