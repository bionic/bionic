package health

import (
	"github.com/BionicTeam/bionic/providers/provider"
	"gorm.io/gorm"
	"path"
)

const name = "health"
const tablePrefix = "health_"

type health struct {
	provider.Database
}

func New(db *gorm.DB) provider.Provider {
	return &health{
		Database: provider.NewDatabase(db),
	}
}

func (health) Name() string {
	return name
}

func (health) TablePrefix() string {
	return tablePrefix
}

func (p *health) Migrate() error {
	return p.DB().AutoMigrate(
		&DataExport{},
		&MeRecord{},
		&Device{},
		&Entry{},
		&EntryMetadataItem{},
		&BeatsPerMinute{},
		&Workout{},
		&WorkoutMetadataItem{},
		&WorkoutEvent{},
		&WorkoutRoute{},
		&WorkoutRouteMetadataItem{},
		&ActivitySummary{},
		&WorkoutRouteTrackPoint{},
	)
}

func (p *health) ImportFns(inputPath string) ([]provider.ImportFn, error) {
	directoryProviders := []provider.ImportFn{
		provider.NewImportFn(
			"Data Export",
			p.importDataExportFromDirectory,
			path.Join(inputPath, "export.xml"),
		),
	}
	archiveProviders := []provider.ImportFn{
		provider.NewImportFn(
			"Data Export",
			p.importDataExportFromArchive,
			inputPath,
		),
	}

	if provider.IsPathDir(inputPath) {
		return directoryProviders, nil
	} else {
		return archiveProviders, nil
	}
}
