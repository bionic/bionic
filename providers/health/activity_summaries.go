package health

import (
	"encoding/xml"
	"github.com/BionicTeam/bionic/types"
	"gorm.io/gorm"
)

type ActivitySummary struct {
	gorm.Model
	Date                   types.DateTime `xml:"dateComponents,attr" gorm:"unique"`
	ActiveEnergyBurned     float64        `xml:"activeEnergyBurned,attr"`
	ActiveEnergyBurnedGoal int            `xml:"activeEnergyBurnedGoal,attr"`
	ActiveEnergyBurnedUnit string         `xml:"activeEnergyBurnedUnit,attr"`
	AppleMoveTime          int            `xml:"appleMoveTime,attr"`
	AppleMoveTimeGoal      int            `xml:"appleMoveTimeGoal,attr"`
	AppleExerciseTime      int            `xml:"appleExerciseTime,attr"`
	AppleExerciseTimeGoal  int            `xml:"appleExerciseTimeGoal,attr"`
	AppleStandHours        int            `xml:"appleStandHours,attr"`
	AppleStandHoursGoal    int            `xml:"appleStandHoursGoal,attr"`
}

func (ActivitySummary) TableName() string {
	return tablePrefix + "activity_summaries"
}

func (as ActivitySummary) Conditions() map[string]interface{} {
	return map[string]interface{}{
		"date": as.Date,
	}
}

func (p *health) parseActivitySummary(_ *DataExport, decoder *xml.Decoder, start *xml.StartElement) error {
	var activitySummary ActivitySummary

	if err := decoder.DecodeElement(&activitySummary, start); err != nil {
		return err
	}

	return p.DB().
		FirstOrCreate(&activitySummary, activitySummary.Conditions()).
		Error
}
