package health

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"github.com/BionicTeam/bionic/types"
	"gorm.io/gorm"
	"io"
	"os"
	"path"
	"path/filepath"
)

type DataExport struct {
	gorm.Model
	Locale     string
	ExportDate types.DateTime `gorm:"unique"`
	Me         MeRecord       `xml:"Me"`
	Workouts   []*Workout     `gorm:"-"`
}

func (DataExport) TableName() string {
	return tablePrefix + "data_exports"
}

func (d DataExport) Constraints() map[string]interface{} {
	return map[string]interface{}{
		"export_date": d.ExportDate,
	}
}

func (p *health) importDataExport(r io.Reader) (*DataExport, error) {
	var data DataExport

	decoder := xml.NewDecoder(r)

	for {
		token, err := decoder.Token()
		if token == nil || err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		switch typ := token.(type) {
		case xml.StartElement:
			var parseFn func(*DataExport, *xml.Decoder, *xml.StartElement) error

			switch typ.Name.Local {
			case "HealthData":
				parseFn = p.parseHealthData
			case "ExportDate":
				parseFn = p.parseExportDate
			case "Me":
				parseFn = p.parseMe
			case "Record":
				parseFn = p.parseRecord
			case "Workout":
				parseFn = p.parseWorkout
			case "ActivitySummary":
				parseFn = p.parseActivitySummary
			default:
				continue
			}

			if err := parseFn(&data, decoder, &typ); err != nil {
				return nil, err
			}
		}
	}

	err := p.DB().
		FirstOrCreate(&data, data.Constraints()).
		Error
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (p *health) importDataExportFromArchive(inputPath string) error {
	var export *DataExport

	r, err := zip.OpenReader(inputPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = r.Close()
	}()

	workoutRouteFiles := map[string]io.ReadCloser{}

	for _, f := range r.File {
		if filepath.Base(f.Name) == "export.xml" {
			rc, err := f.Open()
			if err != nil {
				return err
			}

			export, err = p.importDataExport(rc)
			if err != nil {
				return err
			}

			if err := rc.Close(); err != nil {
				return err
			}
		} else if filepath.Base(filepath.Dir(f.Name)) == "workout-routes" {
			rc, err := f.Open()
			if err != nil {
				return nil
			}

			workoutRouteFiles[filepath.Base(f.Name)] = rc
		}
	}

	if export == nil {
		return errors.New("no export.xml file found")
	}

	return p.importWorkoutRoutes(export, workoutRouteFiles)
}

func (p *health) importDataExportFromDirectory(inputPath string) error {
	var export *DataExport

	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	export, err = p.importDataExport(f)
	if err != nil {
		return err
	}

	workoutRouteFiles := map[string]io.ReadCloser{}

	for _, workout := range export.Workouts {
		if route := workout.Route; route != nil {
			r, err := os.Open(path.Join(path.Dir(inputPath), route.FilePath))
			if err != nil {
				return err
			}

			workoutRouteFiles[filepath.Base(route.FilePath)] = r
		}
	}

	return p.importWorkoutRoutes(export, workoutRouteFiles)
}
