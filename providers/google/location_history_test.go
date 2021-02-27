package google

import (
	"github.com/bionic-dev/bionic/providers/provider"
	_ "github.com/bionic-dev/bionic/testinit"
	"github.com/bionic-dev/bionic/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestGoogle_importLocationHistoryFromFile(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:?_loc=UTC"), &gorm.Config{})
	require.NoError(t, err)

	p := google{Database: provider.NewDatabase(db)}
	require.NoError(t, p.Migrate())

	require.NoError(t, p.importLocationHistoryFromFile("testdata/google/Location History/Location History.json"))

	var locationHistory []LocationHistoryItem
	require.NoError(t, db.Preload("Activities.TypeCandidates").Find(&locationHistory).Error)
	require.Len(t, locationHistory, 4)

	assertLocationHistoryItem(t, LocationHistoryItem{
		Accuracy:    65,
		Activities:  nil,
		LatitudeE7:  570799566,
		LongitudeE7: 539860481,
		Time:        types.DateTime(time.Date(2014, 11, 1, 14, 27, 19, 656000000, time.UTC)),
	}, locationHistory[0])
	assertLocationHistoryItem(t, LocationHistoryItem{
		Accuracy: 1217,
		Activities: []LocationActivity{
			{
				Time: types.DateTime(time.Date(2018, 01, 8, 3, 51, 34, 838000000, time.UTC)),
				TypeCandidates: []LocationActivityTypeCandidate{
					{Confidence: 31, Type: LocationActivityUnknown},
					{Confidence: 18, Type: LocationActivityStill},
				},
			},
		},
		LatitudeE7:  555904294,
		LongitudeE7: 372740438,
		Time:        types.DateTime(time.Date(2018, 1, 8, 3, 52, 34, 970000000, time.UTC)),
	}, locationHistory[1])
	assertLocationHistoryItem(t, LocationHistoryItem{
		Accuracy: 800,
		Activities: []LocationActivity{
			{
				Time: types.DateTime(time.Date(2018, 01, 8, 3, 53, 59, 709000000, time.UTC)),
				TypeCandidates: []LocationActivityTypeCandidate{
					{Confidence: 41, Type: LocationActivityInVehicle},
					{Confidence: 41, Type: LocationActivityInRoadVehicle},
				},
			},
			{
				Time: types.DateTime(time.Date(2018, 01, 8, 3, 54, 28, 502000000, time.UTC)),
				TypeCandidates: []LocationActivityTypeCandidate{
					{Confidence: 30, Type: LocationActivityUnknown},
					{Confidence: 24, Type: LocationActivityInVehicle},
				},
			},
		},
		LatitudeE7:  556148023,
		LongitudeE7: 372820195,
		Time:        types.DateTime(time.Date(2018, 1, 8, 3, 53, 59, 35000000, time.UTC)),
	}, locationHistory[2])
	assertLocationHistoryItem(t, LocationHistoryItem{
		Accuracy:         17,
		Velocity:         1,
		Heading:          95,
		Altitude:         127,
		VerticalAccuracy: 23,
		Activities:       nil,
		LatitudeE7:       557463742,
		LongitudeE7:      376261118,
		Time:             types.DateTime(time.Date(2021, 1, 7, 10, 33, 20, 458000000, time.UTC)),
	}, locationHistory[3])

	locationHistory = nil
	require.NoError(t, p.importLocationHistoryFromFile("testdata/google/Location History/Location History.json"))
	require.NoError(t, db.Preload("Activities.TypeCandidates").Find(&locationHistory).Error)
	require.Len(t, locationHistory, 4)
}

func assertLocationHistoryItem(t *testing.T, expected, actual LocationHistoryItem) {
	assert.Equal(t, expected.Accuracy, actual.Accuracy)
	for i, activity := range expected.Activities {
		assertLocationActivity(t, activity, actual.Activities[i])
	}
	assert.Equal(t, expected.Altitude, actual.Altitude)
	assert.Equal(t, expected.Heading, actual.Heading)
	assert.Equal(t, expected.LatitudeE7, actual.LatitudeE7)
	assert.Equal(t, expected.LongitudeE7, actual.LongitudeE7)
	assert.Equal(t, expected.Time, actual.Time)
	assert.Equal(t, expected.Velocity, actual.Velocity)
	assert.Equal(t, expected.VerticalAccuracy, actual.VerticalAccuracy)
}

func assertLocationActivity(t *testing.T, expected, actual LocationActivity) {
	assert.Equal(t, expected.Time, actual.Time)
	for i, candidate := range expected.TypeCandidates {
		assertLocationActivityTypeCandidate(t, candidate, actual.TypeCandidates[i])
	}
}

func assertLocationActivityTypeCandidate(t *testing.T, expected, actual LocationActivityTypeCandidate) {
	assert.Equal(t, expected.Confidence, actual.Confidence)
	assert.Equal(t, expected.Type, actual.Type)
}
