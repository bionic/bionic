package chrome

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const dbRowSelectLimit = 100

func (p *chrome) importDB(inputPath string) error {
	db, err := gorm.Open(sqlite.Open(inputPath), &gorm.Config{})

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
