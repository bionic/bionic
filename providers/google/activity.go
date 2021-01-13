package google

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"github.com/shekhirin/bionic-cli/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io/ioutil"
	"path/filepath"
)

type Action struct {
	gorm.Model
	Header        string         `json:"header"`
	Title         string         `json:"title"`
	TitleUrl      string         `json:"titleUrl"`
	Time          types.DateTime `json:"time"`
	Products      []Product      `json:"products" gorm:"many2many:google_action_products;"`
	LocationInfos []LocationInfo
	Subtitles     []Subtitle
	Details       []Detail `json:"details" gorm:"many2many:google_action_details;"`
}

func (a Action) TableName() string {
	return "google_activity"
}

type Product struct {
	gorm.Model
	Name string `gorm:"unique"`
}

func (p Product) TableName() string {
	return "google_activity_products"
}

type LocationInfo struct {
	gorm.Model
	ActionID  int
	Action    Action
	Name      string `json:"name"`
	Url       string `json:"url" `
	Source    string `json:"source"`
	SourceUrl string `json:"sourceUrl"`
}

func (i LocationInfo) TableName() string {
	return "google_activity_location_infos"
}

type Subtitle struct {
	gorm.Model
	ActionID int
	Action   Action
	Name     string `json:"name"`
	Url      string `json:"url"`
}

func (s Subtitle) TableName() string {
	return "google_activity_subtitles"
}

type Detail struct {
	gorm.Model
	Name string `json:"name" gorm:"unique"`
}

func (d Detail) TableName() string {
	return "google_activity_details"
}

func (p *Product) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}

	*p = Product{Name: str}
	return nil
}

// TODO: Rebuild on https://stackoverflow.com/questions/31794355/decode-large-stream-json

func (p *google) importActivity(inputPath string) error {
	r, err := zip.OpenReader(inputPath)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}

		filename := filepath.Base(f.Name)

		if filename != "MyActivity.json" {
			continue
		}

		bytes, err := ioutil.ReadAll(rc)
		if err != nil {
			return nil
		}

		fmt.Println(f.Name)

		var actions []Action
		if err := json.Unmarshal(bytes, &actions); err != nil {
			return err
		}

		err = p.DB().
			Clauses(clause.OnConflict{
				DoNothing: true,
			}).
			CreateInBatches(actions, 1000).
			Error
		if err != nil {
			return err
		}

		if err := rc.Close(); err != nil {
			return err
		}

	}

	return nil
}
