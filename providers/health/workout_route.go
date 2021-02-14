package health

import (
	"encoding/xml"
	"github.com/BionicTeam/bionic/types"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"path/filepath"
)

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

func (tp WorkoutRouteTrackPoint) Conditions() map[string]interface{} {
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

func (p *health) importWorkoutRoutes(export *DataExport, files map[string]io.ReadCloser) error {
	var g errgroup.Group

	for i := range export.Workouts {
		workoutRoute := export.Workouts[i].Route

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

	for _, workout := range export.Workouts {
		if workout.Route == nil {
			continue
		}

		err := p.DB().
			Select("ID").
			Find(workout.Route, workout.Route.Conditions()).
			Error
		if err != nil {
			return err
		}

		for i := range workout.Route.TrackPoints {
			trackPoint := &workout.Route.TrackPoints[i]
			trackPoint.WorkoutRouteID = workout.Route.ID

			err = p.DB().
				FirstOrCreate(trackPoint, trackPoint.Conditions()).
				Error
			if err != nil {
				return err
			}
		}

		err = p.DB().
			Session(&gorm.Session{CreateBatchSize: 100}).
			Omit("MetadataEntries", "TrackPoints").
			Where(workout.Route.Conditions()).
			Updates(workout.Route).
			Error
		if err != nil {
			return err
		}
	}

	return nil
}
