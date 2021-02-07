package health

import (
	"encoding/xml"
)

func (p *health) parseWorkout(data *Data, decoder *xml.Decoder, start *xml.StartElement) error {
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

	route := workout.Route
	if route != nil {
		route.WorkoutID = workout.ID
	}

	data.Workouts = append(data.Workouts, &workout)

	return p.DB().
		FirstOrCreate(&workout, workout.Constraints()).
		Error
}
