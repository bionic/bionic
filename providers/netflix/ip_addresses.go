package netflix

import (
	"github.com/gocarina/gocsv"
	"github.com/shekhirin/bionic-cli/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
)

type IPAddress struct {
	gorm.Model
	ESN                        string         `csv:"esn" gorm:"uniqueIndex:idx_ip_addresses"`
	Country                    string         `csv:"country"`
	LocalizedDeviceDescription string         `csv:"localizedDeviceDescription"`
	DeviceDescription          string         `csv:"deviceDescription"`
	IP                         string         `csv:"ip"`
	RegionCodeDisplayName      string         `csv:"regionCodeDisplayName"`
	Time                       types.DateTime `csv:"ts" gorm:"uniqueIndex:idx_ip_addresses"`
}

func (r IPAddress) TableName() string {
	return "netflix_ip_addresses"
}

func (p *netflix) importIPAddresses(inputPath string) error {
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return nil
	}

	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	var items []IPAddress

	if err := gocsv.UnmarshalFile(file, &items); err != nil {
		return err
	}

	err = p.DB().
		Clauses(clause.OnConflict{
			DoNothing: true,
		}).
		CreateInBatches(items, 100).
		Error
	if err != nil {
		return err
	}

	return nil
}
