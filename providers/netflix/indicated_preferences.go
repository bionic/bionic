package netflix

import (
	"github.com/gocarina/gocsv"
	"github.com/shekhirin/bionic-cli/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
)

type IndicatedPreference struct {
	gorm.Model
	ProfileName  string         `csv:"Profile Name" gorm:"uniqueIndex:idx_preference"`
	Show         string         `csv:"Show" gorm:"uniqueIndex:idx_preference"`
	HasWatched   bool           `csv:"Has Watched" gorm:"uniqueIndex:idx_preference"`
	IsInterested bool           `csv:"Is Interested" gorm:"uniqueIndex:idx_preference"`
	EventDate    types.DateTime `csv:"Event Date" gorm:"uniqueIndex:idx_preference"`
}

func (r IndicatedPreference) TableName() string {
	return "netflix_indicated_preferences"
}

func (p *netflix) importIndicatedPreferences(inputPath string) error {
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return nil
	}

	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	var preferences []IndicatedPreference

	if err := gocsv.UnmarshalFile(file, &preferences); err != nil {
		return err
	}

	err = p.DB().
		Clauses(clause.OnConflict{
			DoNothing: true,
		}).
		CreateInBatches(preferences, 100).
		Error
	if err != nil {
		return err
	}

	return nil
}
