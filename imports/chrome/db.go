package chrome

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"net/url"
	"os"
)

const dbRowSelectLimit = 100

func (p *chrome) importDB(inputPath string) error {
	u, err := url.Parse(inputPath)
	if err != nil {
		return err
	}

	sourceFile, err := os.Open(u.Path)
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

	db, err := gorm.Open(sqlite.Open(tmpFile.Name()+"?"+u.RawQuery), &gorm.Config{})

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
