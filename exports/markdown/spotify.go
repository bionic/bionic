package markdown

import (
	"fmt"
	"github.com/bionic-dev/bionic/imports/spotify"
	"strings"
	"time"
)

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
