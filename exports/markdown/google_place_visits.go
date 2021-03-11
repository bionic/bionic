package markdown

import (
	"fmt"
	"github.com/bionic-dev/bionic/imports/google"
	"time"
)

func (p *markdown) googlePlaceVisits() error {
	var data []struct {
		Date   string
		Location string
	}

	p.DB().
		Model(google.PlaceVisit{}).
		Distinct(
			"STRFTIME('%Y-%m-%d', duration_start_timestamp_ms) date",
			"location_name location",
		).
		Find(&data)

	locationPages := map[string]*Page{}

	for _, item := range data {
		date, err := time.Parse("2006-01-02", item.Date)
		if err != nil {
			return err
		}

		datePage := p.pageForDate(date)

		var locationPage *Page

		if page, ok := locationPages[item.Location]; ok {
			locationPage = page
		} else {
			locationPage = &Page{
				Title: item.Location,
				Tag:   "location",
			}

			locationPages[item.Location] = locationPage
			p.pages = append(p.pages, locationPage)
		}

		datePage.Children = append(datePage.Children, Child{
			String: fmt.Sprintf("[[%s]]", item.Location),
			Type:   ChildGooglePlaceVisit,
		})
	}

	return nil
}
