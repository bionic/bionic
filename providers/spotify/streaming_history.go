package spotify

import (
	"encoding/json"
	"github.com/BionicTeam/bionic/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io/ioutil"
	"path"
	"path/filepath"
)

const historyFileMask = "StreamingHistory*.json"

type StreamingHistoryItem struct {
	gorm.Model
	EndTime    types.DateTime `json:"endTime" gorm:"uniqueIndex:spotify_streaming_activity_key"`
	ArtistName string         `json:"artistName" gorm:"uniqueIndex:spotify_streaming_activity_key"`
	TrackName  string         `json:"trackName" gorm:"uniqueIndex:spotify_streaming_activity_key"`
	MsPlayed   int            `json:"msPlayed"`
}

func (StreamingHistoryItem) TableName() string {
	return tablePrefix + "streaming_history"
}

func (p *spotify) importStreamingHistory(inputPath string) error {
	files, err := filepath.Glob(path.Join(inputPath, historyFileMask))
	if err != nil {
		return err
	}

	for _, f := range files {
		bytes, err := ioutil.ReadFile(f)
		if err != nil {
			return err
		}

		var items []StreamingHistoryItem
		err = json.Unmarshal(bytes, &items)
		if err != nil {
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
	}

	return nil
}
