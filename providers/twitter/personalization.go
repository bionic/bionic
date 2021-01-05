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
	Demographics    Demographics `json:"demographics"`
	Interests       Interests    `json:"interests"`
	LocationHistory []Location
	InferredAgeInfo InferredAgeInfo `json:"inferredAgeInfo"`
}

func (Personalization) TableName() string {
	return "twitter_personalizations"
}

func (p *Personalization) UnmarshalJSON(b []byte) error {
	type alias Personalization

	var data struct {
		alias
		LocationHistory []string `json:"locationHistory"`
	}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	*p = Personalization(data.alias)

	for _, location := range data.LocationHistory {
		p.LocationHistory = append(p.LocationHistory, Location{
			Name: location,
		})
	}

	return nil
}

type Demographics struct {
	gorm.Model
	PersonalizationID int
	Languages         []DemographicsLanguage `json:"languages"`
	GenderInfo        GenderInfo             `json:"genderInfo"`
}

func (Demographics) TableName() string {
	return "twitter_personalization_demographics"
}

type DemographicsLanguage struct {
	gorm.Model
	DemographicsID int
	Language       string `json:"language"`
	IsDisabled     bool   `json:"isDisabled"`
}

func (DemographicsLanguage) TableName() string {
	return "twitter_personalization_demographic_languages"
}

type GenderInfo struct {
	gorm.Model
	DemographicsID int
	Gender         string `json:"gender"`
}

func (GenderInfo) TableName() string {
	return "twitter_personalization_gender_infos"
}

type Interests struct {
	gorm.Model
	PersonalizationID      int
	Interests              []Interest             `json:"interests"`
	AudienceAndAdvertisers AudienceAndAdvertisers `json:"audienceAndAdvertisers"`
	Shows                  []Show
}

func (Interests) TableName() string {
	return "twitter_personalization_interestses"
}

func (i *Interests) UnmarshalJSON(b []byte) error {
	type alias Interests

	var data struct {
		alias
		Shows []string `json:"shows"`
	}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	*i = Interests(data.alias)

	for _, show := range data.Shows {
		i.Shows = append(i.Shows, Show{
			Name: show,
		})
	}

	return nil
}

type Interest struct {
	gorm.Model
	InterestsID int
	Name        string `json:"name"`
	IsDisabled  bool   `json:"isDisabled"`
}

func (Interest) TableName() string {
	return "twitter_personalization_interests"
}

type AudienceAndAdvertisers struct {
	gorm.Model
	InterestsID          int
	NumAudiences         int          `json:"numAudiences,string"`
	Advertisers          []Advertiser `gorm:"foreignKey:AdvertisersID"`
	LookalikeAdvertisers []Advertiser `gorm:"foreignKey:LookalikeAdvertisersID"`
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
				Name: lookalikeAdvertiser,
			},
		)
	}

	return nil
}

type Advertiser struct {
	gorm.Model
	AdvertisersID          *int
	LookalikeAdvertisersID *int
	Name                   string
}

func (Advertiser) TableName() string {
	return "twitter_personalization_advertisers"
}

type Show struct {
	gorm.Model
	InterestsID int
	Name        string
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
