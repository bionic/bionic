package health

import (
	"encoding/xml"
	"io/ioutil"
	"path"
)

func (p *health) parseWorkoutRoute(inputPath string, route *WorkoutRoute) error {
	gpxRoute := WorkoutRouteGPX(*route)

	bytes, err := ioutil.ReadFile(path.Join(path.Dir(inputPath), gpxRoute.FilePath))
	if err != nil {
		return err
	}

	if err := xml.Unmarshal(bytes, &gpxRoute); err != nil {
		return err
	}

	*route = WorkoutRoute(gpxRoute)

	return nil
}
