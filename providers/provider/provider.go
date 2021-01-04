package provider

import (
	"errors"
	"gorm.io/gorm/schema"
	"os"
)

var ErrInputPathShouldBeDirectory = errors.New("input path should be directory")

type Provider interface {
	Models() []schema.Tabler
	Process(inputPath string) error
}

func IsPathDir(inputPath string) bool {
	info, err := os.Stat(inputPath)
	if err != nil {
		return false
	}

	return info.IsDir()
}
