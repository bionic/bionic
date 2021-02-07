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
		&Data{},
		&MeRecord{},
		&Entry{},
		&BeatsPerMinute{},
		&Workout{},
		&WorkoutEvent{},
		&WorkoutRoute{},
		&ActivitySummary{},
		&MetadataEntry{},
		&WorkoutRouteTrackPoint{},
	)
}

func (p *health) ImportFns(inputPath string) ([]provider.ImportFn, error) {
	if !provider.IsPathDir(inputPath) {
		return nil, provider.ErrInputPathShouldBeDirectory
	}

	return []provider.ImportFn{
		provider.NewImportFn(
			"Data",
			p.importData,
			path.Join(inputPath, "export.xml"),
		),
	}, nil
}

