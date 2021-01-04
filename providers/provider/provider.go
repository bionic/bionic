package provider

import (
	"errors"
	"os"
)

var ErrInputPathShouldBeDirectory = errors.New("input path should be directory")

type Provider interface {
	Models() []interface{}
	Process(inputPath string) error
}

func IsPathDir(inputPath string) bool {
	info, err := os.Stat(inputPath)
	if err != nil {
		return false
	}

	return info.IsDir()
}
