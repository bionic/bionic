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

type ActivitySegment struct {
	gorm.Model
	Activities   []ActivityTypeCandidate `json:"activities"`
	ActivityType string                  `json:"activityType"`
	Confidence   string                  `json:"confidence"`
	Distance     int                     `json:"distance"`

	DurationEndTimestampMs   types.DateTime `gorm:"uniqueIndex:google_activity_segments_key"`
	DurationStartTimestampMs types.DateTime `gorm:"uniqueIndex:google_activity_segments_key"`

	EndLocationLatitudeE7  int
	EndLocationLongitudeE7 int

	ParkingEventLocationAccuracyMetres int
	ParkingEventLocationLatitudeE7     int
	ParkingEventLocationLongitudeE7    int
	ParkingEventTimestampMs            types.DateTime

	SimplifiedRawPathPoints []ActivityPathPoint

	StartLocationLatitudeE7  int
	StartLocationLongitudeE7 int

	TransitPathHexRgbColor string
	TransitPathName        string
	TransitStops           []TransitStop

	Waypoints []Waypoint
}

func (s ActivitySegment) Conditions() map[string]interface{} {
	return map[string]interface{}{
		"start_location_latitude_e7":  s.StartLocationLatitudeE7,
		"start_location_longitude_e7": s.StartLocationLongitudeE7,
		"duration_start_timestamp_ms": s.DurationStartTimestampMs,
	}
}

func (ActivitySegment) TableName() string {
	return tablePrefix + "activity_segments"
}

func (s *ActivitySegment) UnmarshalJSON(b []byte) error {
	type Alias ActivitySegment

	var data struct {
		Alias
		Duration struct {
			EndTimestampMs   types.DateTime `json:"endTimestampMs"`
			StartTimestampMs types.DateTime `json:"startTimestampMs"`
		} `json:"duration"`
		EndLocation struct {
			LatitudeE7  int `json:"latitudeE7"`
			LongitudeE7 int `json:"longitudeE7"`
		} `json:"endLocation"`
		ParkingEvent struct {
			Location struct {
				AccuracyMetres int `json:"accuracyMetres"`
				LatitudeE7     int `json:"latitudeE7"`
				LongitudeE7    int `json:"longitudeE7"`
			} `json:"location"`
			TimestampMs types.DateTime `json:"timestampMs"`
		} `json:"parkingEvent"`
		SimplifiedRawPath struct {
			Points []ActivityPathPoint `json:"points"`
		} `json:"simplifiedRawPath"`
		StartLocation struct {
			LatitudeE7  int `json:"latitudeE7"`
			LongitudeE7 int `json:"longitudeE7"`
		} `json:"startLocation"`
		TransitPath struct {
			HexRgbColor  string        `json:"hexRgbColor"`
			Name         string        `json:"name"`
			TransitStops []TransitStop `json:"transitStops"`
		} `json:"transitPath"`
		WaypointPath struct {
			Waypoints []Waypoint `json:"waypoints"`
		} `json:"waypointPath"`
	}

	err := json.Unmarshal(b, &data)
	if err != nil {
		return err
	}

	*s = ActivitySegment(data.Alias)
	s.DurationEndTimestampMs = data.Duration.EndTimestampMs
	s.DurationStartTimestampMs = data.Duration.StartTimestampMs
	s.EndLocationLatitudeE7 = data.EndLocation.LatitudeE7
	s.EndLocationLongitudeE7 = data.EndLocation.LongitudeE7
	s.ParkingEventLocationAccuracyMetres = data.ParkingEvent.Location.AccuracyMetres
	s.ParkingEventLocationLatitudeE7 = data.ParkingEvent.Location.LatitudeE7
	s.ParkingEventLocationLongitudeE7 = data.ParkingEvent.Location.LongitudeE7
	s.ParkingEventTimestampMs = data.ParkingEvent.TimestampMs
	s.SimplifiedRawPathPoints = data.SimplifiedRawPath.Points
	s.StartLocationLatitudeE7 = data.StartLocation.LatitudeE7
	s.StartLocationLongitudeE7 = data.StartLocation.LongitudeE7
	s.TransitPathHexRgbColor = data.TransitPath.HexRgbColor
	s.TransitPathName = data.TransitPath.Name
	s.TransitStops = data.TransitPath.TransitStops
	s.Waypoints = data.WaypointPath.Waypoints

	return nil
}

