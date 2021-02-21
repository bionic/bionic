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

func TestInstagram_importComments(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	p := instagram{Database: provider.NewDatabase(db)}
	require.NoError(t, p.Migrate())

	require.NoError(t, p.importComments("testdata/instagram/comments.json"))

	var comments []Comment
	require.NoError(t, db.
		Preload("UserMentions", func(db *gorm.DB) *gorm.DB { return db.Preload(clause.Associations) }).
		Preload("HashtagMentions", func(db *gorm.DB) *gorm.DB { return db.Preload(clause.Associations) }).
		Preload(clause.Associations).
		Find(&comments).
		Error)
	require.Len(t, comments, 1)

	assertComment(t, Comment{
		Target: CommentMedia,
		User: User{
			Username: "shekhirin",
		},
		Text: "@shekhirin nice #look dude. You look almost like @sevazhidkov",
		UserMentions: []CommentUserMention{
			{
				User: User{
					Username: "shekhirin",
				},
				FromIdx: 0,
				ToIdx:   10,
			},
			{
				User: User{
					Username: "sevazhidkov",
				},
				FromIdx: 49,
				ToIdx:   61,
			},
		},
		HashtagMentions: []CommentHashtagMention{
			{
				Hashtag: Hashtag{
					Text: "look",
				},
				FromIdx: 16,
				ToIdx:   21,
			},
		},
		Timestamp: types.DateTime(time.Date(2021, 1, 8, 9, 12, 6, 0, time.UTC)),
	}, comments[0])

	assert.Equal(t, comments[0].UserID, comments[0].UserMentions[0].UserID)
}

func assertComment(t *testing.T, expected, actual Comment) {
	assert.Equal(t, expected.Target, actual.Target)
	assert.Equal(t, expected.User.Username, actual.User.Username)

	require.Equal(t, len(expected.UserMentions), len(actual.UserMentions))
	for i := range expected.UserMentions {
		assertCommentUserMention(t, expected.UserMentions[i], actual.UserMentions[i])
	}

	require.Equal(t, len(expected.HashtagMentions), len(actual.HashtagMentions))
	for i := range expected.HashtagMentions {
		assertCommentHashtagMention(t, expected.HashtagMentions[i], actual.HashtagMentions[i])
	}

	assert.Equal(t, expected.Timestamp, actual.Timestamp)
}

func assertCommentUserMention(t *testing.T, expected, actual CommentUserMention) {
	assert.Equal(t, expected.User.Username, actual.User.Username)
	assert.Equal(t, expected.FromIdx, actual.FromIdx)
	assert.Equal(t, expected.ToIdx, actual.ToIdx)
}

func assertCommentHashtagMention(t *testing.T, expected, actual CommentHashtagMention) {
	assert.Equal(t, expected.Hashtag.Text, actual.Hashtag.Text)
	assert.Equal(t, expected.FromIdx, actual.FromIdx)
	assert.Equal(t, expected.ToIdx, actual.ToIdx)
}
