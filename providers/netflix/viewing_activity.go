package netflix

import (
	"github.com/gocarina/gocsv"
	"github.com/shekhirin/bionic-cli/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
)

type ViewingAction struct {
	gorm.Model
	ProfileName           string         `csv:"Profile Name"`
	StartTime             types.DateTime `csv:"Start Time"`
	Duration              Duration       `csv:"Duration"`
	Attributes            string         `csv:"Attributes"`
	Title                 string         `csv:"Title"`
	SupplementalVideoType string         `csv:"Supplemental Video Type"`
	DeviceType            string         `csv:"Device Type"`
	Bookmark              Duration       `csv:"Bookmark"`
	LatestBookmark        Duration       `csv:"LatestBookmark"`
	Country               string         `csv:"Country"`
}

func (r ViewingAction) TableName() string {
	return "netflix_viewing_activity"
}

func (p *netflix) processViewingActivity(inputPath string) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	var actions []ViewingAction

	if err := gocsv.UnmarshalFile(file, &actions); err != nil { // Load clients from file
		return err
	}

	err = p.db.
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
