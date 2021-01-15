package google

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"github.com/shekhirin/bionic-cli/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

const actionBatchSize = 100

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
		filename := filepath.Base(f.Name)

		if filename != "MyActivity.json" {
			continue
		}
		fmt.Println(f.Name)

		err := processActionsFile(p.DB(), f)
		if err != nil {
			return err
		}
	}

	return nil
}

func processActionsFile(db *gorm.DB, file *zip.File) error {
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer func() {
		if err := rc.Close(); err != nil {
			panic(err)
		}
	}()

	decoder := json.NewDecoder(rc)
	if _, err := decoder.Token(); err != nil {
		return err
	} // Skip first token, which is opening the list

	var actionBatch []Action

	for decoder.More() {
		var action Action
		err := decoder.Decode(&action)
		if err != nil {
			return err
		}

		actionBatch = append(actionBatch, action)
		if len(actionBatch) >= actionBatchSize {
			if err := saveActions(db, actionBatch); err != nil {
				return err
			}
			actionBatch = nil
		}
	}

	if err := saveActions(db, actionBatch); err != nil {
		return err
	}

	return nil
}

func saveActions(db *gorm.DB, actions []Action) error {
	for i, action := range actions {
		for j, product := range action.Products {
			err := db.
				FirstOrCreate(&actions[i].Products[j], map[string]interface{}{"name": product.Name}).
				Error
			if err != nil {
				return err
			}
		}
	}

	err := db.
		Clauses(clause.OnConflict{
			DoNothing: true,
		}).
		CreateInBatches(actions, 1000).
		Error
	return err
}
