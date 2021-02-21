package netflix

import (
	"github.com/bionic-dev/bionic/types"
	"github.com/gocarina/gocsv"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
)

type InteractiveTitle struct {
	gorm.Model
	ProfileName     string         `csv:"Profile Name" gorm:"uniqueIndex:netflix_interactive_titles_key"`
	TitleDesc       string         `csv:"Title Desc"`
	SelectionType   string         `csv:"Selection Type"`
	ChoiceSegmentId string         `csv:"Choice Segment Id"`
	HasWatched      bool           `csv:"Has Watched"`
	Source          string         `csv:"Source"`
	Time            types.DateTime `csv:"Utc Timestamp" gorm:"uniqueIndex:netflix_interactive_titles_key"`
}

func (r InteractiveTitle) TableName() string {
	return tablePrefix + "interactive_titles"
}

func (p *netflix) importInteractiveTitles(inputPath string) error {
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return nil
	}

	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	var titles []InteractiveTitle

	if err := gocsv.UnmarshalFile(file, &titles); err != nil {
		return err
	}

	err = p.DB().
		Clauses(clause.OnConflict{
			DoNothing: true,
		}).
		CreateInBatches(titles, 100).
		Error
	if err != nil {
		return err
	}

	return nil
}