type ActivityTypeCandidate struct {
	gorm.Model
	ActivitySegment   ActivitySegment
	ActivitySegmentID int
	ActivityType      string  `json:"activityType"`
	Probability       float64 `json:"probability"`
}

func (ActivityTypeCandidate) TableName() string {
	return tablePrefix + "activity_type_candidates"
}

type ActivityPathPoint struct {
	gorm.Model
	ActivitySegment   ActivitySegment
	ActivitySegmentID int
	AccuracyMeters    int            `json:"accuracyMeters"`
	LatE7             int            `json:"latE7"`
	LngE7             int            `json:"lngE7"`
	Time              types.DateTime `json:"timestampMs"`
}

func (ActivityPathPoint) TableName() string {
	return tablePrefix + "activity_path_points"
}

type TransitStop struct {
	gorm.Model
	ActivitySegment   ActivitySegment
	ActivitySegmentID int
	LatitudeE7        int    `json:"latitudeE7"`
	LongitudeE7       int    `json:"longitudeE7"`
	Name              string `json:"name"`
	PlaceID           string `json:"placeId"`
}

func (TransitStop) TableName() string {
	return tablePrefix + "transit_stops"
}

type Waypoint struct {
	gorm.Model
	ActivitySegment   ActivitySegment
	ActivitySegmentID int
	LatE7             int `json:"latE7"`
	LngE7             int `json:"lngE7"`
}

func (Waypoint) TableName() string {
	return tablePrefix + "waypoints"
}

type PlaceVisit struct {
	gorm.Model

	CenterLatE7 int `json:"centerLatE7" gorm:"uniqueIndex:google_place_visits_key"`
	CenterLngE7 int `json:"centerLngE7" gorm:"uniqueIndex:google_place_visits_key"`

	DurationEndTimestampMs   types.DateTime `gorm:"uniqueIndex:google_place_visits_key"`
	DurationStartTimestampMs types.DateTime `gorm:"uniqueIndex:google_place_visits_key"`

	EditConfirmationStatus string

	LocationAddress             string
	LocationLatitudeE7          int
	LocationLocationConfidence  float64
	LocationLongitudeE7         int
	LocationName                string
	LocationPlaceID             string `gorm:"uniqueIndex:google_place_visits_key"`
	LocationSourceInfoDeviceTag int

	OtherCandidateLocations []CandidateLocation

	PlaceConfidence string `json:"placeConfidence"`
	VisitConfidence int    `json:"visitConfidence"`

	SimplifiedRawPathPoints []PlacePathPoint // only for top level (no parent_id)

	ChildVisits []*PlaceVisit `json:"childVisits"`

	PlaceVisit   *PlaceVisit
	PlaceVisitID int
}

func (PlaceVisit) TableName() string {
	return tablePrefix + "place_visits"
}

func (p PlaceVisit) Conditions() map[string]interface{} {
	return map[string]interface{}{
		"center_lat_e7":               p.CenterLatE7,
		"center_lng_e7":               p.CenterLngE7,
		"duration_end_timestamp_ms":   p.DurationEndTimestampMs,
		"duration_start_timestamp_ms": p.DurationStartTimestampMs,
		"location_place_id":           p.LocationPlaceID,
	}
}

func (p *PlaceVisit) UnmarshalJSON(b []byte) error {
	type Alias PlaceVisit

	var data struct {
		Alias

		Duration struct {
			EndTimestampMs   types.DateTime `json:"endTimestampMs"`
			StartTimestampMs types.DateTime `json:"startTimestampMs"`
		} `json:"duration"`

		Location struct {
			Address            string  `json:"address"`
			LatitudeE7         int     `json:"latitudeE7"`
			LocationConfidence float64 `json:"locationConfidence"`
			LongitudeE7        int     `json:"longitudeE7"`
			Name               string  `json:"name"`
			PlaceID            string  `json:"placeId"`
			SourceInfo         struct {
				DeviceTag int `json:"deviceTag"`
			} `json:"sourceInfo"`
		} `json:"location"`

		SimplifiedRawPath struct {
			Points []PlacePathPoint `json:"points"`
		} `json:"simplifiedRawPath"`
	}

	err := json.Unmarshal(b, &data)
	if err != nil {
		return err
	}

	*p = PlaceVisit(data.Alias)
	p.DurationEndTimestampMs = data.Duration.EndTimestampMs
	p.DurationStartTimestampMs = data.Duration.StartTimestampMs
	p.LocationAddress = data.Location.Address
	p.LocationLatitudeE7 = data.Location.LatitudeE7
	p.LocationLocationConfidence = data.Location.LocationConfidence
	p.LocationLongitudeE7 = data.Location.LongitudeE7
	p.LocationName = data.Location.Name
	p.LocationPlaceID = data.Location.PlaceID
	p.LocationSourceInfoDeviceTag = data.Location.SourceInfo.DeviceTag
	p.SimplifiedRawPathPoints = data.SimplifiedRawPath.Points

	return nil
}

