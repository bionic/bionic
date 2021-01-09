package netflix

import (
	"github.com/shekhirin/bionic-cli/providers/provider"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"path"
)

type netflix struct {
	provider.Database
}

func New(db *gorm.DB) provider.Provider {
	return &netflix{
		Database: provider.NewDatabase(db),
	}
}

func (p *netflix) Models() []schema.Tabler {
	return []schema.Tabler{
		&ViewingAction{},
	}
}

func (p *netflix) ImportFns(inputPath string) ([]provider.ImportFn, error) {
	if !provider.IsPathDir(inputPath) {
		return nil, provider.ErrInputPathShouldBeDirectory
	}

	return []provider.ImportFn{
		provider.NewImportFn(
			"Viewing Activity",
			p.importViewingActivity,
			path.Join(inputPath, "Content_Interaction", "ViewingActivity.csv"),
		),
	}, nil
}
