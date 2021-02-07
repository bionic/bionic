package health

import (
	"encoding/xml"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"path/filepath"
)

func (p *health) importWorkoutRoutes(data *Data, files map[string]io.ReadCloser) error {
	var g errgroup.Group

	for i := range data.Workouts {
		workoutRoute := data.Workouts[i].Route

		g.Go(func() error {
			if workoutRoute != nil {
				if r, ok := files[filepath.Base(workoutRoute.FilePath)]; ok {
					gpxRoute := WorkoutRouteGPX(*workoutRoute)

					bytes, err := ioutil.ReadAll(r)
					if err != nil {
						return err
					}

					if err := r.Close(); err != nil {
						return err
					}

					if err := xml.Unmarshal(bytes, &gpxRoute); err != nil {
						return err
					}

					*workoutRoute = WorkoutRoute(gpxRoute)
				}
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	for _, workout := range data.Workouts {
		if workout.Route == nil {
			continue
		}

		err := p.DB().
			Select("ID").
			Find(workout.Route, workout.Route.Constraints()).
			Error
		if err != nil {
			return err
		}

		for i := range workout.Route.TrackPoints {
			trackPoint := &workout.Route.TrackPoints[i]
			trackPoint.WorkoutRouteID = workout.Route.ID

			err = p.DB().
				FirstOrCreate(trackPoint, trackPoint.Constraints()).
				Error
			if err != nil {
				return err
			}
		}

		err = p.DB().
			Session(&gorm.Session{CreateBatchSize: 100}).
			Omit("MetadataEntries", "TrackPoints").
			Where(workout.Route.Constraints()).
			Updates(workout.Route).
			Error
		if err != nil {
			return err
		}
	}

	return nil
}