type PlacePathPoint struct {
	gorm.Model
	PlaceVisit   PlaceVisit
	PlaceVisitID int

	AccuracyMeters int            `json:"accuracyMeters"`
	LatE7          int            `json:"latE7"`
	LngE7          int            `json:"lngE7"`
	Time           types.DateTime `json:"timestampMs"`
}

func (PlacePathPoint) TableName() string {
	return tablePrefix + "place_path_points"
}

type CandidateLocation struct {
	gorm.Model
	PlaceVisit   PlaceVisit
	PlaceVisitID int

	LatitudeE7         int     `json:"latitudeE7"`
	LocationConfidence float64 `json:"locationConfidence"`
	LongitudeE7        int     `json:"longitudeE7"`
	PlaceID            string  `json:"placeId"`
}

func (CandidateLocation) TableName() string {
	return tablePrefix + "candidate_locations"
}

func (p *google) importSemanticLocationHistoryFromArchive(inputPath string) error {
	r, err := zip.OpenReader(inputPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = r.Close()
	}()

	for _, f := range r.File {
		// files are located in 2-nd level directories: Semantic Location History/2020/*.json
		dataTypeDirectory := filepath.Base(filepath.Dir(filepath.Dir(f.Name)))
		if dataTypeDirectory != semanticLocationDirectoryName {
			continue
		}
		if filepath.Ext(f.Name) != ".json" {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		if err := p.processSemanticLocationFile(rc); err != nil {
			return err
		}
		if err := rc.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (p *google) importSemanticLocationHistoryFromDirectory(inputPath string) error {
	if !provider.IsPathDir(inputPath) {
		return nil
	}

	err := filepath.Walk(inputPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filepath.Ext(info.Name()) != ".json" {
				return nil
			}

			rc, err := os.Open(path)
			if err != nil {
				return err
			}

			err = p.processSemanticLocationFile(rc)
			if err != nil {
				return err
			}

			return nil
		})
	if err != nil {
		return err
	}
	return nil
}

func (p *google) processSemanticLocationFile(rc io.ReadCloser) error {
	bytes, err := ioutil.ReadAll(rc)
	if err != nil {
		return err
	}

	var data struct {
		TimelineObjects []struct {
			ActivitySegment *ActivitySegment `json:"activitySegment"`
			PlaceVisit      *PlaceVisit      `json:"placeVisit"`
		} `json:"timelineObjects"`
	}

	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return err
	}

	for _, timelineObject := range data.TimelineObjects {
		if timelineObject.ActivitySegment != nil {
			var prevActivitySegment ActivitySegment

			activitySegment := *timelineObject.ActivitySegment
			err := p.DB().
				Limit(1).
				Find(&prevActivitySegment, activitySegment.Conditions()).
				Error
			if err != nil {
				return err
			}

			if prevActivitySegment.ID == 0 {
				err = p.DB().
					Clauses(clause.OnConflict{
						DoNothing: true,
					}).
					Create(&activitySegment).
					Error
				if err != nil {
					return err
				}
			}
		}

		if timelineObject.PlaceVisit != nil {
			var prevPlaceVisit PlaceVisit

			placeVisit := *timelineObject.PlaceVisit
			err := p.DB().
				Limit(1).
				Find(&prevPlaceVisit, placeVisit.Conditions()).
				Error
			if err != nil {
				return err
			}

			if prevPlaceVisit.ID == 0 {
				err = p.DB().
					Clauses(clause.OnConflict{
						DoNothing: true,
					}).
					FirstOrCreate(&placeVisit).
					Error
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
