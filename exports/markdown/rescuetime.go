package markdown

import (
	"fmt"
	"github.com/bionic-dev/bionic/imports/rescuetime"
	"gorm.io/gorm"
	"time"
)

const rescuetimeMinSecondsDuration = 5 * 60

func (p *markdown) rescueTime() error {
	var data []rescuetime.ActivityHistoryItem

	categories := map[string]bool{}
	classes := map[string]bool{}

	p.DB().
		Where("duration > ?", rescuetimeMinSecondsDuration).
		Group("category,class,timestamp").
		Select("category", "class", "timestamp", "sum(duration) duration", "min(id) id").
		FindInBatches(&data, 100, func(tx *gorm.DB, batch int) error {
			for _, item := range data {
				timestamp := time.Time(item.Timestamp)
				_, timestampOffset := timestamp.Zone()
				_, localOffset := time.Now().Zone()

				utcTime := timestamp.UTC().Add(time.Duration(timestampOffset) * time.Second)
				localTime := utcTime.Local().Add(time.Duration(-localOffset) * time.Second)

				datePage := p.pageForDate(localTime)

				if !categories[item.Category] {
					p.pages = append(p.pages, &Page{
						Title: item.Category,
						Tag:   "category",
					})
					categories[item.Category] = true
				}

				if !classes[item.Class] {
					p.pages = append(p.pages, &Page{
						Title: item.Class,
						Tag:   "class",
					})
					classes[item.Class] = true
				}

				datePage.Children = append(datePage.Children, Child{
					String: fmt.Sprintf(
						"[[%s]], [[%s]] for %s",
						item.Category,
						item.Class,
						(time.Second * time.Duration(item.Duration)).String(),
					),
					Type: ChildRescueTime,
					Time: localTime,
				})
			}

			return nil
		})

	return nil
}
