package health

import (
	"encoding/xml"
	"github.com/BionicTeam/bionic/types"
	"gorm.io/gorm"
)

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

func (p *health) parseWorkout(export *DataExport, decoder *xml.Decoder, start *xml.StartElement) error {
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

	if workout.Device != nil {
		err = p.DB().
			FirstOrCreate(workout.Device, workout.Device.Constraints()).
			Error
		if err != nil {
			return err
		}
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

	route := workout.Route
	if route != nil {
		route.WorkoutID = workout.ID
	}

	export.Workouts = append(export.Workouts, &workout)

	return p.DB().
		FirstOrCreate(&workout, workout.Constraints()).
		Error
}
