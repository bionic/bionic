package netflix

import (
	"github.com/shekhirin/bionic-cli/providers/provider"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"path"
)

type netflix struct {
	db *gorm.DB
}

func New(db *gorm.DB) provider.Provider {
	return &netflix{
		db: db,
	}
}

func (p *netflix) Models() []schema.Tabler {
	return []schema.Tabler{
		&ViewingAction{},
	}
}

func (p *netflix) Process(inputPath string) error {
	if !provider.IsPathDir(inputPath) {
		return provider.ErrInputPathShouldBeDirectory
	}

	if err := p.processViewingActivity(path.Join(inputPath, "Content_Interaction", "ViewingActivity.csv")); err != nil {
		return err
	}

	return nil
}
