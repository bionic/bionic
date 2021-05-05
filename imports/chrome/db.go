package chrome

import (
	"github.com/bionic-dev/bionic/types"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const dbRowSelectLimit = 100

type DbURL struct {
	ID            int
	URL           string
	Title         string
	VisitCount    int
	TypedCount    int
	LastVisitTime types.DateTime
	Hidden        bool
}

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
