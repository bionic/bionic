package markdown

import (
	"fmt"
	"github.com/bionic-dev/bionic/imports/google"
	"time"
)

func (p *markdown) googlePlaceVisits() error {
	var visits []google.PlaceVisit

	p.DB().
		Model(google.PlaceVisit{}).
		Find(&visits)

	locations := map[string]bool{}

	for _, visit := range visits {
		datePage := p.pageForDate(time.Time(visit.DurationStartTimestampMs))

		if !locations[visit.LocationName] {
			p.pages = append(p.pages, &Page{
				Title: visit.LocationName,
				Tag:   "location",
			})
			locations[visit.LocationName] = true
		}

		datePage.Children = append(datePage.Children, Child{
			String: fmt.Sprintf("[[%s]]", visit.LocationName),
			Type:   ChildGooglePlaceVisit,
		})
	}

	return nil
}
