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

func TestGoogle_importSemanticLocationHistoryFromDirectory(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:?_loc=UTC"), &gorm.Config{})
	require.NoError(t, err)

	p := google{Database: provider.NewDatabase(db)}
	require.NoError(t, p.Migrate())

	require.NoError(t, p.importSemanticLocationHistoryFromDirectory("testdata/google/Location History/Semantic Location History/"))

	var placeVisits []PlaceVisit
	var activitySegments []ActivitySegment
	require.NoError(t, db.
		Preload("OtherCandidateLocations").
		Preload("SimplifiedRawPathPoints").
		Preload("ChildVisits.OtherCandidateLocations").
		Order("id asc").
		Find(&placeVisits).
		Error)
	require.NoError(t, db.
		Preload("Activities").
		Preload("SimplifiedRawPathPoints").
		Preload("TransitStops").
		Preload("Waypoints").
		Find(&activitySegments).
		Error)

	require.Len(t, placeVisits, 2)
	require.Len(t, activitySegments, 1)

	expectedPlaceVisit := PlaceVisit{
		CenterLatE7:                 340154991,
		CenterLngE7:                 -1184860798,
		DurationEndTimestampMs:      types.DateTime(time.Date(2020, 2, 4, 18, 12, 44, 675000000, time.UTC)),
		DurationStartTimestampMs:    types.DateTime(time.Date(2020, 2, 4, 17, 39, 12, 182000000, time.UTC)),
		EditConfirmationStatus:      "NOT_CONFIRMED",
		LocationAddress:             "1351 3rd Street Promenade\nSanta Monica, CA 90401\nUSA",
		LocationLatitudeE7:          340165480,
		LocationLocationConfidence:  55.464928,
		LocationLongitudeE7:         -1184961984,
		LocationName:                "Downtown Santa Monica",
		LocationPlaceID:             "ChIJFwxnxM-kwoARewIypbH6UsU",
		LocationSourceInfoDeviceTag: 59293564,
		PlaceConfidence:             "MEDIUM_CONFIDENCE",
		VisitConfidence:             75,
		SimplifiedRawPathPoints: []PlacePathPoint{
			{
				AccuracyMeters: 24, LatE7: 340470606, LngE7: -1185254605,
				Time: types.DateTime(time.Date(2020, 2, 11, 17, 54, 14, 505000000, time.UTC)),
			},
		},
		OtherCandidateLocations: []CandidateLocation{},
		ChildVisits: []*PlaceVisit{
			{
				CenterLatE7:                 340154991,
				CenterLngE7:                 -1184860798,
				DurationEndTimestampMs:      types.DateTime(time.Date(2020, 2, 4, 18, 12, 44, 675000000, time.UTC)),
				DurationStartTimestampMs:    types.DateTime(time.Date(2020, 2, 4, 17, 39, 12, 182000000, time.UTC)),
				EditConfirmationStatus:      "NOT_CONFIRMED",
				LocationAddress:             "1654 Lincoln Blvd\nSanta Monica, CA 90404\nUSA",
				LocationLatitudeE7:          340157520,
				LocationLocationConfidence:  25.894972,
				LocationLongitudeE7:         -1184872646,
				LocationName:                "Box 'N Burn Santa monica Boxing Fitness gym",
				LocationPlaceID:             "ChIJD0XpyNKkwoARDAiHsXtjcvg",
				LocationSourceInfoDeviceTag: 59293564,
				PlaceConfidence:             "LOW_CONFIDENCE",
				VisitConfidence:             75,
				OtherCandidateLocations: []CandidateLocation{
					{LatitudeE7: 340154424, LocationConfidence: 11.792331, LongitudeE7: -1184863484, PlaceID: "ChIJj0JGzdKkwoARw0F60V1_eQE"},
				},
			},
		},
	}
	assertPlaceVisit(t, expectedPlaceVisit, placeVisits[0])

	expectedActivitySegment := ActivitySegment{
		Activities: []ActivityTypeCandidate{
			{
				ActivityType: "IN_PASSENGER_VEHICLE",
				Probability:  93.55993270874023,
			},
		},
		ActivityType:                       "IN_PASSENGER_VEHICLE",
		Confidence:                         "HIGH",
		Distance:                           13958,
		DurationEndTimestampMs:             types.DateTime(time.Date(2020, 2, 1, 20, 1, 43, 307000000, time.UTC)),
		DurationStartTimestampMs:           types.DateTime(time.Date(2020, 2, 1, 19, 22, 14, 623000000, time.UTC)),
		EndLocationLatitudeE7:              340322573,
		EndLocationLongitudeE7:             -1184817511,
		ParkingEventLocationAccuracyMetres: 125,
		ParkingEventLocationLatitudeE7:     340320032,
		ParkingEventLocationLongitudeE7:    -1184827155,
		ParkingEventTimestampMs:            types.DateTime(time.Date(2020, 2, 1, 20, 4, 20, 47000000, time.UTC)),
		SimplifiedRawPathPoints: []ActivityPathPoint{
			{
				AccuracyMeters: 65, LatE7: 340779839, LngE7: -1184296112,
				Time: types.DateTime(time.Date(2020, 2, 1, 19, 26, 45, 571000000, time.UTC)),
			},
		},
		StartLocationLatitudeE7:  340902329,
		StartLocationLongitudeE7: -1184383023,
		TransitPathHexRgbColor:   "ED7D31",
		TransitPathName:          "33",
		TransitStops: []TransitStop{
			{
				LatitudeE7: 339888109, LongitudeE7: -1184626624,
				PlaceID: "ChIJZYma8pW6woARV34ul5SMRSk",
				Name:    "Venice / Abbot Kinney",
			},
		},
		Waypoints: []Waypoint{
			{LatE7: 340902481, LngE7: -1184381866},
		},
	}
	assertActivitySegments(t, expectedActivitySegment, activitySegments[0])

	var newPlaceVisits []PlaceVisit
	var newActivitySegments []ActivitySegment
	require.NoError(t, p.importSemanticLocationHistoryFromDirectory("testdata/google/Location History/Semantic Location History/"))
	require.NoError(t, db.
		Preload("OtherCandidateLocations").
		Preload("SimplifiedRawPathPoints").
		Preload("ChildVisits.OtherCandidateLocations").
		Order("id asc").
		Find(&newPlaceVisits).
		Error)
	require.NoError(t, db.
		Preload("Activities").
		Preload("SimplifiedRawPathPoints").
		Preload("TransitStops").
		Preload("Waypoints").
		Find(&newActivitySegments).
		Error)
	require.Len(t, newPlaceVisits, 2)
	assertPlaceVisit(t, placeVisits[0], newPlaceVisits[0])
	require.Len(t, newActivitySegments, 1)
	assertActivitySegments(t, activitySegments[0], newActivitySegments[0])
}

