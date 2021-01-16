package provider

import (
	"errors"
	"gorm.io/gorm/schema"
	"os"
)

var ErrInputPathShouldBeDirectory = errors.New("input path should be directory")

type ImportFn struct {
	name      string
	fn        func(inputPath string) error
	inputPath string
}

func NewImportFn(name string, fn func(inputPath string) error, inputPath string) ImportFn {
	return ImportFn{
		name:      name,
		fn:        fn,
		inputPath: inputPath,
	}
}

func (fn ImportFn) Name() string {
	return fn.name
}

func (fn ImportFn) Call() error {
	return fn.fn(fn.inputPath)
}

type Provider interface {
	Database
	Name() string
	TablePrefix() string
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
