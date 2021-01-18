package netflix

import (
	"github.com/gocarina/gocsv"
	"github.com/shekhirin/bionic-cli/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
)

type MyListItem struct {
	gorm.Model
	ProfileName  string         `csv:"Profile Name" gorm:"uniqueIndex:netflix_my_list_key"`
	TitleName    string         `csv:"Title Name" gorm:"uniqueIndex:netflix_my_list_key"`
	Country      string         `csv:"Country"`
	TitleAddDate types.DateTime `csv:"Utc Title Add Date"`
}

func (r MyListItem) TableName() string {
	return "netflix_my_list"
}

func (p *netflix) importMyList(inputPath string) error {
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return nil
	}

	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	var items []MyListItem

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
