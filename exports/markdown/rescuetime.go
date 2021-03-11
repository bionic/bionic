package markdown

import (
	"fmt"
	"github.com/bionic-dev/bionic/imports/rescuetime"
	"time"
)

func (p *markdown) rescueTime() error {
	var data []struct {
		Date string
		//Activity string
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

	categoryPages := map[string]*Page{}
	classPages := map[string]*Page{}

	for _, item := range data {
		date, err := time.Parse("2006-01-02", item.Date)
		if err != nil {
			return err
		}

		datePage := p.pageForDate(date)

		var categoryPage *Page

		if page, ok := categoryPages[item.Category]; ok {
			categoryPage = page
		} else {
			categoryPage = &Page{
				Title: item.Category,
				Tag:   "category",
			}

			categoryPages[item.Category] = categoryPage
			p.pages = append(p.pages, categoryPage)
		}

		var classPage *Page

		if page, ok := classPages[item.Class]; ok {
			classPage = page
		} else {
			classPage = &Page{
				Title: item.Class,
				Tag:   "class",
			}

			classPages[item.Class] = classPage
			p.pages = append(p.pages, classPage)
		}

		datePage.Children = append(datePage.Children, Child{
			String: fmt.Sprintf("[[%s]], [[%s]]", item.Category, item.Class),
			Type:   ChildRescueTime,
		})
	}

	return nil
}
