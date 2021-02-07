package health

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"github.com/BionicTeam/bionic/types"
	"gorm.io/gorm"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Data struct {
	gorm.Model
	Locale            string
	ExportDate        types.DateTime     `gorm:"unique"`
	Me                MeRecord           `xml:"Me"`
	Workouts          []*Workout         `gorm:"-"`
}

func (Data) TableName() string {
	return tablePrefix + "data"
}

func (d Data) Constraints() map[string]interface{} {
	return map[string]interface{}{
		"export_date": d.ExportDate,
	}
}

type MeRecord struct {
	gorm.Model
	DataID                      uint           `gorm:"unique"`
	DateOfBirth                 types.DateTime `xml:"HKCharacteristicTypeIdentifierDateOfBirth,attr"`
	BiologicalSex               string         `xml:"HKCharacteristicTypeIdentifierBiologicalSex,attr"`
	BloodType                   string         `xml:"HKCharacteristicTypeIdentifierBloodType,attr"`
	FitzpatrickSkinType         string         `xml:"HKCharacteristicTypeIdentifierFitzpatrickSkinType,attr"`
	CardioFitnessMedicationsUse string         `xml:"HKCharacteristicTypeIdentifierCardioFitnessMedicationsUse,attr"`
}

func (MeRecord) TableName() string {
	return tablePrefix + "me_records"
}

func (m MeRecord) Constraints() map[string]interface{} {
	return map[string]interface{}{
		"data_id": m.DataID,
	}
}

type Device struct {
	gorm.Model
	Name         string
	Manufacturer string
	DeviceModel  string `gorm:"column:model"`
	Hardware     string
	Software     string
}

func (Device) TableName() string {
	return tablePrefix + "devices"
}

func (d Device) Constraints() map[string]interface{} {
	return map[string]interface{}{
		"name":         d.Name,
		"manufacturer": d.Manufacturer,
		"model":        d.DeviceModel,
		"hardware":     d.Hardware,
		"software":     d.Software,
	}
}

func (d *Device) UnmarshalText(b []byte) error {
	text := string(b)

	if len(text) < 3 {
		return nil
	}

	parts := strings.Split(text[1:len(text)-1], ", ")
	if len(parts) < 2 {
		return nil
	}

	attributes := parts[1:]

	for _, attr := range attributes {
		attrParts := strings.Split(attr, ":")
		if len(attrParts) != 2 {
			continue
		}

		key, value := attrParts[0], attrParts[1]

		switch key {
		case "name":
			d.Name = value
		case "manufacturer":
			d.Manufacturer = value
		case "model":
			d.DeviceModel = value
		case "hardware":
			d.Hardware = value
		case "software":
			d.Software = value
		}
	}

	return nil
}

type Entry struct {
	gorm.Model
	Type            string         `xml:"type,attr" gorm:"uniqueIndex:health_entries_key"`
	SourceName      string         `xml:"sourceName,attr"`
	SourceVersion   string         `xml:"sourceVersion,attr"`
	Unit            string         `xml:"unit,attr"`
	CreationDate    types.DateTime `xml:"creationDate,attr" gorm:"uniqueIndex:health_entries_key"`
	StartDate       types.DateTime `xml:"startDate,attr"`
	EndDate         types.DateTime `xml:"endDate,attr"`
	Value           string         `xml:"value,attr"`
	DeviceID        *int
	Device          *Device          `xml:"device,attr"`
	MetadataEntries []MetadataEntry  `xml:"MetadataEntry" gorm:"polymorphic:Parent"`
	BeatsPerMinutes []BeatsPerMinute `xml:"HeartRateVariabilityMetadataList"`
}

func (Entry) TableName() string {
	return tablePrefix + "entries"
}

func (e Entry) Constraints() map[string]interface{} {
	return map[string]interface{}{
		"type":          e.Type,
		"creation_date": e.CreationDate,
	}
}

func (e *Entry) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	type Alias Entry

	var data struct {
		Alias
		HeartRateVariabilityMetadataList struct {
			InstantaneousBeatsPerMinute []BeatsPerMinute `xml:"InstantaneousBeatsPerMinute"`
		} `xml:"HeartRateVariabilityMetadataList"`
	}

	if err := decoder.DecodeElement(&data, &start); err != nil {
		return err
	}

	*e = Entry(data.Alias)

	e.BeatsPerMinutes = data.HeartRateVariabilityMetadataList.InstantaneousBeatsPerMinute

	return nil
}

type BeatsPerMinute struct {
	gorm.Model
	EntryID uint   `gorm:"uniqueIndex:health_beats_per_minutes_key"`
	BPM     int    `xml:"bpm,attr"`
	Time    string `xml:"time,attr" gorm:"uniqueIndex:health_beats_per_minutes_key"`
}

func (BeatsPerMinute) TableName() string {
	return tablePrefix + "beats_per_minutes"
}

func (bpm BeatsPerMinute) Constraints() map[string]interface{} {
	return map[string]interface{}{
		"entry_id": bpm.EntryID,
		"time":     bpm.Time,
	}
}