func assertPlaceVisit(t *testing.T, expected, actual PlaceVisit) {
	expected = convertPlaceVisit(expected)
	actual = convertPlaceVisit(actual)
	assert.EqualValues(t, expected, actual)
}

func convertPlaceVisit(visit PlaceVisit) PlaceVisit {
	visit.Model = gorm.Model{}
	visit.PlaceVisitID = 0
	for i := range visit.SimplifiedRawPathPoints {
		visit.SimplifiedRawPathPoints[i].Model = gorm.Model{}
		visit.SimplifiedRawPathPoints[i].PlaceVisitID = 0
		visit.SimplifiedRawPathPoints[i].PlaceVisit = PlaceVisit{}
	}
	for i := range visit.OtherCandidateLocations {
		visit.OtherCandidateLocations[i].Model = gorm.Model{}
		visit.OtherCandidateLocations[i].PlaceVisitID = 0
	}
	for i, childVisit := range visit.ChildVisits {
		newChildVisit := convertPlaceVisit(*childVisit)
		visit.ChildVisits[i] = &newChildVisit
	}

	return visit
}

func assertActivitySegments(t *testing.T, expected, actual ActivitySegment) {
	expected = convertActivitySegment(expected)
	actual = convertActivitySegment(actual)
	assert.EqualValues(t, expected, actual)
}

func convertActivitySegment(segment ActivitySegment) ActivitySegment {
	segment.Model = gorm.Model{}
	for i := range segment.Activities {
		segment.Activities[i].Model = gorm.Model{}
		segment.Activities[i].ActivitySegmentID = 0
		segment.Activities[i].ActivitySegment = ActivitySegment{}
	}
	for i := range segment.SimplifiedRawPathPoints {
		segment.SimplifiedRawPathPoints[i].Model = gorm.Model{}
		segment.SimplifiedRawPathPoints[i].ActivitySegmentID = 0
		segment.SimplifiedRawPathPoints[i].ActivitySegment = ActivitySegment{}
	}
	for i := range segment.TransitStops {
		segment.TransitStops[i].Model = gorm.Model{}
		segment.TransitStops[i].ActivitySegmentID = 0
		segment.TransitStops[i].ActivitySegment = ActivitySegment{}
	}
	for i := range segment.Waypoints {
		segment.Waypoints[i].Model = gorm.Model{}
		segment.Waypoints[i].ActivitySegmentID = 0
		segment.Waypoints[i].ActivitySegment = ActivitySegment{}
	}
	return segment
}
