package instagram

import (
	"github.com/bionic-dev/bionic/providers/provider"
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

	p := instagram{Database: provider.NewDatabase(db)}
	require.NoError(t, p.Migrate())

	require.NoError(t, p.importStoriesActivities("testdata/instagram/stories_activities.json"))

	var storyActivity []StoryActivityItem
	require.NoError(t, db.Preload(clause.Associations).Find(&storyActivity).Error)
	require.Len(t, storyActivity, 5)

	assertStoryActivity(t, StoryActivityItem{
		Type: StoryActivityPoll,
		User: User{
			Username: "shekhirin",
		},
		Timestamp: types.DateTime(time.Date(2021, 1, 3, 12, 36, 33, 0, time.UTC)),
	}, storyActivity[0])
	assertStoryActivity(t, StoryActivityItem{
		Type: StoryActivityEmojiSlider,
		User: User{
			Username: "sevazhidkov",
		},
		Timestamp: types.DateTime(time.Date(2020, 8, 24, 13, 48, 53, 0, time.UTC)),
	}, storyActivity[1])
	assertStoryActivity(t, StoryActivityItem{
		Type: StoryActivityQuestion,
		User: User{
			Username: "lexfridman",
		},
		Timestamp: types.DateTime(time.Date(2019, 4, 25, 18, 55, 34, 0, time.UTC)),
	}, storyActivity[2])
	assertStoryActivity(t, StoryActivityItem{
		Type: StoryActivityCountdown,
		User: User{
			Username: "zuck",
		},
		Timestamp: types.DateTime(time.Date(2019, 5, 6, 19, 0, 45, 0, time.UTC)),
	}, storyActivity[3])
	assertStoryActivity(t, StoryActivityItem{
		Type: StoryActivityQuiz,
		User: User{
			Username: "shekhirin",
		},
		Timestamp: types.DateTime(time.Date(2020, 12, 27, 12, 10, 23, 0, time.UTC)),
	}, storyActivity[4])

	assert.Equal(t, storyActivity[0].UserID, storyActivity[4].UserID)
}

func assertStoryActivity(t *testing.T, expected, actual StoryActivityItem) {
	assert.Equal(t, expected.Type, actual.Type)
	assert.Equal(t, expected.User.Username, actual.User.Username)
	assert.Equal(t, expected.Timestamp, actual.Timestamp)
}
