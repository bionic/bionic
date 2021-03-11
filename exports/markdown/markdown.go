package markdown

import (
	"context"
	"github.com/bionic-dev/bionic/exports/provider"
	"github.com/bionic-dev/bionic/internal/provider/database"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"os"
	"time"
)

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

	errs, _ := errgroup.WithContext(context.Background())

	for _, page := range p.pages {
		page := page
		errs.Go(func() error {
			return page.Write(outputPath)
		})
	}

	return errs.Wait()
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
