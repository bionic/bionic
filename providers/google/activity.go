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
	Header        string         `json:"header" gorm:"uniqueIndex:idx_google_activity"`
	Title         string         `json:"title" gorm:"uniqueIndex:idx_google_activity"`
	TitleUrl      string         `json:"titleUrl"`
	Time          types.DateTime `json:"time" gorm:"uniqueIndex:idx_google_activity"`
	Products      []Product      `json:"products" gorm:"many2many:google_activity_products_assoc;"`
	LocationInfos []LocationInfo
	Subtitles     []Subtitle
	Details       []Detail `json:"details"`
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

type ActionProductAssoc struct {
	ActionID  int `gorm:"primaryKey;not null"`
	ProductID int `gorm:"primaryKey;not null"`
}

func (a ActionProductAssoc) TableName() string {
	return "google_activity_products_assoc"
}

func (p *Product) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}

	*p = Product{Name: str}
	return nil
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
	ActionID int
	Action   Action
	Name     string `json:"name"`
}

func (d Detail) TableName() string {
	return "google_activity_details"
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

		for i, action := range actions {
			for j, product := range action.Products {
				err = p.DB().
					FirstOrCreate(&actions[i].Products[j], map[string]interface{}{"name": product.Name}).
					Error
				if err != nil {
					return err
				}
			}
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
