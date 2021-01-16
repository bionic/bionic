package netflix

import (
	"github.com/gocarina/gocsv"
	"github.com/shekhirin/bionic-cli/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
)

type Rating struct {
	gorm.Model
	ProfileName    string         `csv:"Profile Name" gorm:"uniqueIndex:idx_rating"`
	TitleName      string         `csv:"Title Name" gorm:"uniqueIndex:idx_rating"`
	RatingType     string         `csv:"Rating Type"`
	StarValue      int            `csv:"Star Value"`
	ThumbsValue    int            `csv:"Thumbs Value"`
	DeviceModel    string         `csv:"Device Model"`
	EventTime      types.DateTime `csv:"Event Utc Ts" gorm:"uniqueIndex:idx_rating"`
	RegionViewDate types.DateTime `csv:"Region View Date"`
}

func (r Rating) TableName() string {
	return "netflix_ratings"
}

func (p *netflix) importRatings(inputPath string) error {
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return nil
	}

	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	var ratings []Rating

	if err := gocsv.UnmarshalFile(file, &ratings); err != nil {
		return err
	}

	err = p.DB().
		Clauses(clause.OnConflict{
			DoNothing: true,
		}).
		CreateInBatches(ratings, 100).
		Error
	if err != nil {
		return err
	}

	return nil
}
