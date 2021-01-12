package netflix

import (
	"encoding/json"
	"github.com/gocarina/gocsv"
	"github.com/shekhirin/bionic-cli/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
)

type PlaybackRelatedEvent struct {
	gorm.Model
	ProfileName       string         `csv:"Profile Name" gorm:"uniqueIndex:idx_playback"`
	TitleDescription  string         `csv:"Title Description" gorm:"uniqueIndex:idx_playback"`
	Device            string         `csv:"Device"`
	Country           string         `csv:"Country"`
	PlaybackStartTime types.DateTime `csv:"Playback Start Utc Ts" gorm:"uniqueIndex:idx_playback"`
	Playtraces        PlaytracesJSON `csv:"Playtraces"`
}

type Playtrace struct {
	gorm.Model
	PlaybackRelatedEventID int
	PlaybackRelatedEvent   PlaybackRelatedEvent
	EventType              string `json:"eventType"`
	SessionOffsetMs        int    `json:"sessionOffsetMs"`
	MediaOffsetMs          int    `json:"mediaOffsetMs"`
}

type PlaytracesJSON []Playtrace

func (p *PlaytracesJSON) UnmarshalCSV(csv string) error {
	err := json.Unmarshal([]byte(csv), &p)
	if err != nil {
		return err
	}

	return nil
}

func (r PlaybackRelatedEvent) TableName() string {
	return "netflix_playback_related_events"
}

func (r Playtrace) TableName() string {
	return "netflix_playtraces"
}

func (p *netflix) importPlaybackRelatedEvents(inputPath string) error {
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
