package provider

import (
	"errors"
	"gorm.io/gorm/schema"
	"os"
)

var ErrInputPathShouldBeDirectory = errors.New("input path should be directory")

type ImportFn struct {
	Fn        func(inputPath string) error
	InputPath string
}

func (pf ImportFn) Call() error {
	return pf.Fn(pf.InputPath)
}

type Provider interface {
	Database
	Models() []schema.Tabler
	ImportFns(inputPath string) ([]ImportFn, error)
}

func IsPathDir(inputPath string) bool {
	info, err := os.Stat(inputPath)
	if err != nil {
		return false
	}

	return info.IsDir()
}
