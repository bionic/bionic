package chrome

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

const dbRowSelectLimit = 100

func (p *chrome) importDB(inputPath string) error {
	// strip sqlite arguments
	var filePath string
	var arguments string
	parts := strings.Split(inputPath, "?")
	if len(parts) > 1 {
		filePath = strings.Join(parts[:len(parts)-1], "?")
		arguments = parts[len(parts)-1]
	} else {
		filePath = parts[0]
	}

	sourceFile, err := os.Open(filePath)
	if err != nil {
		return err
	}

	tmpFile, err := ioutil.TempFile("", "bionic-chrome-copy.*.sqlite")
	if err != nil {
		return err
	}

	defer os.Remove(tmpFile.Name())

	_, err = io.Copy(tmpFile, sourceFile)
	if err != nil {
		return err
	}

	db, err := gorm.Open(sqlite.Open(tmpFile.Name()+"?"+arguments), &gorm.Config{})

	if err != nil {
		return err
	}

	if err := p.importURLs(db); err != nil {
		return err
	}

	if err := p.importSegments(db); err != nil {
		return err
	}

	if err := p.importVisits(db); err != nil {
		return err
	}

	return nil
}
