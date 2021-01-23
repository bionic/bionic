package twitter

import (
	"encoding/json"
	"github.com/BionicTeam/bionic/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
)

type AgeInfoRecord struct {
	gorm.Model
	Age       int
	BirthDate types.DateTime `json:"birthDate"`
}

func (AgeInfoRecord) TableName() string {
	return tablePrefix + "age_info_records"
}

func (ai *AgeInfoRecord) UnmarshalJSON(b []byte) error {
	type alias AgeInfoRecord

	var data struct {
		AgeMeta struct {
			AgeInfo struct {
				alias
				Age []string `json:"age"`
			} `json:"ageInfo"`
		} `json:"ageMeta"`
	}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	ageInfo := data.AgeMeta.AgeInfo

	*ai = AgeInfoRecord(ageInfo.alias)

	if len(ageInfo.Age) == 1 {
		age, err := strconv.Atoi(ageInfo.Age[0])
		if err != nil {
			return err
		}

		ai.Age = age
	}

	return nil
}

func (p *twitter) importAgeInfo(inputPath string) error {
	var ageInfo []AgeInfoRecord

	if err := readJSON(
		inputPath,
		"window.YTD.ageinfo.part0 = ",
		&ageInfo,
	); err != nil {
		return err
	}

	err := p.DB().
		Clauses(clause.OnConflict{
			DoNothing: true,
		}).
		Create(ageInfo).
		Error
	if err != nil {
		return err
	}

	return nil
}
