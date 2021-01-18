package google

import (
	"archive/zip"
	"encoding/json"
	"github.com/shekhirin/bionic-cli/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io"
	"os"
	"path/filepath"
)

const actionBatchSize = 100
const targetFilename = "MyActivity.json"

type Action struct {
	gorm.Model
	Header        string         `json:"header" gorm:"uniqueIndex:google_activity_key"`
	Title         string         `json:"title" gorm:"uniqueIndex:google_activity_key"`
	TitleURL      string         `json:"titleUrl"`
	Time          types.DateTime `json:"time" gorm:"uniqueIndex:google_activity_key"`
	Products      []Product      `json:"products" gorm:"many2many:google_activity_products_assoc;"`
	LocationInfos []LocationInfo
	Subtitles     []Subtitle
	Details       []Detail `json:"details"`
}

func (Action) TableName() string {
	return "google_activity"
}

type Product struct {
	gorm.Model
	Name string `gorm:"unique"`
}

func (Product) TableName() string {
	return "google_activity_products"
}

type ActionProductAssoc struct {
	ActionID  int `gorm:"primaryKey;not null"`
	ProductID int `gorm:"primaryKey;not null"`
}

func (ActionProductAssoc) TableName() string {
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
	URL       string `json:"url" `
	Source    string `json:"source"`
	SourceURL string `json:"sourceUrl"`
}

func (LocationInfo) TableName() string {
	return "google_activity_location_infos"
}

type Subtitle struct {
	gorm.Model
	ActionID int
	Action   Action
	Name     string `json:"name"`
	URL      string `json:"url"`
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

func (p *google) importActivityFromArchive(inputPath string) error {
	r, err := zip.OpenReader(inputPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = r.Close()
	}()

	for _, f := range r.File {
		filename := filepath.Base(f.Name)
		if filename != targetFilename {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		if err := p.processActionsFile(rc); err != nil {
			return err
		}
		if err := rc.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (p *google) importActivityFromDirectory(inputPath string) error {
	err := filepath.Walk(inputPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Name() != targetFilename {
				return nil
			}

			rc, err := os.Open(path)
			if err != nil {
				return err
			}

			err = p.processActionsFile(rc)
			if err != nil {
				return err
			}

			return nil
		})
	if err != nil {
		return err
	}
	return nil
}

func (p *google) processActionsFile(rc io.ReadCloser) error {
	decoder := json.NewDecoder(rc)
	if _, err := decoder.Token(); err != nil {
		return err
	} // Skip first token, which is opening the list

	var batch []Action

	for decoder.More() {
		var action Action
		err := decoder.Decode(&action)
		if err != nil {
			return err
		}

		batch = append(batch, action)
		if len(batch) >= actionBatchSize {
			if err := p.saveActions(batch); err != nil {
				return err
			}
			batch = nil
		}
	}

	if err := p.saveActions(batch); err != nil {
		return err
	}

	return nil
}

func (p *google) saveActions(actions []Action) error {
	for i, action := range actions {
		for j, product := range action.Products {
			err := p.DB().
				FirstOrCreate(&actions[i].Products[j], map[string]interface{}{"name": product.Name}).
				Error
			if err != nil {
				return err
			}
		}
	}

	err := p.DB().
		Clauses(clause.OnConflict{
			DoNothing: true,
		}).
		CreateInBatches(actions, 1000).
		Error
	return err
}
