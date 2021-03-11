package markdown

import (
	"fmt"
	"github.com/bionic-dev/bionic/exports/provider"
	"github.com/bionic-dev/bionic/imports/google"
	"github.com/bionic-dev/bionic/imports/rescuetime"
	"github.com/bionic-dev/bionic/imports/spotify"
	"github.com/bionic-dev/bionic/internal/provider/database"
	"gorm.io/gorm"
	"os"
	"path"
	"sort"
	"strings"
	"time"
)

type ChildType int

const (
	ChildSpotify ChildType = iota + 1
	ChildRescueTime
	ChildGooglePlaceVisit
)

func (ct ChildType) String() string {
	switch ct {
	case ChildSpotify:
		return "Spotify"
	default:
		return ""
	}
}

type Page struct {
	Title    string
	Children []Child
	Tag      string
}

type Child struct {
	String string
	Type   ChildType
}

type markdown struct {
	database.Database

	pages     []*Page
	datePages map[time.Time]*Page
}

func New(db *gorm.DB) provider.Provider {
	return &markdown{
		Database:  database.New(db),
		datePages: make(map[time.Time]*Page),
	}
}

func (p *markdown) Name() string {
	return "markdown"
}

func (p *markdown) ExportDescription() string {
	return "Suitable for import to Roam Research, Obsidian, Athens, etc."
}

func (p *markdown) Export(outputPath string) error {
	if err := p.googlePlaceVisits(); err != nil {
		return err
	}
	if err := p.spotify(); err != nil {
		return err
	}
	if err := p.rescueTime(); err != nil {
		return err
	}

	if err := os.MkdirAll(outputPath, 0755); err != nil && !os.IsExist(err) {
		return err
	}

	for _, page := range p.pages {
		title := strings.Replace(page.Title, "/", "\\", -1)

		file, err := os.OpenFile(path.Join(outputPath, fmt.Sprintf("%s.md", title)), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}

		sort.Slice(page.Children, func(i, j int) bool {
			return page.Children[i].Type < page.Children[j].Type
		})

		var previousChildren ChildType

		for _, children := range page.Children {
			if children.Type != previousChildren {
				file.WriteString(fmt.Sprintf("# %s\n", children.Type))
				previousChildren = children.Type
			}
			file.WriteString(fmt.Sprintf("- %s\n", strings.Replace(children.String, "/", "\\", -1)))
		}

		if page.Tag != "" {
			file.WriteString(fmt.Sprintf("\n#%s", page.Tag))
		}

		file.Close()
	}

	return nil
}

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

func (p *markdown) spotify() error {
	var data []struct {
		Date   string
		Artist string
		Track  string
	}

	p.DB().
		Model(spotify.StreamingHistoryItem{}).
		Distinct(
			"STRFTIME('%Y-%m-%d', end_time) date",
			"artist_name artist",
			"track_name track",
		).
		Find(&data)

	artistPages := map[string]*Page{}
	trackPages := map[string]*Page{}

	for _, item := range data {
		date, err := time.Parse("2006-01-02", item.Date)
		if err != nil {
			return err
		}

		datePage := p.pageForDate(date)

		var artistPage *Page

		artistName := strings.TrimLeft(item.Artist, "#")

		if page, ok := artistPages[artistName]; ok {
			artistPage = page
		} else {
			artistPage = &Page{
				Title: artistName,
				Tag:   "artist",
			}

			artistPages[artistName] = artistPage
			p.pages = append(p.pages, artistPage)
		}

		var trackPage *Page

		trackName := fmt.Sprintf("%s – %s", artistName, item.Track)

		if page, ok := trackPages[trackName]; ok {
			trackPage = page
		} else {
			trackPage = &Page{
				Title: trackName,
				Tag:   "track",
				Children: []Child{
					{
						String: fmt.Sprintf("[[%s]]", artistName),
					},
				},
			}

			trackPages[trackName] = trackPage
			p.pages = append(p.pages, trackPage)
		}

		datePage.Children = append(datePage.Children, Child{
			String: fmt.Sprintf("[[%s]], [[%s]]", artistName, trackName),
			Type:   ChildSpotify,
		})
	}

	return nil
}

func (p *markdown) rescueTime() error {
	var data []struct {
		Date     string
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
		//if domainRegex.Match([]byte(item.Activity)) {
		//	continue
		//}

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

func (p *markdown) pageForDate(date time.Time) *Page {
	if page, ok := p.datePages[date]; ok {
		return page
	} else {
		page := &Page{
			Title: date.Format("2006-01-02"),
			Tag:   "date",
		}

		p.datePages[date] = page
		p.pages = append(p.pages, page)

		return page
	}
}
