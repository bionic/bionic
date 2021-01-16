package netflix

import (
	"github.com/gocarina/gocsv"
	"github.com/shekhirin/bionic-cli/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
)

type ClickstreamAction struct {
	gorm.Model
	ProfileName     string         `csv:"Profile Name" gorm:"uniqueIndex:idx_action"`
	Source          string         `csv:"Source" gorm:"uniqueIndex:idx_action"`
	NavigationLevel string         `csv:"Navigation Level" gorm:"uniqueIndex:idx_action"`
	ReferrerUrl     string         `csv:"Referrer Url"`
	WebpageUrl      string         `csv:"Webpage Url"`
	ClickTime       types.DateTime `csv:"Click Utc Ts" gorm:"uniqueIndex:idx_action"`
}

func (r ClickstreamAction) TableName() string {
	return "netflix_clickstream"
}

func (p *netflix) importClickstream(inputPath string) error {
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return nil
	}

	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	var actions []ClickstreamAction

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
