package instagram

import (
	"github.com/bionic-dev/bionic/internal/provider/database"
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

func TestInstagram_importStoriesActivities(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	p := instagram{Database: database.New(db)}
	require.NoError(t, p.Migrate())

	require.NoError(t, p.importStoriesActivities("testdata/instagram/stories_activities.json"))

	var storiesActivity []StoriesActivityItem
	require.NoError(t, db.Preload(clause.Associations).Find(&storiesActivity).Error)
	require.Len(t, storiesActivity, 5)

	assertStoriesActivity(t, StoriesActivityItem{
		Type: StoriesActivityPoll,
		User: User{
			Username: "shekhirin",
		},
		Timestamp: types.DateTime(time.Date(2021, 1, 3, 12, 36, 33, 0, time.UTC)),
	}, storiesActivity[0])
	assertStoriesActivity(t, StoriesActivityItem{
		Type: StoriesActivityEmojiSlider,
		User: User{
			Username: "sevazhidkov",
		},
		Timestamp: types.DateTime(time.Date(2020, 8, 24, 13, 48, 53, 0, time.UTC)),
	}, storiesActivity[1])
	assertStoriesActivity(t, StoriesActivityItem{
		Type: StoriesActivityQuestion,
		User: User{
			Username: "lexfridman",
		},
		Timestamp: types.DateTime(time.Date(2019, 4, 25, 18, 55, 34, 0, time.UTC)),
	}, storiesActivity[2])
	assertStoriesActivity(t, StoriesActivityItem{
		Type: StoriesActivityCountdown,
		User: User{
			Username: "zuck",
		},
		Timestamp: types.DateTime(time.Date(2019, 5, 6, 19, 0, 45, 0, time.UTC)),
	}, storiesActivity[3])
	assertStoriesActivity(t, StoriesActivityItem{
		Type: StoriesActivityQuiz,
		User: User{
			Username: "shekhirin",
		},
		Timestamp: types.DateTime(time.Date(2020, 12, 27, 12, 10, 23, 0, time.UTC)),
	}, storiesActivity[4])

	assert.Equal(t, storiesActivity[0].UserID, storiesActivity[4].UserID)
}

func assertStoriesActivity(t *testing.T, expected, actual StoriesActivityItem) {
	assert.Equal(t, expected.Type, actual.Type)
	assert.Equal(t, expected.User.Username, actual.User.Username)
	assert.Equal(t, expected.Timestamp, actual.Timestamp)
}
