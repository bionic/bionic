package chrome

import (
	"gorm.io/gorm"
)

type Segment struct {
	gorm.Model
	Name  string `gorm:"uniqueIndex:chrome_segments_key"`
	URLID int    `gorm:"uniqueIndex:chrome_segments_key"`
	URL   URL
}

func (Segment) TableName() string {
	return tablePrefix + "segments"
}

func (s Segment) Conditions() map[string]interface{} {
	return map[string]interface{}{
		"url_id": s.URLID,
	}
}

func (p *chrome) importSegments(db *gorm.DB) error {
	selection := "id, name, url_id"

	var segments []Segment
	err := db.
		Raw("select "+selection+" from segments order by id limit ?", dbRowSelectLimit).
		Scan(&segments).
		Error
	if err != nil {
		return err
	}
	if err := p.saveSegments(segments); err != nil {
		return err
	}
	for len(segments) == dbRowSelectLimit {
		lastSegment := segments[len(segments)-1]
		err = db.
			Raw("select "+selection+" from segments where id > ? order by id limit ?", lastSegment.ID, dbRowSelectLimit).
			Scan(&segments).
			Error
		if err != nil {
			return err
		}
		if err := p.saveSegments(segments); err != nil {
			return err
		}
	}

	return nil
}

func (p *chrome) saveSegments(segments []Segment) error {
	for i, segment := range segments {
		err := p.DB().
			FirstOrCreate(&segments[i], segment.Conditions()).
			Error
		if err != nil {
			return err
		}
	}

	return nil
}
