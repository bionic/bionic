package health

import (
	_ "github.com/BionicTeam/bionic/testinit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"testing"
	"time"
)

var tz0300 = time.FixedZone("", 10800)
var creationDate = time.Date(2019, 01, 19, 16, 57, 15, 0, tz0300)
var startDate = time.Date(2019, 01, 19, 16, 56, 13, 0, tz0300)
var endDate = time.Date(2019, 01, 19, 16, 57, 15, 0, tz0300)

func TestHealth_Import(t *testing.T) {
	t.Run("directory", func(t *testing.T) {
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		require.NoError(t, err)

		provider := New(db)
		require.NoError(t, provider.Migrate())

		importFns, err := provider.ImportFns("testdata/health/apple_health_export")
		require.NoError(t, err)

		for _, importFn := range importFns {
			require.NoError(t, importFn.Call(), importFn.Name())
		}

		assertModels(t, db)
	})

	t.Run("archive", func(t *testing.T) {
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		require.NoError(t, err)

		provider := New(db)
		require.NoError(t, provider.Migrate())

		importFns, err := provider.ImportFns("testdata/health/export.zip")
		require.NoError(t, err)

		for _, importFn := range importFns {
			require.NoError(t, importFn.Call(), importFn.Name())
		}

		assertModels(t, db)
	})
}

func assertModels(t *testing.T, db *gorm.DB) {
	var data []DataExport
	require.NoError(t, db.Preload(clause.Associations).Find(&data).Error)
	require.Len(t, data, 1)
	assertData(t, data[0])

	var entries []Entry
	require.NoError(t, db.Preload(clause.Associations).Find(&entries).Error)
	require.Len(t, entries, 1)
	assertEntry(t, entries[0])

	var workouts []Workout
	require.NoError(t, db.
		Preload("Route", func(db *gorm.DB) *gorm.DB { return db.Preload(clause.Associations) }).
		Preload(clause.Associations).
		Find(&workouts).
		Error)
	require.Len(t, workouts, 1)
	assertWorkout(t, workouts[0])

	var activitySummaries []ActivitySummary
	require.NoError(t, db.Find(&activitySummaries).Error)
	require.Len(t, activitySummaries, 1)
	assertActivitySummary(t, activitySummaries[0])
}

func assertData(t *testing.T, data DataExport) {
	assert.Equal(t, "en_RU", data.Locale)
	assert.EqualValues(t, time.Date(2021, 01, 11, 12, 06, 40, 0, tz0300), data.ExportDate)
	assert.EqualValues(t, time.Date(2000, 07, 19, 0, 0, 0, 0, time.UTC), data.Me.DateOfBirth)
	assert.Equal(t, "HKBiologicalSexMale", data.Me.BiologicalSex)
	assert.Equal(t, "HKBloodTypeAPositive", data.Me.BloodType)
	assert.Equal(t, "HKFitzpatrickSkinTypeNotSet", data.Me.FitzpatrickSkinType)
	assert.Equal(t, "None", data.Me.CardioFitnessMedicationsUse)
}

func assertDevice(t *testing.T, device *Device) {
	require.NotNil(t, device)

	assert.Equal(t, "Apple Watch", device.Name)
	assert.Equal(t, "Apple", device.Manufacturer)
	assert.Equal(t, "Watch", device.DeviceModel)
	assert.Equal(t, "Watch3,4", device.Hardware)
	assert.Equal(t, "5.1.2", device.Software)
}

func assertMetadataEntry(t *testing.T, entry MetadataEntry) {
	assert.Equal(t, "HKMetadataKeySyncVersion", entry.Key)
	assert.Equal(t, "2", entry.Value)
}

func assertEntry(t *testing.T, entry Entry) {
	assert.Equal(t, "HKQuantityTypeIdentifierHeartRateVariabilitySDNN", entry.Type)
	assert.Equal(t, "Alexey’s Apple Watch", entry.SourceName)
	assert.Equal(t, "5.1.2", entry.SourceVersion)
	assert.Equal(t, "ms", entry.Unit)
	assert.EqualValues(t, creationDate, entry.CreationDate)
	assert.EqualValues(t, startDate, entry.StartDate)
	assert.EqualValues(t, endDate, entry.EndDate)
	assert.Equal(t, "35.7133", entry.Value)
	assertDevice(t, entry.Device)
	require.Len(t, entry.MetadataEntries, 1)
	assertMetadataEntry(t, entry.MetadataEntries[0])
	require.Len(t, entry.BeatsPerMinutes, 1)
	assert.Equal(t, 70, entry.BeatsPerMinutes[0].BPM)
	assert.Equal(t, "4:56:15,46 PM", entry.BeatsPerMinutes[0].Time)
}

