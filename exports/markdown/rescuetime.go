package markdown

import (
	"fmt"
	"github.com/bionic-dev/bionic/imports/rescuetime"
	"gorm.io/gorm"
	"strings"
	"time"
)

const rescuetimeMinSecondsDuration = 5 * 60

func (p *markdown) rescueTime() error {
	var data []rescuetime.ActivityHistoryItem

	categories := map[string]bool{}
	classes := map[string]bool{}

	err := p.DB().
		Exec(fmt.Sprintf(`
CREATE TEMP TABLE rescuetime_activity_history_agg AS
SELECT ROW_NUMBER() OVER (ORDER BY category, class, timestamp) id, *
FROM (
	SELECT category, class, timestamp, SUM(duration) duration
	FROM %s
	GROUP BY category, class, timestamp
) t
`, rescuetime.ActivityHistoryItem{}.TableName())).
		Error
	if err != nil {
		return err
	}

	var mergedItem *rescuetime.ActivityHistoryItem
	var offset int

	itemsByTimestamps := map[time.Time][]*rescuetime.ActivityHistoryItem{}

	err = p.DB().
		Unscoped().
		Model(&rescuetime.ActivityHistoryItem{}).
		Table("rescuetime_activity_history_agg").
		FindInBatches(&data, 100, func(tx *gorm.DB, batch int) error {
			for i, item := range data {
				offset++

				if mergedItem != nil {
					if mergedItem.Category == item.Category && mergedItem.Class == item.Class &&
						time.Time(mergedItem.Timestamp).Add(time.Duration(offset)*time.Hour).Equal(time.Time(item.Timestamp)) {
						mergedItem.Duration += item.Duration
						continue
					} else {
						p.insertRescuetimeItem(categories, classes, itemsByTimestamps, mergedItem)
					}
				}

				mergedItem = &data[i]
				offset = 0
			}

			return nil
		}).
		Error
	if err != nil {
		return err
	}

	if mergedItem != nil {
		p.insertRescuetimeItem(categories, classes, itemsByTimestamps, mergedItem)
	}

	for timestamp, items := range itemsByTimestamps {
		datePage := p.pageForDate(timestamp)

		activity := make([]string, len(items))

		for i, item := range items {
			activity[i] = fmt.Sprintf(
				"[[%s]] [[%s]] for %s",
				item.Category,
				item.Class,
				formatDuration(time.Second*time.Duration(item.Duration), time.Minute),
			)
		}

		datePage.Children = append(datePage.Children, Child{
			String: strings.Join(activity, ", "),
			Type:   ChildRescueTime,
			Time:   timestamp,
		})
	}

	return nil
}

func (p *markdown) insertRescuetimeItem(
	categories, classes map[string]bool,
	items map[time.Time][]*rescuetime.ActivityHistoryItem,
	item *rescuetime.ActivityHistoryItem,
) {
	if item.Duration < rescuetimeMinSecondsDuration {
		return
	}

	timestamp := time.Time(item.Timestamp)
	_, timestampOffset := timestamp.Zone()
	_, localOffset := time.Now().Zone()

	utcTime := timestamp.UTC().Add(time.Duration(timestampOffset) * time.Second)
	localTime := utcTime.Local().Add(time.Duration(-localOffset) * time.Second)

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

	if _, ok := items[localTime]; !ok {
		items[localTime] = make([]*rescuetime.ActivityHistoryItem, 0)
	}

	items[localTime] = append(items[localTime], item)
}
