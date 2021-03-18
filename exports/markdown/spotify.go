package markdown

import (
	"fmt"
	"github.com/bionic-dev/bionic/imports/spotify"
	"gorm.io/gorm"
	"strings"
	"time"
)

const spotifyMinMsPlayed = 20000

func (p *markdown) spotify() error {
	var items []spotify.StreamingHistoryItem

	p.DB().
		Model(spotify.StreamingHistoryItem{}).
		Where("ms_played > ?", spotifyMinMsPlayed).
		FindInBatches(&items, 100, func(tx *gorm.DB, batch int) error {
			for _, item := range items {
				datePage := p.pageForDate(time.Time(item.EndTime))

				artistName := strings.TrimLeft(item.ArtistName, "#")

				trackName := fmt.Sprintf("%s – %s", artistName, item.TrackName)

				datePage.Children = append(datePage.Children, Child{
					String: trackName,
					Type:   ChildSpotify,
					Time:   time.Time(item.EndTime),
				})
			}
			return nil
		})

	return nil
}
