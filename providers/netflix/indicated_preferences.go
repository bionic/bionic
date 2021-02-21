package netflix

import (
	"github.com/bionic-dev/bionic/types"
	"github.com/gocarina/gocsv"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
)

type IndicatedPreference struct {
	gorm.Model
	ProfileName  string         `csv:"Profile Name" gorm:"uniqueIndex:netflix_indicated_preferences_key"`
	Show         string         `csv:"Show" gorm:"uniqueIndex:netflix_indicated_preferences_key"`
	HasWatched   bool           `csv:"Has Watched" gorm:"uniqueIndex:netflix_indicated_preferences_key"`
	IsInterested bool           `csv:"Is Interested" gorm:"uniqueIndex:netflix_indicated_preferences_key"`
	EventDate    types.DateTime `csv:"Event Date" gorm:"uniqueIndex:netflix_indicated_preferences_key"`
}

func (r IndicatedPreference) TableName() string {
	return tablePrefix + "indicated_preferences"
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
