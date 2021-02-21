package instagram

import (
	"github.com/bionic-dev/bionic/providers/provider"
	"gorm.io/gorm"
	"path/filepath"
)

const name = "instagram"
const tablePrefix = "instagram_"

type instagram struct {
	provider.Database
}

func New(db *gorm.DB) provider.Provider {
	return &instagram{
		Database: provider.NewDatabase(db),
	}
}

func (instagram) Name() string {
	return name
}

func (instagram) TablePrefix() string {
	return tablePrefix
}

func (p *instagram) Migrate() error {
	err := p.DB().AutoMigrate(
		&AccountHistoryItem{},
		&RegistrationInfo{},
		&User{},
		&CommentUserMention{},
		&Hashtag{},
		&CommentHashtagMention{},
		&Like{},
		&Comment{},
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *instagram) ImportFns(inputPath string) ([]provider.ImportFn, error) {
	if !provider.IsPathDir(inputPath) {
		return nil, provider.ErrInputPathShouldBeDirectory
	}

	var importFns []provider.ImportFn

	globs := map[string]func(path string) provider.ImportFn{
		"account_history.json": func(path string) provider.ImportFn {
			return provider.NewImportFn(
				"Account History",
				p.importAccountHistory,
				path,
			)
		},
		"comments.json": func(path string) provider.ImportFn {
			return provider.NewImportFn(
				"Comments",
				p.importComments,
				path,
			)
		},
		"likes.json": func(path string) provider.ImportFn {
			return provider.NewImportFn(
				"Likes",
				p.importLikes,
				path,
			)
		},
	}

	for glob, importFnForPath := range globs {
		files, err := filepath.Glob(filepath.Join(inputPath, glob))
		if err != nil {
			return nil, err
		}

		if files != nil {
			importFns = append(importFns, importFnForPath(files[0]))
		}
	}

	return importFns, nil
}