type Workout struct {
	gorm.Model
	ActivityType          string         `xml:"workoutActivityType,attr"`
	Duration              float64        `xml:"duration,attr"`
	DurationUnit          string         `xml:"durationUnit,attr"`
	TotalDistance         float64        `xml:"totalDistance,attr"`
	TotalDistanceUnit     string         `xml:"totalDistanceUnit,attr"`
	TotalEnergyBurned     float64        `xml:"totalEnergyBurned,attr"`
	TotalEnergyBurnedUnit string         `xml:"totalEnergyBurnedUnit,attr"`
	SourceName            string         `xml:"sourceName,attr"`
	SourceVersion         string         `xml:"sourceVersion,attr"`
	CreationDate          types.DateTime `xml:"creationDate,attr" gorm:"unique"`
	StartDate             types.DateTime `xml:"startDate,attr"`
	EndDate               types.DateTime `xml:"endDate,attr"`
	DeviceID              *int
	Device                *Device         `xml:"device,attr"`
	MetadataEntries       []MetadataEntry `xml:"MetadataEntry" gorm:"polymorphic:Parent"`
	Events                []WorkoutEvent  `xml:"WorkoutEvent"`
	Route                 *WorkoutRoute   `xml:"WorkoutRoute"`
}

func (Workout) TableName() string {
	return tablePrefix + "workouts"
}

func (w Workout) Constraints() map[string]interface{} {
	return map[string]interface{}{
		"creation_date": w.CreationDate,
	}
}

type WorkoutEvent struct {
	gorm.Model
	WorkoutID    uint           `gorm:"uniqueIndex:health_workout_events_key"`
	Type         string         `xml:"type,attr" gorm:"uniqueIndex:health_workout_events_key"`
	Date         types.DateTime `xml:"date,attr" gorm:"uniqueIndex:health_workout_events_key"`
	Duration     float64        `xml:"duration,attr"`
	DurationUnit string         `xml:"durationUnit,attr"`
}

func (WorkoutEvent) TableName() string {
	return tablePrefix + "workout_events"
}

func (e WorkoutEvent) Constraints() map[string]interface{} {
	return map[string]interface{}{
		"workout_id": e.WorkoutID,
		"type":       e.Type,
		"date":       e.Date,
	}
}

type WorkoutRoute struct {
	gorm.Model
	WorkoutID       uint            `gorm:"uniqueIndex:health_workout_routes_key"`
	SourceName      string          `xml:"sourceName,attr"`
	SourceVersion   string          `xml:"sourceVersion,attr"`
	CreationDate    types.DateTime  `xml:"creationDate,attr" gorm:"uniqueIndex:health_workout_routes_key"`
	StartDate       types.DateTime  `xml:"startDate,attr"`
	EndDate         types.DateTime  `xml:"endDate,attr"`
	MetadataEntries []MetadataEntry `xml:"MetadataEntry" gorm:"polymorphic:Parent"`
	FilePath        string
	Time            types.DateTime
	TrackName       string
	TrackPoints     []WorkoutRouteTrackPoint
}

func (WorkoutRoute) TableName() string {
	return tablePrefix + "workout_routes"
}

func (wr WorkoutRoute) Constraints() map[string]interface{} {
	return map[string]interface{}{
		"workout_id":    wr.WorkoutID,
		"creation_date": wr.CreationDate,
	}
}

func (wr *WorkoutRoute) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	type Alias WorkoutRoute

	var data struct {
		Alias
		FileReference struct {
			Path string `xml:"path,attr"`
		} `xml:"FileReference"`
	}

	if err := decoder.DecodeElement(&data, &start); err != nil {
		return err
	}

	*wr = WorkoutRoute(data.Alias)

	wr.FilePath = data.FileReference.Path

	return nil
}

type WorkoutRouteGPX WorkoutRoute

func (wr *WorkoutRouteGPX) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	type Alias WorkoutRoute

	var data = struct {
		Alias
		XMLName  xml.Name `xml:"gpx"`
		Metadata struct {
			Time types.DateTime `xml:"time"`
		} `xml:"metadata"`
		Track struct {
			Name    string `xml:"name"`
			Segment struct {
				Points []struct {
					WorkoutRouteTrackPoint
					Extensions WorkoutRouteTrackPointExtensions `xml:"extensions"`
				} `xml:"trkpt"`
			} `xml:"trkseg"`
		} `xml:"trk"`
	}{
		Alias: Alias(*wr),
	}

	if err := decoder.DecodeElement(&data, &start); err != nil {
		return err
	}

	*wr = WorkoutRouteGPX(data.Alias)

	wr.Time = data.Metadata.Time
	wr.TrackName = data.Track.Name
	for _, point := range data.Track.Segment.Points {
		trackPoint := point.WorkoutRouteTrackPoint
		trackPoint.WorkoutRouteTrackPointExtensions = point.Extensions
		wr.TrackPoints = append(wr.TrackPoints, trackPoint)
	}

	return nil
}

type WorkoutRouteTrackPoint struct {
	gorm.Model
	WorkoutRouteID uint           `gorm:"uniqueIndex:health_track_points_key"`
	Lon            float64        `xml:"lon,attr"`
	Lat            float64        `xml:"lat,attr"`
	Ele            float64        `xml:"ele"`
	Time           types.DateTime `xml:"time" gorm:"uniqueIndex:health_track_points_key"`
	WorkoutRouteTrackPointExtensions
}

