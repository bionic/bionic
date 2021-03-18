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
	datePages map[string]*Page
}

func New(db *gorm.DB) provider.Provider {
	return &markdown{
		Database:  database.New(db),
		datePages: make(map[string]*Page),
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
	if err := p.googleSearches(); err != nil {
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

	// Avoid too many opened files
	ch := make(chan struct{}, 100)

	for _, page := range p.pages {
		page := page
		errs.Go(func() error {
			ch <- struct{}{}
			err := page.Write(outputPath)
			<-ch
			return err
		})
	}

	return errs.Wait()
}

func (p *markdown) pageForDate(date time.Time) *Page {
	dateString := date.Format("2006-01-02")

	if page, ok := p.datePages[dateString]; ok {
		return page
	} else {
		page := &Page{
			Title: dateString,
		}

		p.datePages[dateString] = page
		p.pages = append(p.pages, page)

		return page
	}
}
