package netflix

import (
	"github.com/gocarina/gocsv"
	"github.com/shekhirin/bionic-cli/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
)

type InteractiveTitle struct {
	gorm.Model
	ProfileName     string         `csv:"Profile Name" gorm:"uniqueIndex:idx_title"`
	TitleDesc       string         `csv:"Title Desc"`
	SelectionType   string         `csv:"Selection Type"`
	ChoiceSegmentId string         `csv:"Choice Segment Id"`
	HasWatched      bool           `csv:"Has Watched"`
	Source          string         `csv:"Source"`
	Time            types.DateTime `csv:"Utc Timestamp" gorm:"uniqueIndex:idx_title"`
}

func (r InteractiveTitle) TableName() string {
	return "netflix_interactive_titles"
}

func (p *netflix) importInteractiveTitles(inputPath string) error {
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
