package google

import (
	"github.com/bionic-dev/bionic/providers/google"
	"github.com/bionic-dev/bionic/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

const searchQueryPrefix = "Searched for "

type Search struct {
	gorm.Model
	Text     string
	Time     types.DateTime
	ActionID int `gorm:"unique"`
	Action   google.Action
}

func (Search) TableName() string {
	return tablePrefix + "searches"
}

func (Search) Update(db *gorm.DB) error {
	var results []google.Action
	query := db.
		Model(&google.Action{}).
		Where("Title like ? AND Header = 'Search'", "%"+searchQueryPrefix+"%")
	query.FindInBatches(&results, 100, func(tx *gorm.DB, batch int) error {
		var items []Search
		for _, action := range results {
			items = append(items, Search{
				Text:   strings.TrimPrefix(action.Title, searchQueryPrefix),
				Time:   action.Time,
				Action: action,
			})
		}
		err := db.
			Clauses(clause.OnConflict{
				DoNothing: true,
			}).
			Create(items).
			Error
		return err
	})

	return nil
}

func (Search) Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&Search{})
	return err
}
