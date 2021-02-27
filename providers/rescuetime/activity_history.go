package rescuetime

import (
	"bytes"
	"database/sql"
	"github.com/bionic-dev/bionic/types"
	"github.com/gocarina/gocsv"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io"
	"os"
)

type ActivityHistoryItem struct {
	gorm.Model
	Activity string `csv:"activity" gorm:"uniqueIndex:rescuetime_activity_history_key"`
	// Workaround for NULL in unique index: https://stackoverflow.com/a/8289327
	Details   ActivityHistoryDetails `csv:"details" gorm:"uniqueIndex:rescuetime_activity_history_key,expression:COALESCE(details\\, '')"`
	Category  string                 `csv:"category"`
	Class     string                 `csv:"class"`
	Duration  int                    `csv:"duration"`
	Timestamp types.DateTime         `csv:"timestamp" gorm:"uniqueIndex:rescuetime_activity_history_key"`
}

func (ActivityHistoryItem) TableName() string {
	return tablePrefix + "activity_history"
}

type ActivityHistoryDetails struct {
	sql.NullString
}

func (ahd *ActivityHistoryDetails) UnmarshalText(text []byte) error {
	str := string(text)

	if str == "No Details" {
		return nil
	}

	ahd.Valid = true
	ahd.String = str

	return nil
}

func (p *rescuetime) importActivityHistory(inputPath string) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	reader := io.MultiReader(
		bytes.NewBufferString(`"timestamp","activity","details","category","class","duration"`+"\n"),
		file,
	)

	var activityHistory []ActivityHistoryItem
	if err := gocsv.Unmarshal(reader, &activityHistory); err != nil {
		return err
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
