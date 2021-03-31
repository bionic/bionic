package provider

import (
	"github.com/bionic-dev/bionic/internal/provider/database"
)

type Provider interface {
	database.Database
	Name() string
	Export(outputPath string) error
}
