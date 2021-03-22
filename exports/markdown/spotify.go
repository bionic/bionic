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
		Where("ms_played > ?", spotifyMinMsPlayed).
		FindInBatches(&items, 100, func(tx *gorm.DB, batch int) error {
			for _, item := range items {
				localTime := time.Time(item.EndTime).Local()

				datePage := p.pageForDate(localTime)

				artistName := strings.TrimLeft(item.ArtistName, "#")

				trackName := fmt.Sprintf(
					"[[%s]] – %s for %s",
					artistName,
					item.TrackName,
					formatDuration(time.Millisecond*time.Duration(item.MsPlayed), time.Second),
				)

				datePage.Children = append(datePage.Children, Child{
					String: trackName,
					Type:   ChildSpotify,
					Time:   localTime,
				})
			}
			return nil
		})

	return nil
}
