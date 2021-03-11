package google

import (
	"github.com/bionic-dev/bionic/imports/provider"
	"github.com/bionic-dev/bionic/internal/provider/database"
	"gorm.io/gorm"
	"path"
)

const name = "google"
const tablePrefix = "google_"

const locationHistoryFile = "Location History.json"
const semanticLocationDirectoryName = "Semantic Location History"

type google struct {
	database.Database
}

func New(db *gorm.DB) provider.Provider {
	return &google{
		Database: database.New(db),
	}
}

func (google) Name() string {
	return name
}

func (google) TablePrefix() string {
	return tablePrefix
}

func (google) ImportDescription() string {
	return "https://takeout.google.com/"
}

func (p *google) Migrate() error {
	err := p.DB().AutoMigrate(
		&Action{},
		&Product{},
		&LocationInfo{},
		&Subtitle{},
		&Detail{},
		&LocationHistoryItem{},
		&LocationActivity{},
		&LocationActivityTypeCandidate{},
		&ActivitySegment{},
		&ActivityTypeCandidate{},
		&ActivityPathPoint{},
		&TransitStop{},
		&Waypoint{},
		&PlaceVisit{},
		&PlacePathPoint{},
		&CandidateLocation{},
	)
	if err != nil {
		return err
	}

	if err := p.DB().SetupJoinTable(&Action{}, "Products", &ActionProductAssoc{}); err != nil {
		return err
	}

	return nil
}

func (p *google) ImportFns(inputPath string) ([]provider.ImportFn, error) {
	directoryProviders := []provider.ImportFn{
		provider.NewImportFn(
			"Activity",
			p.importActivityFromDirectory,
			path.Join(inputPath, "My Activity"),
		),
		provider.NewImportFn(
			"Location History",
			p.importLocationHistoryFromFile,
			path.Join(inputPath, "Location History", locationHistoryFile),
		),
		provider.NewImportFn(
			"Semantic Location History",
			p.importSemanticLocationHistoryFromDirectory,
			path.Join(inputPath, "Location History", semanticLocationDirectoryName),
		),
	}
	archiveProviders := []provider.ImportFn{
		provider.NewImportFn(
			"Activity",
			p.importActivityFromArchive,
			inputPath,
		),
		provider.NewImportFn(
			"Location History",
			p.importLocationHistoryFromArchive,
			inputPath,
		),
		provider.NewImportFn(
			"Semantic Location History",
			p.importSemanticLocationHistoryFromArchive,
			inputPath,
		),
	}

	if provider.IsPathDir(inputPath) {
		return directoryProviders, nil
	} else {
		return archiveProviders, nil
	}
}
