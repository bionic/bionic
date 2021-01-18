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
	Languages              []LanguageRecord
	GenderInfoID           int
	GenderInfo             GenderInfo
	Interests              []InterestRecord
	AudienceAndAdvertisers AudienceAndAdvertiserRecord
	Shows                  []Show                `gorm:"many2many:twitter_personalization_shows"`
	LocationHistory        []Location            `gorm:"many2many:twitter_personalization_locations"`
	InferredAgeInfo        InferredAgeInfoRecord `json:"inferredAgeInfo"`
}

func (PersonalizationRecord) TableName() string {
	return tablePrefix + "personalization_records"
}

func (p *PersonalizationRecord) UnmarshalJSON(b []byte) error {
	type alias PersonalizationRecord

	var data struct {
		alias
		Demographics struct {
			Languages  []LanguageRecord `json:"languages"`
			GenderInfo GenderInfo       `json:"genderInfo"`
		} `json:"demographics"`
		Interests struct {
			Interests              []InterestRecord            `json:"interests"`
			AudienceAndAdvertisers AudienceAndAdvertiserRecord `json:"audienceAndAdvertisers"`
			Shows                  []string                    `json:"shows"`
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
		p.Shows = append(p.Shows, Show{
			Name: show,
		})
	}

	for _, location := range data.LocationHistory {
		p.LocationHistory = append(p.LocationHistory, Location{
			Name: location,
		})
	}

	return nil
}

type LanguageRecord struct {
	gorm.Model
	PersonalizationRecordID int
	Language                string `json:"language"`
	IsDisabled              bool   `json:"isDisabled"`
}

func (LanguageRecord) TableName() string {
	return tablePrefix + "language_records"
}

type GenderInfo struct {
	gorm.Model
	Gender string `json:"gender" gorm:"unique"`
}

func (GenderInfo) TableName() string {
	return tablePrefix + "gender_info"
}

type InterestRecord struct {
	gorm.Model
	PersonalizationRecordID int
	Name                    string `json:"name"`
	IsDisabled              bool   `json:"isDisabled"`
}

func (InterestRecord) TableName() string {
	return tablePrefix + "interest_records"
}

type AudienceAndAdvertiserRecord struct {
	gorm.Model
	PersonalizationRecordID int
	NumAudiences            int          `json:"numAudiences,string"`
	Advertisers             []Advertiser `gorm:"many2many:twitter_audience_and_advertisers"`
	LookalikeAdvertisers    []Advertiser `gorm:"many2many:twitter_audience_and_lookalike_advertisers"`
}

func (AudienceAndAdvertiserRecord) TableName() string {
	return tablePrefix + "audience_and_advertiser_records"
}

func (aaa *AudienceAndAdvertiserRecord) UnmarshalJSON(b []byte) error {
	type alias AudienceAndAdvertiserRecord

	var data struct {
		alias
		Advertisers          []string `json:"advertisers"`
		LookalikeAdvertisers []string `json:"lookalikeAdvertisers"`
	}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	*aaa = AudienceAndAdvertiserRecord(data.alias)

	for _, advertiser := range data.Advertisers {
		aaa.Advertisers = append(
			aaa.Advertisers,
			Advertiser{
				Name: advertiser,
			},
		)
	}

	for _, lookalikeAdvertiser := range data.LookalikeAdvertisers {
		aaa.LookalikeAdvertisers = append(
			aaa.LookalikeAdvertisers,
			Advertiser{
				Name:      lookalikeAdvertiser,
				Lookalike: true,
			},
		)
	}

	return nil
}

type Advertiser struct {
	gorm.Model
	Name      string `gorm:"uniqueIndex:twitter_advertisers_key"`
	Lookalike bool   `gorm:"uniqueIndex:twitter_advertisers_key"`
}

func (Advertiser) TableName() string {
	return tablePrefix + "advertisers"
}

type Show struct {
	gorm.Model
	Name string `gorm:"unique"`
}

func (Show) TableName() string {
	return tablePrefix + "shows"
}

type Location struct {
	gorm.Model
	Name string `json:"name" gorm:"unique"`
}

func (Location) TableName() string {
	return tablePrefix + "locations"
}

type InferredAgeInfoRecord struct {
	gorm.Model
	PersonalizationRecordID int
	Age                     []string `json:"age" gorm:"type:text"`
	BirthDate               string   `json:"birthDate"`
}

func (InferredAgeInfoRecord) TableName() string {
	return tablePrefix + "inferred_age_info_records"
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
		FirstOrCreate(&personalization.GenderInfo, map[string]interface{}{
			"gender": personalization.GenderInfo.Gender,
		}).
		Error
	if err != nil {
		return err
	}

	var advertisers []*Advertiser
	for i := range personalization.AudienceAndAdvertisers.Advertisers {
		advertisers = append(advertisers, &personalization.AudienceAndAdvertisers.Advertisers[i])
	}
	for i := range personalization.AudienceAndAdvertisers.LookalikeAdvertisers {
		advertisers = append(advertisers, &personalization.AudienceAndAdvertisers.LookalikeAdvertisers[i])
	}
	for _, advertiser := range advertisers {
		err = p.DB().
			FirstOrCreate(advertiser, map[string]interface{}{
				"name":      advertiser.Name,
				"lookalike": advertiser.Lookalike,
			}).
			Error
		if err != nil {
			return err
		}
	}

	for i := range personalization.Shows {
		show := &personalization.Shows[i]
		err = p.DB().
			FirstOrCreate(show, map[string]interface{}{
				"name": show.Name,
			}).
			Error
		if err != nil {
			return err
		}
	}

	for i := range personalization.LocationHistory {
		location := &personalization.LocationHistory[i]
		err = p.DB().
			FirstOrCreate(location, map[string]interface{}{
				"name": location.Name,
			}).
			Error
		if err != nil {
			return err
		}
	}

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
