package instagram

import (
	"github.com/bionic-dev/bionic/internal/provider/database"
	"github.com/bionic-dev/bionic/pkg/ptr"
	_ "github.com/bionic-dev/bionic/testinit"
	"github.com/bionic-dev/bionic/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"testing"
	"time"
)

func TestInstagram_importMedia(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	p := instagram{Database: database.New(db)}
	require.NoError(t, p.Migrate())

	require.NoError(t, p.importMedia("testdata/instagram/media.json"))

	var media []MediaItem
	require.NoError(t, db.
		Preload("UserMentions", func(db *gorm.DB) *gorm.DB { return db.Preload(clause.Associations) }).
		Preload("HashtagMentions", func(db *gorm.DB) *gorm.DB { return db.Preload(clause.Associations) }).
		Preload(clause.Associations).
		Find(&media).
		Error)
	require.Len(t, media, 3)

	var profilePhoto []ProfilePhoto
	require.NoError(t, db.Find(&profilePhoto).Error)
	require.Len(t, profilePhoto, 1)

	assertMedia(t, MediaItem{
		Type:    MediaStory,
		Caption: ptr.String("how do you like my #outfit? it's @asos"),
		TakenAt: types.DateTime(time.Date(2021, 1, 5, 17, 35, 26, 0, time.UTC)),
		Path:    "stories/202101/7707a9ca67f53984bf7e26eb64c41e21.jpg",
		UserMentions: []MediaUserMention{
			{
				User: User{
					Username: "asos",
				},
				FromIdx: 33,
				ToIdx:   38,
			},
		},
		HashtagMentions: []MediaHashtagMention{
			{
				Hashtag: Hashtag{
					Text: "outfit",
				},
				FromIdx: 19,
				ToIdx:   26,
			},
		},
	}, media[0])

	assertMedia(t, MediaItem{
		Type:    MediaVideo,
		Caption: ptr.String("#istanbul is beautiful! @sevazhidkov"),
		TakenAt: types.DateTime(time.Date(2020, 8, 2, 20, 35, 46, 0, time.UTC)),
		Path:    "videos/202008/8458fe5533068d764783ba3904fd2fea.mp4",
		UserMentions: []MediaUserMention{
			{
				User: User{
					Username: "sevazhidkov",
				},
				FromIdx: 24,
				ToIdx:   36,
			},
		},
		HashtagMentions: []MediaHashtagMention{
			{
				Hashtag: Hashtag{
					Text: "istanbul",
				},
				FromIdx: 0,
				ToIdx:   9,
			},
		},
	}, media[1])

	assertMedia(t, MediaItem{
		Type:     MediaPhoto,
		Caption:  ptr.String("i'm in #turkey now with @sevazhidkov"),
		TakenAt:  types.DateTime(time.Date(2020, 9, 27, 18, 37, 9, 0, time.UTC)),
		Location: ptr.String("Emirgân Korusu"),
		Path:     "photos/202009/1b31609a81d7749a1462e38ee425ea55.jpg",
		UserMentions: []MediaUserMention{
			{
				User: User{
					Username: "sevazhidkov",
				},
				FromIdx: 24,
				ToIdx:   36,
			},
		},
		HashtagMentions: []MediaHashtagMention{
			{
				Hashtag: Hashtag{
					Text: "turkey",
				},
				FromIdx: 7,
				ToIdx:   14,
			},
		},
	}, media[2])

	assert.EqualValues(t, time.Date(2020, 10, 7, 8, 57, 48, 0, time.UTC), profilePhoto[0].TakenAt)
	assert.Equal(t, true, profilePhoto[0].IsActiveProfile)
	assert.Equal(t, "profile/202010/fd47638085821857192d11805bd71898.jpg", profilePhoto[0].Path)
}

func assertMedia(t *testing.T, expected, actual MediaItem) {
	assert.Equal(t, expected.Caption, actual.Caption)
	assert.Equal(t, expected.TakenAt, actual.TakenAt)
	assert.Equal(t, expected.Location, actual.Location)
	assert.Equal(t, expected.Path, actual.Path)

	require.Equal(t, len(expected.UserMentions), len(actual.UserMentions))
	for i := range expected.UserMentions {
		assertMediaUserMention(t, expected.UserMentions[i], actual.UserMentions[i])
	}

	require.Equal(t, len(expected.HashtagMentions), len(actual.HashtagMentions))
	for i := range expected.HashtagMentions {
		assertMediaHashtagMention(t, expected.HashtagMentions[i], actual.HashtagMentions[i])
	}
}

func assertMediaUserMention(t *testing.T, expected, actual MediaUserMention) {
	assert.Equal(t, expected.User.Username, actual.User.Username)
	assert.Equal(t, expected.FromIdx, actual.FromIdx)
	assert.Equal(t, expected.ToIdx, actual.ToIdx)
}

func assertMediaHashtagMention(t *testing.T, expected, actual MediaHashtagMention) {
	assert.Equal(t, expected.Hashtag.Text, actual.Hashtag.Text)
	assert.Equal(t, expected.FromIdx, actual.FromIdx)
	assert.Equal(t, expected.ToIdx, actual.ToIdx)
}
