package netflix

import (
	"github.com/gocarina/gocsv"
	"github.com/BionicTeam/bionic/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
)

type SearchHistoryItem struct {
	gorm.Model
	ProfileName    string         `csv:"Profile Name" gorm:"uniqueIndex:netflix_search_history_key"`
	CountryIsoCode string         `csv:"Country Iso Code"`
	Device         string         `csv:"Device" gorm:"uniqueIndex:netflix_search_history_key"`
	IsKids         bool           `csv:"Is Kids"`
	QueryTyped     string         `csv:"Query Typed"`
	DisplayedName  string         `csv:"Displayed Name"`
	Action         string         `csv:"Action"`
	Section        string         `csv:"Section"`
	Time           types.DateTime `csv:"Utc Timestamp" gorm:"uniqueIndex:netflix_search_history_key"`
}

func (r SearchHistoryItem) TableName() string {
	return "netflix_search_history"
}

func (p *netflix) importSearchHistory(inputPath string) error {
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return nil
	}

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
