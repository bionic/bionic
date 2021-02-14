package health

import "gorm.io/gorm"

type MetadataEntry struct {
	gorm.Model
	ParentID   uint   `gorm:"uniqueIndex:health_metadata_entries_key"`
	ParentType string `gorm:"uniqueIndex:health_metadata_entries_key"`
	Key        string `xml:"key,attr" gorm:"uniqueIndex:health_metadata_entries_key"`
	Value      string `xml:"value,attr"`
}

func (MetadataEntry) TableName() string {
	return tablePrefix + "metadata_entries"
}

func (e MetadataEntry) Conditions() map[string]interface{} {
	return map[string]interface{}{
		"parent_id":   e.ParentID,
		"parent_type": e.ParentType,
		"key":         e.Key,
	}
}
