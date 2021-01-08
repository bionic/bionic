package twitter

import (
	"encoding/json"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io/ioutil"
	"strings"
)

type Personalization struct {
	gorm.Model
	Languages              []DemographicsLanguage
	GenderInfo             GenderInfo
	Interests              []Interest
	AudienceAndAdvertisers AudienceAndAdvertisers
	Shows                  []Show
	LocationHistory        []Location
	InferredAgeInfo        InferredAgeInfo `json:"inferredAgeInfo"`
}

func (Personalization) TableName() string {
	return "twitter_personalizations"
}

func (p *Personalization) UnmarshalJSON(b []byte) error {
	type alias Personalization

	var data struct {
		alias
		Demographics struct {
			Languages  []DemographicsLanguage `json:"languages"`
			GenderInfo GenderInfo             `json:"genderInfo"`
		} `json:"demographics"`
		Interests struct {
			Interests              []Interest             `json:"interests"`
			AudienceAndAdvertisers AudienceAndAdvertisers `json:"audienceAndAdvertisers"`
			Shows                  []string               `json:"shows"`
		} `json:"interests"`
		LocationHistory []string `json:"locationHistory"`
	}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	*p = Personalization(data.alias)

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

type DemographicsLanguage struct {
	gorm.Model
	PersonalizationID int
	Language          string `json:"language"`
	IsDisabled        bool   `json:"isDisabled"`
}

func (DemographicsLanguage) TableName() string {
	return "twitter_personalization_languages"
}

type GenderInfo struct {
	gorm.Model
	PersonalizationID int
	Gender            string `json:"gender"`
}

func (GenderInfo) TableName() string {
	return "twitter_personalization_gender_infos"
}

type Interest struct {
	gorm.Model
	PersonalizationID int
	Name              string `json:"name"`
	IsDisabled        bool   `json:"isDisabled"`
}

func (Interest) TableName() string {
	return "twitter_personalization_interests"
}

type AudienceAndAdvertisers struct {
	gorm.Model
	PersonalizationID    int
	NumAudiences         int          `json:"numAudiences,string"`
	Advertisers          []Advertiser `gorm:"foreignKey:AudienceAndAdvertisersID"`
	LookalikeAdvertisers []Advertiser `gorm:"foreignKey:AudienceAndAdvertisersID"`
}

func (AudienceAndAdvertisers) TableName() string {
	return "twitter_personalization_audience_and_advertisers"
}

func (aaa *AudienceAndAdvertisers) UnmarshalJSON(b []byte) error {
	type alias AudienceAndAdvertisers

	var data struct {
		alias
		Advertisers          []string `json:"advertisers"`
		LookalikeAdvertisers []string `json:"lookalikeAdvertisers"`
	}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	*aaa = AudienceAndAdvertisers(data.alias)

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
	AudienceAndAdvertisersID *int
	Name                     string
	Lookalike                bool
}

func (Advertiser) TableName() string {
	return "twitter_personalization_advertisers"
}

type Show struct {
	gorm.Model
	PersonalizationID int
	Name              string
}

func (Show) TableName() string {
	return "twitter_personalization_shows"
}

type Location struct {
	gorm.Model
	PersonalizationID int
	Name              string `json:"name"`
}

func (Location) TableName() string {
	return "twitter_personalization_locations"
}

type InferredAgeInfo struct {
	gorm.Model
	PersonalizationID int
	Age               []string `json:"age" gorm:"type:text"`
	BirthDate         string   `json:"birthDate"`
}

func (InferredAgeInfo) TableName() string {
	return "twitter_personalization_inferred_age_infos"
}

func (p *twitter) importPersonalization(inputPath string) error {
	var fileData []struct {
		P13nData Personalization `json:"p13nData"`
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
