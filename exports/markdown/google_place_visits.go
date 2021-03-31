package markdown

import (
	"fmt"
	"github.com/bionic-dev/bionic/imports/google"
	"gorm.io/gorm"
	"time"
)

func (p *markdown) googlePlaceVisits() error {
	var visits []google.PlaceVisit

	locations := map[string]bool{}

	p.DB().
		FindInBatches(&visits, 100, func(tx *gorm.DB, batch int) error {
			for _, visit := range visits {
				localTime := time.Time(visit.DurationStartTimestampMs).Local()

				datePage := p.pageForDate(localTime)

				if !locations[visit.LocationName] {
					p.pages = append(p.pages, &Page{
						Title: visit.LocationName,
						Tag:   "location",
					})
					locations[visit.LocationName] = true
				}

				datePage.Children = append(datePage.Children, Child{
					String: fmt.Sprintf(
						"[[%s]] for %s",
						visit.LocationName,
						formatDuration(
							time.Time(visit.DurationEndTimestampMs).Sub(time.Time(visit.DurationStartTimestampMs)),
							time.Minute,
						),
					),
					Type: ChildGooglePlaceVisit,
					Time: localTime,
				})
			}
			return nil
		})

	return nil
}
