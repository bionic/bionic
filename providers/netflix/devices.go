package netflix

import (
	"github.com/BionicTeam/bionic/types"
	"github.com/gocarina/gocsv"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
)

type Device struct {
	gorm.Model
	ProfileName                                   string         `csv:"profileName" gorm:"uniqueIndex:netflix_devices_key"`
	ESN                                           string         `csv:"esn" gorm:"uniqueIndex:netflix_devices_key"`
	DeviceType                                    string         `csv:"deviceType"`
	AcctFirstPlaybackDate                         types.DateTime `csv:"acctFirstPlaybackDate"`
	AcctLastPlaybackDate                          types.DateTime `csv:"acctLastPlaybackDate"`
	AcctFirstPlaybackDateForUserGeneratedPlays    types.DateTime `csv:"acctFirstPlaybackDateForUserGeneratedPlays"`
	AcctLastPlaybackDateForUserGeneratedPlays     types.DateTime `csv:"acctLastPlaybackDateForUserGeneratedPlays"`
	ProfileFirstPlaybackDate                      types.DateTime `csv:"profileFirstPlaybackDate"`
	ProfileLastPlaybackDate                       types.DateTime `csv:"profileLastPlaybackDate"`
	ProfileFirstPlaybackDateForUserGeneratedPlays types.DateTime `csv:"profileFirstPlaybackDateForUserGeneratedPlays"`
	ProfileLastPlaybackDateForUserGeneratedPlays  types.DateTime `csv:"profileLastPlaybackDateForUserGeneratedPlays"`
	DeactivationTime                              types.DateTime `csv:"deactivationTime"`
}

func (r Device) TableName() string {
	return tablePrefix + "devices"
}

func (p *netflix) importDevices(inputPath string) error {
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return nil
	}

	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	var devices []Device

	if err := gocsv.UnmarshalFile(file, &devices); err != nil {
		return err
	}

	err = p.DB().
		Clauses(clause.OnConflict{
			DoNothing: true,
		}).
		CreateInBatches(devices, 100).
		Error
	if err != nil {
		return err
	}

	return nil
}