func assertWorkout(t *testing.T, workout Workout) {
	assert.Equal(t, "HKWorkoutActivityTypeWalking", workout.ActivityType)
	assert.Equal(t, 16.49007770021757, workout.Duration)
	assert.Equal(t, "min", workout.DurationUnit)
	assert.Equal(t, 1.154875562449862, workout.TotalDistance)
	assert.Equal(t, "km", workout.TotalDistanceUnit)
	assert.Equal(t, 52.07101376026529, workout.TotalEnergyBurned)
	assert.Equal(t, "kcal", workout.TotalEnergyBurnedUnit)
	assert.Equal(t, "Alexey’s Apple Watch", workout.SourceName)
	assert.Equal(t, "5.1.2", workout.SourceVersion)
	assert.EqualValues(t, creationDate, workout.CreationDate)
	assert.EqualValues(t, startDate, workout.StartDate)
	assert.EqualValues(t, endDate, workout.EndDate)
	assertDevice(t, workout.Device)
	require.Len(t, workout.MetadataEntries, 1)
	assertMetadataEntry(t, workout.MetadataEntries[0])
	require.Len(t, workout.Events, 1)
	assert.Equal(t, "HKWorkoutEventTypeSegment", workout.Events[0].Type)
	assert.EqualValues(t, time.Date(2019, 01, 22, 20, 20, 16, 0, tz0300), workout.Events[0].Date)
	assert.Equal(t, 14.84274098277092, workout.Events[0].Duration)
	assert.Equal(t, "min", workout.Events[0].DurationUnit)
	assertWorkoutRoute(t, workout.Route)
}

func assertWorkoutRoute(t *testing.T, route *WorkoutRoute) {
	require.NotNil(t, route)

	assert.Equal(t, "Alexey’s Apple Watch", route.SourceName)
	assert.Equal(t, "12.1.2", route.SourceVersion)
	assert.EqualValues(t, creationDate, route.CreationDate)
	assert.EqualValues(t, startDate, route.StartDate)
	assert.EqualValues(t, endDate, route.EndDate)
	require.Len(t, route.MetadataEntries, 1)
	assertMetadataEntry(t, route.MetadataEntries[0])
	assert.Equal(t, "/workout-routes/route_2019-01-22_8.32pm.gpx", route.FilePath)
	assert.EqualValues(t, time.Date(2021, 01, 11, 9, 07, 01, 0, time.UTC), route.Time)
	assert.Equal(t, "Route 2019-01-22 8:32pm", route.TrackName)
	require.Len(t, route.TrackPoints, 1)
	assert.Equal(t, 30.359697, route.TrackPoints[0].Lon)
	assert.Equal(t, 59.92849, route.TrackPoints[0].Lat)
	assert.Equal(t, 10.150643, route.TrackPoints[0].Ele)
	assert.EqualValues(t, time.Date(2019, 01, 22, 17, 30, 26, 0, time.UTC), route.TrackPoints[0].Time)
	assert.Equal(t, 1.596127, route.TrackPoints[0].Speed)
	assert.Equal(t, 17.516218, route.TrackPoints[0].Course)
	assert.Equal(t, 4.325758, route.TrackPoints[0].HAcc)
	assert.Equal(t, 3.928791, route.TrackPoints[0].VAcc)
}

func assertActivitySummary(t *testing.T, summary ActivitySummary) {
	assert.EqualValues(t, 10, summary.ActiveEnergyBurned)
	assert.EqualValues(t, 15, summary.ActiveEnergyBurnedGoal)
	assert.EqualValues(t, "kcal", summary.ActiveEnergyBurnedUnit)
	assert.EqualValues(t, 20, summary.AppleMoveTime)
	assert.EqualValues(t, 40, summary.AppleMoveTimeGoal)
	assert.EqualValues(t, 30, summary.AppleExerciseTime)
	assert.EqualValues(t, 50, summary.AppleExerciseTimeGoal)
	assert.EqualValues(t, 90, summary.AppleStandHours)
	assert.EqualValues(t, 12, summary.AppleStandHoursGoal)
}
