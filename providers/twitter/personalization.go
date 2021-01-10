package twitter

import (
	"encoding/json"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io/ioutil"
	"strings"
)

type PersonalizationRecord struct {
	gorm.Model
	Languages              []PersonalizationLanguageRecord
	GenderInfo             PersonalizationGenderInfoRecord
	Interests              []PersonalizationInterestRecord
	AudienceAndAdvertisers PersonalizationAudienceAndAdvertiserRecord
	Shows                  []PersonalizationShowRecord
	LocationHistory        []PersonalizationLocationRecord
	InferredAgeInfo        PersonalizationInferredAgeInfoRecord `json:"inferredAgeInfo"`
}

func (PersonalizationRecord) TableName() string {
	return "twitter_personalization_records"
}

func (p *PersonalizationRecord) UnmarshalJSON(b []byte) error {
	type alias PersonalizationRecord

	var data struct {
		alias
		Demographics struct {
			Languages  []PersonalizationLanguageRecord `json:"languages"`
			GenderInfo PersonalizationGenderInfoRecord `json:"genderInfo"`
		} `json:"demographics"`
		Interests struct {
			Interests              []PersonalizationInterestRecord            `json:"interests"`
			AudienceAndAdvertisers PersonalizationAudienceAndAdvertiserRecord `json:"audienceAndAdvertisers"`
			Shows                  []string                                   `json:"shows"`
		} `json:"interests"`
		LocationHistory []string `json:"locationHistory"`
	}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	*p = PersonalizationRecord(data.alias)

	p.Languages = data.Demographics.Languages
	p.GenderInfo = data.Demographics.GenderInfo

	p.Interests = data.Interests.Interests
	p.AudienceAndAdvertisers = data.Interests.AudienceAndAdvertisers

	for _, show := range data.Interests.Shows {
		p.Shows = append(p.Shows, PersonalizationShowRecord{
			Name: show,
		})
	}

	for _, location := range data.LocationHistory {
		p.LocationHistory = append(p.LocationHistory, PersonalizationLocationRecord{
			Name: location,
		})
	}

	return nil
}

type PersonalizationLanguageRecord struct {
	gorm.Model
	PersonalizationRecordID int
	Language                string `json:"language"`
	IsDisabled              bool   `json:"isDisabled"`
}

func (PersonalizationLanguageRecord) TableName() string {
	return "twitter_personalization_language_records"
}

type PersonalizationGenderInfoRecord struct {
	gorm.Model
	PersonalizationRecordID int
	Gender                  string `json:"gender"`
}

func (PersonalizationGenderInfoRecord) TableName() string {
	return "twitter_personalization_gender_info_records"
}

type PersonalizationInterestRecord struct {
	gorm.Model
	PersonalizationRecordID int
	Name                    string `json:"name"`
	IsDisabled              bool   `json:"isDisabled"`
}

func (PersonalizationInterestRecord) TableName() string {
	return "twitter_personalization_interest_records"
}

type PersonalizationAudienceAndAdvertiserRecord struct {
	gorm.Model
	PersonalizationRecordID int
	NumAudiences            int                               `json:"numAudiences,string"`
	Advertisers             []PersonalizationAdvertiserRecord `gorm:"foreignKey:AudienceAndAdvertisersID"`
	LookalikeAdvertisers    []PersonalizationAdvertiserRecord `gorm:"foreignKey:AudienceAndAdvertisersID"`
}

func (PersonalizationAudienceAndAdvertiserRecord) TableName() string {
	return "twitter_personalization_audience_and_advertiser_records"
}

func (aaa *PersonalizationAudienceAndAdvertiserRecord) UnmarshalJSON(b []byte) error {
	type alias PersonalizationAudienceAndAdvertiserRecord

	var data struct {
		alias
		Advertisers          []string `json:"advertisers"`
		LookalikeAdvertisers []string `json:"lookalikeAdvertisers"`
	}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	*aaa = PersonalizationAudienceAndAdvertiserRecord(data.alias)

	for _, advertiser := range data.Advertisers {
		aaa.Advertisers = append(
			aaa.Advertisers,
			PersonalizationAdvertiserRecord{
				Name: advertiser,
			},
		)
	}

	for _, lookalikeAdvertiser := range data.LookalikeAdvertisers {
		aaa.LookalikeAdvertisers = append(
			aaa.LookalikeAdvertisers,
			PersonalizationAdvertiserRecord{
				Name:      lookalikeAdvertiser,
				Lookalike: true,
			},
		)
	}

	return nil
}

type PersonalizationAdvertiserRecord struct {
	gorm.Model
	AudienceAndAdvertisersID int
	Name                     string
	Lookalike                bool
}

func (PersonalizationAdvertiserRecord) TableName() string {
	return "twitter_personalization_advertiser_records"
}

type PersonalizationShowRecord struct {
	gorm.Model
	PersonalizationRecordID int
	Name                    string
}

func (PersonalizationShowRecord) TableName() string {
	return "twitter_personalization_show_records"
}

type PersonalizationLocationRecord struct {
	gorm.Model
	PersonalizationRecordID int
	Name                    string `json:"name"`
}

func (PersonalizationLocationRecord) TableName() string {
	return "twitter_personalization_location_records"
}

type PersonalizationInferredAgeInfoRecord struct {
	gorm.Model
	PersonalizationRecordID int
	Age                     []string `json:"age" gorm:"type:text"`
	BirthDate               string   `json:"birthDate"`
}

func (PersonalizationInferredAgeInfoRecord) TableName() string {
	return "twitter_personalization_inferred_age_info_records"
}

func (p *twitter) importPersonalization(inputPath string) error {
	var fileData []struct {
		P13nData PersonalizationRecord `json:"p13nData"`
	}

	bytes, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return err
	}

	data := string(bytes)
	data = strings.TrimPrefix(data, "window.YTD.personalization.part0 = ")

	if err := json.Unmarshal([]byte(data), &fileData); err != nil {
		return err
	}

	personalization := fileData[0].P13nData

	err = p.DB().
		Clauses(clause.OnConflict{
			DoNothing: true,
		}).
		Create(&personalization).
		Error
	if err != nil {
		return err
	}

	return nil
}
