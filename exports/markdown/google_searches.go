package markdown

import (
	"fmt"
	"github.com/bionic-dev/bionic/views/google"
	"time"
)

func (p *markdown) googleSearches() error {
	var searches []google.Search

	p.DB().
		Model(google.Search{}).
		Find(&searches)

	for _, search := range searches {
		datePage := p.pageForDate(time.Time(search.Time))

		datePage.Children = append(datePage.Children, Child{
			String: fmt.Sprintf("Searched in Google for '%s'", search.Text),
			Type:   ChildGooglePlaceVisit,
			Time:   time.Time(search.Time),
		})
	}

	return nil
}
