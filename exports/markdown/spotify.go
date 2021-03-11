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

	artists := map[string]bool{}
	tracks := map[string]bool{}

	for _, item := range data {
		date, err := time.Parse("2006-01-02", item.Date)
		if err != nil {
			return err
		}

		datePage := p.pageForDate(date)

		artistName := strings.TrimLeft(item.Artist, "#")

		if !artists[artistName] {
			p.pages = append(p.pages, &Page{
				Title: artistName,
				Tag:   "artist",
			})
			artists[artistName] = true
		}

		trackName := fmt.Sprintf("%s – %s", artistName, item.Track)

		if !tracks[trackName] {
			p.pages = append(p.pages, &Page{
				Title: trackName,
				Tag:   "track",
				Children: []Child{
					{
						String: fmt.Sprintf("[[%s]]", artistName),
					},
				},
			})
			tracks[trackName] = true
		}

		datePage.Children = append(datePage.Children, Child{
			String: fmt.Sprintf("[[%s]], [[%s]]", artistName, trackName),
			Type:   ChildSpotify,
		})
	}

	return nil
}
