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
		Order("place_visit_id asc").
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
			{AccuracyMeters: 24, LatE7: 340470606, LngE7: -1185254605, Time: types.DateTime(time.Date(2020, 2, 11, 17, 54, 14, 505000000, time.UTC))},
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
	placeVisits[0].Model = gorm.Model{}
	placeVisits[0].SimplifiedRawPathPoints[0].Model = gorm.Model{}
	placeVisits[0].SimplifiedRawPathPoints[0].PlaceVisitID = 0
	placeVisits[0].SimplifiedRawPathPoints[0].PlaceVisit = PlaceVisit{}
	placeVisits[0].ChildVisits[0].Model = gorm.Model{}
	placeVisits[0].ChildVisits[0].PlaceVisitID = 0
	placeVisits[0].ChildVisits[0].PlaceVisit = nil
	placeVisits[0].ChildVisits[0].OtherCandidateLocations[0].Model = gorm.Model{}
	placeVisits[0].ChildVisits[0].OtherCandidateLocations[0].PlaceVisitID = 0
	assert.EqualValues(t, expectedPlaceVisit, placeVisits[0])

	expectedActivitySegment := ActivitySegment{
		Activities: []ActivityTypeCandidate{
			{ActivityType: "IN_PASSENGER_VEHICLE", Probability: 93.55993270874023},
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
	activitySegments[0].Model = gorm.Model{}
	activitySegments[0].Activities[0].Model = gorm.Model{}
	activitySegments[0].Activities[0].ActivitySegmentID = 0
	activitySegments[0].Activities[0].ActivitySegment = ActivitySegment{}
	activitySegments[0].SimplifiedRawPathPoints[0].Model = gorm.Model{}
	activitySegments[0].SimplifiedRawPathPoints[0].ActivitySegmentID = 0
	activitySegments[0].SimplifiedRawPathPoints[0].ActivitySegment = ActivitySegment{}
	activitySegments[0].TransitStops[0].Model = gorm.Model{}
	activitySegments[0].TransitStops[0].ActivitySegmentID = 0
	activitySegments[0].TransitStops[0].ActivitySegment = ActivitySegment{}
	activitySegments[0].Waypoints[0].Model = gorm.Model{}
	activitySegments[0].Waypoints[0].ActivitySegmentID = 0
	activitySegments[0].Waypoints[0].ActivitySegment = ActivitySegment{}
	assert.EqualValues(t, expectedActivitySegment, activitySegments[0])

	placeVisits = nil
	activitySegments = nil
	require.NoError(t, p.importSemanticLocationHistoryFromDirectory("testdata/google/Location History/Semantic Location History/"))
	require.NoError(t, db.Find(&placeVisits).Error)
	require.NoError(t, db.Find(&activitySegments).Error)
	require.Len(t, placeVisits, 2)
	require.Len(t, activitySegments, 1)
}
