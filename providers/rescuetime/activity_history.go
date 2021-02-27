package rescuetime

import (
	"github.com/bionic-dev/bionic/types"
	"github.com/gocarina/gocsv"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
)

type ActivityHistoryItem struct {
	gorm.Model
	Activity   string `gorm:"uniqueIndex:rescuetime_activity_history_key"`
	Details    *string `gorm:"uniqueIndex:rescuetime_activity_history_key"`
	Category   string
	Class      string
	Duration   int
	Timestamp  types.DateTime `gorm:"uniqueIndex:rescuetime_activity_history_key"`
}

func (ActivityHistoryItem) TableName() string {
	return tablePrefix + "activity_history"
}

func (p *rescuetime) importActivityHistory(inputPath string) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	var rawActivityHistory []struct {
		Timestamp types.DateTime
		Activity  string
		Details   string
		Category  string
		Class     string
		Duration  int `json:",string"`
	}
	if err := gocsv.UnmarshalWithoutHeaders(file, &rawActivityHistory); err != nil {
		return err
	}

	var activityHistory []ActivityHistoryItem

	for _, item := range rawActivityHistory {
		activityHistoryItem := ActivityHistoryItem{
			Activity:  item.Activity,
			Category:  item.Category,
			Class:     item.Class,
			Duration:  item.Duration,
			Timestamp: item.Timestamp,
		}

		if item.Details != "No Details" {
			activityHistoryItem.Details = &item.Details
		}

		activityHistory = append(activityHistory, activityHistoryItem)
	}

	err = p.DB().
		Clauses(clause.OnConflict{
			DoNothing: true,
		}).
		CreateInBatches(activityHistory, 100).
		Error
	if err != nil {
		return err
	}

	return nil
}
