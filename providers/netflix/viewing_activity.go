package netflix

import (
	"github.com/gocarina/gocsv"
	"github.com/BionicTeam/bionic/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
)

type ViewingAction struct {
	gorm.Model
	ProfileName           string         `csv:"Profile Name" gorm:"uniqueIndex:netflix_viewing_activity_key"`
	StartTime             types.DateTime `csv:"Start Time" gorm:"uniqueIndex:netflix_viewing_activity_key"`
	Duration              Duration       `csv:"Duration"`
	Attributes            string         `csv:"Attributes"`
	Title                 string         `csv:"Title" gorm:"uniqueIndex:netflix_viewing_activity_key"`
	SupplementalVideoType string         `csv:"Supplemental Video Type"`
	DeviceType            string         `csv:"Device Type"`
	Bookmark              Duration       `csv:"Bookmark"`
	LatestBookmark        Duration       `csv:"Latest Bookmark"`
	Country               string         `csv:"Country"`
}

func (r ViewingAction) TableName() string {
	return tablePrefix + "viewing_activity"
}

func (p *netflix) importViewingActivity(inputPath string) error {
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return nil
	}

	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	var actions []ViewingAction

	if err := gocsv.UnmarshalFile(file, &actions); err != nil {
		return err
	}

	err = p.DB().
		Clauses(clause.OnConflict{
			DoNothing: true,
		}).
		CreateInBatches(actions, 100).
		Error
	if err != nil {
		return err
	}

	return nil
}
