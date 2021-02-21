package netflix

import (
	"github.com/bionic-dev/bionic/types"
	"github.com/gocarina/gocsv"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
)

type IPAddress struct {
	gorm.Model
	ESN                        string         `csv:"esn" gorm:"uniqueIndex:netflix_ip_addresses_key"`
	Country                    string         `csv:"country"`
	LocalizedDeviceDescription string         `csv:"localizedDeviceDescription"`
	DeviceDescription          string         `csv:"deviceDescription"`
	IP                         string         `csv:"ip"`
	RegionCodeDisplayName      string         `csv:"regionCodeDisplayName"`
	Time                       types.DateTime `csv:"ts" gorm:"uniqueIndex:netflix_ip_addresses_key"`
}

func (r IPAddress) TableName() string {
	return tablePrefix + "ip_addresses"
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