func (WorkoutRouteTrackPoint) TableName() string {
	return tablePrefix + "workout_route_track_points"
}

func (tp WorkoutRouteTrackPoint) Constraints() map[string]interface{} {
	return map[string]interface{}{
		"workout_route_id": tp.WorkoutRouteID,
		"time":             tp.Time,
	}
}

type WorkoutRouteTrackPointExtensions struct {
	Speed  float64 `xml:"speed"`
	Course float64 `xml:"course"`
	HAcc   float64 `xml:"hAcc"`
	VAcc   float64 `xml:"vAcc"`
}

type ActivitySummary struct {
	gorm.Model
	Date                   types.DateTime `xml:"dateComponents,attr" gorm:"unique"`
	ActiveEnergyBurned     float64        `xml:"activeEnergyBurned,attr"`
	ActiveEnergyBurnedGoal int            `xml:"activeEnergyBurnedGoal,attr"`
	ActiveEnergyBurnedUnit string         `xml:"activeEnergyBurnedUnit,attr"`
	AppleMoveTime          int            `xml:"appleMoveTime,attr"`
	AppleMoveTimeGoal      int            `xml:"appleMoveTimeGoal,attr"`
	AppleExerciseTime      int            `xml:"appleExerciseTime,attr"`
	AppleExerciseTimeGoal  int            `xml:"appleExerciseTimeGoal,attr"`
	AppleStandHours        int            `xml:"appleStandHours,attr"`
	AppleStandHoursGoal    int            `xml:"appleStandHoursGoal,attr"`
}

func (ActivitySummary) TableName() string {
	return tablePrefix + "activity_summaries"
}

func (as ActivitySummary) Constraints() map[string]interface{} {
	return map[string]interface{}{
		"date": as.Date,
	}
}

type MetadataEntry struct {
	gorm.Model
	ParentID   uint   `gorm:"uniqueIndex:health_metadata_entries_key"`
	ParentType string `gorm:"uniqueIndex:health_metadata_entries_key"`
	Key        string `xml:"key,attr" gorm:"uniqueIndex:health_metadata_entries_key"`
	Value      string `xml:"value,attr"`
}

func (MetadataEntry) TableName() string {
	return tablePrefix + "metadata_entries"
}

func (e MetadataEntry) Constraints() map[string]interface{} {
	return map[string]interface{}{
		"parent_id":   e.ParentID,
		"parent_type": e.ParentType,
		"key":         e.Key,
	}
}

func (p *health) importDataFromArchive(inputPath string) error {
	var data *Data

	r, err := zip.OpenReader(inputPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = r.Close()
	}()

	workoutRouteFiles := map[string]io.ReadCloser{}

	for _, f := range r.File {
		if filepath.Base(f.Name) == "export.xml" {
			rc, err := f.Open()
			if err != nil {
				return err
			}

			data, err = p.importData(rc)
			if err != nil {
				return err
			}

			if err := rc.Close(); err != nil {
				return err
			}
		} else if filepath.Base(filepath.Dir(f.Name)) == "workout-routes" {
			rc, err := f.Open()
			if err != nil {
				return nil
			}

			workoutRouteFiles[filepath.Base(f.Name)] = rc
		}
	}

	if data == nil {
		return errors.New("no export.xml file found")
	}

	return p.importWorkoutRoutes(data, workoutRouteFiles)
}

func (p *health) importDataFromDirectory(inputPath string) error {
	var data *Data

	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	data, err = p.importData(f)
	if err != nil {
		return err
	}

	workoutRouteFiles := map[string]io.ReadCloser{}

	for _, workout := range data.Workouts {
		if route := workout.Route; route != nil {
			r, err := os.Open(path.Join(path.Dir(inputPath), route.FilePath))
			if err != nil {
				return err
			}

			workoutRouteFiles[filepath.Base(route.FilePath)] = r
		}
	}

	return p.importWorkoutRoutes(data, workoutRouteFiles)
}

func (p *health) importData(r io.Reader) (*Data, error) {
	var data Data

	decoder := xml.NewDecoder(r)

	for {
		token, err := decoder.Token()
		if token == nil || err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		switch typ := token.(type) {
		case xml.StartElement:
			var parseFn func(*Data, *xml.Decoder, *xml.StartElement) error

			switch typ.Name.Local {
			case "HealthData":
				parseFn = p.parseHealthData
			case "ExportDate":
				parseFn = p.parseExportDate
			case "Me":
				parseFn = p.parseMe
			case "Record":
				parseFn = p.parseRecord
			case "Workout":
				parseFn = p.parseWorkout
			case "ActivitySummary":
				parseFn = p.parseActivitySummary
			default:
				continue
			}

			if err := parseFn(&data, decoder, &typ); err != nil {
				return nil, err
			}
		}
	}

	err := p.DB().
		FirstOrCreate(&data, data.Constraints()).
		Error
	if err != nil {
		return nil, err
	}

	return &data, nil
}
