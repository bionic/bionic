package markdown

import (
	"fmt"
	"github.com/bionic-dev/bionic/exports/provider"
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
