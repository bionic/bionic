package markdown

import (
	"fmt"
	"github.com/bionic-dev/bionic/imports/rescuetime"
	"gorm.io/gorm"
	"time"
)

func (p *markdown) rescueTime() error {
	var data []rescuetime.ActivityHistoryItem

	categories := map[string]bool{}
	classes := map[string]bool{}

	p.DB().
		Model(rescuetime.ActivityHistoryItem{}).
		FindInBatches(&data, 100, func(tx *gorm.DB, batch int) error {
			for _, item := range data {
				datePage := p.pageForDate(time.Time(item.Timestamp))

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
					String: fmt.Sprintf("[[%s]], [[%s]]", item.Category, item.Class),
					Type:   ChildRescueTime,
				})
			}

			return nil
		})

	return nil
}
