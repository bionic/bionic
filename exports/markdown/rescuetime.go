package markdown

import (
	"fmt"
	"github.com/bionic-dev/bionic/imports/rescuetime"
	"time"
)

func (p *markdown) rescueTime() error {
	var data []struct {
		Date     string
		Category string
		Class    string
	}

	p.DB().
		Model(rescuetime.ActivityHistoryItem{}).
		Distinct(
			"STRFTIME('%Y-%m-%d', timestamp) date",
			"category",
			"class",
		).
		Find(&data)

	categories := map[string]bool{}
	classes := map[string]bool{}

	for _, item := range data {
		date, err := time.Parse("2006-01-02", item.Date)
		if err != nil {
			return err
		}

		datePage := p.pageForDate(date)

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
}
