package chrome

import (
	"github.com/bionic-dev/bionic/types"
	"gorm.io/gorm"
)

type URL struct {
	gorm.Model
	URL        string `gorm:"unique"`
	Title      string
	VisitCount int
	TypedCount int
	LastVisit  types.DateTime
	Hidden     bool
}

func (URL) TableName() string {
	return tablePrefix + "urls"
}

func (u URL) Conditions() map[string]interface{} {
	return map[string]interface{}{
		"url": u.URL,
	}
}

func (p *chrome) importURLs(db *gorm.DB) error {
	selection := "id, url, title, visit_count, typed_count, " +
		"datetime((last_visit_time/1000000)-11644473600, 'unixepoch') as last_visit, hidden"

	var urls []URL
	err := db.
		Raw("select "+selection+" from urls order by id limit ?", dbRowSelectLimit).
		Scan(&urls).
		Error
	if err != nil {
		return err
	}
	if err := p.saveURLs(urls); err != nil {
		return err
	}
	for len(urls) == dbRowSelectLimit {
		lastUrl := urls[len(urls)-1]
		err = db.
			Raw("select "+selection+" from urls where id > ? order by id limit ?", lastUrl.ID, dbRowSelectLimit).
			Scan(&urls).
			Error
		if err != nil {
			return err
		}
		if err := p.saveURLs(urls); err != nil {
			return err
		}
	}

	return nil
}

func (p *chrome) saveURLs(urls []URL) error {
	for i, url := range urls {
		err := p.DB().
			FirstOrCreate(&urls[i], url.Conditions()).
			Error
		if err != nil {
			return err
		}
	}

	return nil
}
