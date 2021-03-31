package rescuetime

import (
	"database/sql"
	"github.com/bionic-dev/bionic/internal/provider/database"
	_ "github.com/bionic-dev/bionic/testinit"
	"github.com/bionic-dev/bionic/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestInstagram_importAccountHistory(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	p := rescuetime{Database: database.New(db)}
	require.NoError(t, p.Migrate())

	require.NoError(t, p.importActivityHistory("testdata/rescuetime/rescuetime-activity-history.csv"))

	var activityHistory []ActivityHistoryItem
	require.NoError(t, db.Find(&activityHistory).Error)
	require.Len(t, activityHistory, 2)

	assertActivityHistoryItem(t, ActivityHistoryItem{
		Activity:  "telegram",
		Details:   ActivityHistoryDetails{sql.NullString{Valid: false}},
		Category:  "Communication & Scheduling",
		Class:     "Instant Message",
		Duration:  526,
		Timestamp: types.DateTime(time.Date(2017, 9, 22, 13, 0, 0, 0, time.FixedZone("", -25200))),
	}, activityHistory[0])
	assertActivityHistoryItem(t, ActivityHistoryItem{
		Activity:  "tripadvisor.com",
		Details:   ActivityHistoryDetails{sql.NullString{String: "The Top 10 Things to Do in Budapest 2017", Valid: true}},
		Category:  "Reference & Learning",
		Class:     "Travel & Outdoors",
		Duration:  33,
		Timestamp: types.DateTime(time.Date(2017, 9, 22, 13, 0, 0, 0, time.FixedZone("", -25200))),
	}, activityHistory[1])

	require.NoError(t, p.importActivityHistory("testdata/rescuetime/rescuetime-activity-history.csv"))

	var newActivityHistory []ActivityHistoryItem
	require.NoError(t, db.Find(&newActivityHistory).Error)
	assertActivityHistories(t, activityHistory, newActivityHistory)
}

func assertActivityHistoryItem(t *testing.T, expected, actual ActivityHistoryItem) {
	expected.Model = gorm.Model{}
	actual.Model = gorm.Model{}

	assert.EqualValues(t, expected, actual)
}

func assertActivityHistories(t *testing.T, expected, actual []ActivityHistoryItem) {
	require.Equal(t, len(expected), len(actual))

	for i := range expected {
		assertActivityHistoryItem(t, expected[i], actual[i])
	}
}
