package netflix

import (
	"github.com/gocarina/gocsv"
	"github.com/shekhirin/bionic-cli/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
)

type PlaybackRelatedEvent struct {
	gorm.Model
	ProfileName       string         `csv:"Profile Name" gorm:"uniqueIndex:netflix_playback_related_events_key"`
	TitleDescription  string         `csv:"Title Description" gorm:"uniqueIndex:netflix_playback_related_events_key"`
	Device            string         `csv:"Device"`
	Country           string         `csv:"Country"`
	PlaybackStartTime types.DateTime `csv:"Playback Start Utc Ts" gorm:"uniqueIndex:netflix_playback_related_events_key"`
	Playtraces        []Playtrace    `csv:"Playtraces"`
}

func (r PlaybackRelatedEvent) TableName() string {
	return "netflix_playback_related_events"
}

type Playtrace struct {
	gorm.Model
	PlaybackRelatedEventID int
	PlaybackRelatedEvent   PlaybackRelatedEvent
	EventType              string `json:"eventType"`
	SessionOffsetMs        int    `json:"sessionOffsetMs"`
	MediaOffsetMs          int    `json:"mediaOffsetMs"`
}

func (r Playtrace) TableName() string {
	return "netflix_playtraces"
}

func (p *netflix) importPlaybackRelatedEvents(inputPath string) error {
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return nil
	}

	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	var events []PlaybackRelatedEvent

	if err := gocsv.UnmarshalFile(file, &events); err != nil {
		return err
	}

	err = p.DB().
		Clauses(clause.OnConflict{
			DoNothing: true,
		}).
		CreateInBatches(events, 100).
		Error
	if err != nil {
		return err
	}

	return nil
}
