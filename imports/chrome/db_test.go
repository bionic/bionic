package chrome

import (
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

func TestChrome_importDB(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:?_loc=UTC"), &gorm.Config{})
	require.NoError(t, err)

	p := chrome{Database: database.New(db)}
	require.NoError(t, p.Migrate())

	require.NoError(t, p.importDB("testdata/chrome/db.sqlite?_loc=UTC"))

	var visits []Visit
	require.NoError(t, db.Preload("Segment.URL").Preload("URL").Find(&visits).Error)
	require.Len(t, visits, 1)

	expectedURL := URL{
		URL:        "https://mercury.com/",
		Title:      "Mercury | Banking built for startups",
		VisitCount: 46,
		TypedCount: 42,
		LastVisit:  types.DateTime(time.Date(2021, 4, 25, 17, 51, 12, 0, time.UTC)),
		Hidden:     false,
	}

	assertVisit(t, Visit{
		URL:                     expectedURL,
		Time:                    types.DateTime(time.Date(2021, 1, 26, 6, 19, 44, 0, time.UTC)),
		VisitID:                 0,
		Visit:                   nil,
		TransitionType:          "TYPED",
		TransitionQualifierType: "",
		IsRedirect:              false,
		Segment: Segment{
			Name: "http://mercury.com/",
			URL:  expectedURL,
		},
		VisitDuration:                0,
		IncrementedOmniboxTypedScore: true,
		PubliclyRoutable:             true,
	}, visits[0])
}

func assertVisit(t *testing.T, expected, actual Visit) {
	expected = cleanVisit(expected)
	actual = cleanVisit(actual)
	assert.EqualValues(t, expected, actual)
}

func cleanVisit(v Visit) Visit {
	v.Model = gorm.Model{}
	v.URLID = 0
	v.URL.Model = gorm.Model{}
	v.SegmentID = 0
	v.Segment.Model = gorm.Model{}
	v.Segment.URLID = 0
	v.Segment.URL.Model = gorm.Model{}
	return v
}
