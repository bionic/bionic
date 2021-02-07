package health

import (
	"encoding/xml"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"path/filepath"
)

func (p *health) parseWorkoutRoute(f io.Reader, route *WorkoutRoute) error {
	gpxRoute := WorkoutRouteGPX(*route)

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	if err := xml.Unmarshal(bytes, &gpxRoute); err != nil {
		return err
	}

	*route = WorkoutRoute(gpxRoute)

	return nil
}

func (p *health) importWorkoutRoutes(data *Data, getWorkoutRouteReader func(name string) io.Reader) error {
	var g errgroup.Group

	for i := range data.Workouts {
		workoutRoute := data.Workouts[i].Route

		if workoutRoute != nil {
			if rc := getWorkoutRouteReader(filepath.Base(workoutRoute.FilePath)); rc != nil {
				g.Go(func() error {
					return p.parseWorkoutRoute(rc, workoutRoute)
				})
			}
		}
	}

	if err := g.Wait(); err != nil {
		return err
	}

	for _, workout := range data.Workouts {
		if workout.Route == nil {
			continue
		}

		err := p.DB().
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
			FirstOrCreate(workout.Route, workout.Route.Constraints()).
			Error
		if err != nil {
			return err
		}
	}

	return nil
}
