package google

import (
	"archive/zip"
	"encoding/json"
	"github.com/bionic-dev/bionic/providers/provider"
	"github.com/bionic-dev/bionic/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type LocationHistoryItem struct {
	gorm.Model
	Accuracy         int                `json:"accuracy"`
	Activities       []LocationActivity `json:"activity"`
	Altitude         int                `json:"altitude"`
	Heading          int                `json:"heading"`
	LatitudeE7       int                `json:"latitudeE7" gorm:"uniqueIndex:google_location_history_key"`
	LongitudeE7      int                `json:"longitudeE7" gorm:"uniqueIndex:google_location_history_key"`
	Time             types.DateTime     `json:"timestampMs" gorm:"uniqueIndex:google_location_history_key"`
	Velocity         int                `json:"velocity"`
	VerticalAccuracy int                `json:"verticalAccuracy"`
}

func (LocationHistoryItem) TableName() string {
	return tablePrefix + "location_history"
}

type LocationActivity struct {
	gorm.Model
	LocationHistoryItemID int
	LocationHistoryItem   LocationHistoryItem
	TypeCandidates        []LocationActivityTypeCandidate `json:"activity"`
	Time                  types.DateTime                  `json:"timestampMs"`
}

func (LocationActivity) TableName() string {
	return tablePrefix + "location_activity"
}

type LocationActivityType string

const (
	LocationActivityExitingVehicle LocationActivityType = "EXITING_VEHICLE"
	LocationActivityInRailVehicle  LocationActivityType = "IN_RAIL_VEHICLE"
	LocationActivityInRoadVehicle  LocationActivityType = "IN_ROAD_VEHICLE"
	LocationActivityInVehicle      LocationActivityType = "IN_VEHICLE"
	LocationActivityOnBicycle      LocationActivityType = "ON_BICYCLE"
	LocationActivityOnFoot         LocationActivityType = "ON_FOOT"
	LocationActivityRunning        LocationActivityType = "RUNNING"
	LocationActivityStill          LocationActivityType = "STILL"
	LocationActivityTilting        LocationActivityType = "TILTING"
	LocationActivityUnknown        LocationActivityType = "UNKNOWN"
	LocationActivityWalking        LocationActivityType = "WALKING"
)

type LocationActivityTypeCandidate struct {
	gorm.Model
	LocationActivityID int
	LocationActivity   LocationActivity
	Confidence         int                  `json:"confidence"`
	Type               LocationActivityType `json:"type"`
}

func (LocationActivityTypeCandidate) TableName() string {
	return tablePrefix + "location_activity_type_candidates"
}

func (p *google) importLocationHistoryFromArchive(inputPath string) error {
	r, err := zip.OpenReader(inputPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = r.Close()
	}()

	for _, f := range r.File {
		filename := filepath.Base(f.Name)
		if filename != locationHistoryFile {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		if err := p.processLocationHistoryFile(rc); err != nil {
			return err
		}
		if err := rc.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (p *google) importLocationHistoryFromFile(inputPath string) error {
	if !provider.IsPathExists(inputPath) {
		return nil
	}

	rc, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	err = p.processLocationHistoryFile(rc)
	if err != nil {
		return err
	}

	return nil
}

func (p *google) processLocationHistoryFile(rc io.ReadCloser) error {
	bytes, err := ioutil.ReadAll(rc)
	if err != nil {
		return err
	}

	var data struct {
		Locations []LocationHistoryItem `json:"locations"`
	}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return err
	}

	if err := p.saveLocationHistory(data.Locations); err != nil {
		return err
	}

	return nil
}

func (p *google) saveLocationHistory(items []LocationHistoryItem) error {
	err := p.DB().
		Clauses(clause.OnConflict{
			DoNothing: true,
		}).
		CreateInBatches(items, 1000).
		Error
	return err
}
