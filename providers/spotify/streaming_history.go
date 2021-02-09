package spotify

import (
	"encoding/json"
	"github.com/BionicTeam/bionic/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io/ioutil"
	"log"
	"path"
	"strings"
)

const historyFilePrefix = "StreamingHistory"

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
	files, err := ioutil.ReadDir(inputPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if !strings.HasPrefix(f.Name(), historyFilePrefix) {
			continue
		}

		bytes, err := ioutil.ReadFile(path.Join(inputPath, f.Name()))
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
