package netflix

import (
	"github.com/gocarina/gocsv"
	"github.com/shekhirin/bionic-cli/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
)

type SearchHistoryItem struct {
	gorm.Model
	ProfileName    string         `csv:"Profile Name" gorm:"uniqueIndex:idx_search_history"`
	CountryIsoCode string         `csv:"Country Iso Code"`
	Device         string         `csv:"Device" gorm:"uniqueIndex:idx_search_history"`
	IsKids         bool           `csv:"Is Kids"`
	QueryTyped     string         `csv:"Query Typed"`
	DisplayedName  string         `csv:"Displayed Name"`
	Action         string         `csv:"Action"`
	Section        string         `csv:"Section"`
	Time           types.DateTime `csv:"Utc Timestamp" gorm:"uniqueIndex:idx_search_history"`
}

func (r SearchHistoryItem) TableName() string {
	return "netflix_search_history"
}

func (p *netflix) importSearchHistory(inputPath string) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	var items []SearchHistoryItem

	if err := gocsv.UnmarshalFile(file, &items); err != nil {
		return err
	}

	err = p.DB().
		Clauses(clause.OnConflict{
			DoNothing: true,
		}).
		CreateInBatches(items, 100).
		Error
	if err != nil {
		return err
	}

	return nil
}
