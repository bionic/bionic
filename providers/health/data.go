package health

import (
	"encoding/xml"
	"github.com/BionicTeam/bionic/types"
	"gorm.io/gorm"
	"io"
	"os"
	"time"
)

type Data struct {
	gorm.Model
	Locale            string
	ExportDate        types.DateTime    `gorm:"unique"`
	Me                MeRecord          `xml:"Me"`
	Entries           []Entry           `gorm:"-"`
	Workouts          []Workout         `gorm:"-"`
	ActivitySummaries []ActivitySummary `gorm:"-"`
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

type Entry struct {
	gorm.Model
	Type            string           `xml:"type,attr" gorm:"uniqueIndex:health_entries_key"`
	SourceName      string           `xml:"sourceName,attr"`
	SourceVersion   string           `xml:"sourceVersion,attr"`
	Unit            string           `xml:"unit,attr"`
	CreationDate    types.DateTime   `xml:"creationDate,attr" gorm:"uniqueIndex:health_entries_key"`
	StartDate       types.DateTime   `xml:"startDate,attr"`
	EndDate         types.DateTime   `xml:"endDate,attr"`
	Value           string           `xml:"value,attr"`
	Device          string           `xml:"device,attr"`
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
	Bpm     int    `xml:"bpm,attr"`
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
	ActivityType          string          `xml:"workoutActivityType,attr"`
	Duration              float64         `xml:"duration,attr"`
	DurationUnit          string          `xml:"durationUnit,attr"`
	TotalDistance         float64         `xml:"totalDistance,attr"`
	TotalDistanceUnit     string          `xml:"totalDistanceUnit,attr"`
	TotalEnergyBurned     float64         `xml:"totalEnergyBurned,attr"`
	TotalEnergyBurnedUnit string          `xml:"totalEnergyBurnedUnit,attr"`
	SourceName            string          `xml:"sourceName,attr"`
	SourceVersion         string          `xml:"sourceVersion,attr"`
	CreationDate          types.DateTime  `xml:"creationDate,attr" gorm:"unique"`
	StartDate             types.DateTime  `xml:"startDate,attr"`
	EndDate               types.DateTime  `xml:"endDate,attr"`
	Device                string          `xml:"device,attr"`
	MetadataEntries       []MetadataEntry `xml:"MetadataEntry" gorm:"polymorphic:Parent"`
	Events                []WorkoutEvent  `xml:"WorkoutEvent"`
	Route                 WorkoutRoute    `xml:"WorkoutRoute"`
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

func (p *health) importData(inputPath string) error {
	var data Data

	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	decoder := xml.NewDecoder(f)

	for {
		token, err := decoder.Token()
		if token == nil || err == io.EOF {
			break
		} else if err != nil {
			return err
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
				return err
			}
		}
	}

	err = p.DB().
		FirstOrCreate(&data, data.Constraints()).
		Error
	if err != nil {
		return err
	}

	return nil
}

func (p *health) parseHealthData(data *Data, _ *xml.Decoder, start *xml.StartElement) error {
	data.Locale = start.Attr[0].Value

	return nil
}

func (p *health) parseExportDate(data *Data, _ *xml.Decoder, start *xml.StartElement) error {
	if err := data.ExportDate.UnmarshalText([]byte(start.Attr[0].Value)); err != nil {
		return err
	}

	err := p.DB().
		Find(&data, data.Constraints()).
		Error
	if err != nil {
		return err
	}

	data.Me.DataID = data.ID

	return nil
}

func (p *health) parseMe(data *Data, decoder *xml.Decoder, start *xml.StartElement) error {
	if err := decoder.DecodeElement(&data.Me, start); err != nil {
		return err
	}

	return p.DB().
		FirstOrCreate(&data.Me, data.Me.Constraints()).
		Error
}

func (p *health) parseRecord(_ *Data, decoder *xml.Decoder, start *xml.StartElement) error {
	var entry Entry

	if err := decoder.DecodeElement(&entry, start); err != nil {
		return err
	}

	err := p.DB().
		Find(&entry, entry.Constraints()).
		Error
	if err != nil {
		return err
	}

	for i := range entry.MetadataEntries {
		metadataEntry := &entry.MetadataEntries[i]

		metadataEntry.ParentID = entry.ID
		metadataEntry.ParentType = entry.TableName()

		err = p.DB().
			FirstOrCreate(metadataEntry, metadataEntry.Constraints()).
			Error
		if err != nil {
			return err
		}
	}

	for i := range entry.BeatsPerMinutes {
		beatsPerMinute := &entry.BeatsPerMinutes[i]

		beatsPerMinute.EntryID = entry.ID

		err = p.DB().
			FirstOrCreate(beatsPerMinute, beatsPerMinute.Constraints()).
			Error
		if err != nil {
			return err
		}
	}

	return p.DB().
		FirstOrCreate(&entry, entry.Constraints()).
		Error
}

func (p *health) parseWorkout(_ *Data, decoder *xml.Decoder, start *xml.StartElement) error {
	var workout Workout

	if err := decoder.DecodeElement(&workout, start); err != nil {
		return err
	}

	err := p.DB().
		Find(&workout, workout.Constraints()).
		Error
	if err != nil {
		return err
	}

	for i := range workout.MetadataEntries {
		metadataEntry := &workout.MetadataEntries[i]

		metadataEntry.ParentID = workout.ID
		metadataEntry.ParentType = workout.TableName()

		err = p.DB().
			FirstOrCreate(metadataEntry, metadataEntry.Constraints()).
			Error
		if err != nil {
			return err
		}
	}

	for i := range workout.Events {
		event := &workout.Events[i]

		event.WorkoutID = workout.ID

		err = p.DB().
			FirstOrCreate(event, event.Constraints()).
			Error
		if err != nil {
			return err
		}
	}

	route := &workout.Route
	if !time.Time(route.CreationDate).IsZero() {
		route.WorkoutID = workout.ID

		err = p.DB().
			FirstOrCreate(route, route.Constraints()).
			Error
		if err != nil {
			return err
		}
	}

	return p.DB().
		FirstOrCreate(&workout, workout.Constraints()).
		Error
}

func (p *health) parseActivitySummary(_ *Data, decoder *xml.Decoder, start *xml.StartElement) error {
	var activitySummary ActivitySummary

	if err := decoder.DecodeElement(&activitySummary, start); err != nil {
		return err
	}

	return p.DB().
		FirstOrCreate(&activitySummary, activitySummary.Constraints()).
		Error
}
