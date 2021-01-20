package google

import (
	"github.com/BionicTeam/bionic/providers/google"
	"github.com/BionicTeam/bionic/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

const queryPrefix = "Searched for "

type Search struct {
	gorm.Model
	Text     string         `gorm:"uniqueIndex:google_searches_key"`
	Time     types.DateTime `gorm:"uniqueIndex:google_searches_key"`
	ActionID int
	Action   google.Action
}

func (Search) TableName() string {
	return "google_searches"
}

func (Search) Update(db *gorm.DB) error {
	err := db.AutoMigrate(&Search{})
	if err != nil {
		return err
	}
	var results []google.Action
	db.Model(&google.Action{}).Where("Title like '%Searched for%' AND Header = 'Search'").FindInBatches(&results, 100, func(tx *gorm.DB, batch int) error {
		var items []Search
		for _, action := range results {
			items = append(items, Search{
				Text:   strings.TrimPrefix(action.Title, queryPrefix),
				Time:   action.Time,
				Action: action,
			})
		}
		err := db.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(items).Error
		return err
	})

	return nil
}
