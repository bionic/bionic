package markdown

import (
	"fmt"
	"github.com/bionic-dev/bionic/views/google"
	"gorm.io/gorm"
	"time"
)

func (p *markdown) googleSearches() error {
	var searches []google.Search

	p.DB().
		FindInBatches(&searches, 100, func(tx *gorm.DB, batch int) error {
			for _, search := range searches {
				localTime := time.Time(search.Time).Local()

				datePage := p.pageForDate(localTime)
				datePage.Children = append(datePage.Children, Child{
					String: fmt.Sprintf("Searched in Google for '%s'", search.Text),
					Type:   ChildGooglePlaceVisit,
					Time:   localTime,
				})
			}

			return nil
		})

	return nil
}
