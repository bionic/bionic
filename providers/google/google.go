package google

import (
	"github.com/shekhirin/bionic-cli/providers/provider"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type google struct {
	provider.Database
}

func New(db *gorm.DB) provider.Provider {
	return &google{
		Database: provider.NewDatabase(db),
	}
}

func (p *google) Models() []schema.Tabler {
	return []schema.Tabler{
		&Action{},
		&Product{},
		&LocationInfo{},
		&Subtitle{},
		&Detail{},
	}
}

func (p *google) ImportFns(inputPath string) ([]provider.ImportFn, error) {
	if !provider.IsPathDir(inputPath) {
		//return nil, provider.ErrInputPathShouldBeDirectory FIXME
	}

	return []provider.ImportFn{
		provider.NewImportFn(
			"Search Activity",
			p.importActivity,
			inputPath,
		),
	}, nil
}
