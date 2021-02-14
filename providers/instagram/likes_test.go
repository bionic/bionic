package instagram

import (
	"github.com/BionicTeam/bionic/providers/provider"
	_ "github.com/BionicTeam/bionic/testinit"
	"github.com/BionicTeam/bionic/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"testing"
	"time"
)

func TestInstagram_importLikes(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	p := instagram{Database: provider.NewDatabase(db)}
	require.NoError(t, p.Migrate())

	require.NoError(t, p.importLikes("testdata/instagram/likes.json"))

	var likes []Like
	require.NoError(t, db.Preload(clause.Associations).Find(&likes).Error)
	require.Len(t, likes, 4)

	assertLike(t, Like{
		Target:    LikeMedia,
		Timestamp: types.DateTime(time.Date(2021, 1, 7, 11, 41, 24, 0, time.UTC)),
		User: User{
			Username: "shekhirin",
		},
	}, likes[0])
	assertLike(t, Like{
		Target:    LikeMedia,
		Timestamp: types.DateTime(time.Date(2021, 1, 6, 17, 25, 56, 0, time.UTC)),
		User: User{
			Username: "sevazhidkov",
		},
	}, likes[1])
	assertLike(t, Like{
		Target:    LikeComment,
		Timestamp: types.DateTime(time.Date(2020, 12, 23, 14, 53, 56, 0, time.UTC)),
		User: User{
			Username: "shekhirin",
		},
	}, likes[2])
	assertLike(t, Like{
		Target:    LikeComment,
		Timestamp: types.DateTime(time.Date(2020, 12, 22, 2, 34, 13, 0, time.UTC)),
		User: User{
			Username: "lexfridman",
		},
	}, likes[3])

	assert.Equal(t, likes[0].UserID, likes[2].UserID)
}

func assertLike(t *testing.T, expected, actual Like) {
	assert.Equal(t, expected.Target, actual.Target)
	assert.Equal(t, expected.Timestamp, actual.Timestamp)
	assert.Equal(t, expected.User.Username, actual.User.Username)
}
