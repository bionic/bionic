package health

import (
	"encoding/xml"
	"github.com/bionic-dev/bionic/types"
	"gorm.io/gorm"
)

type Entry struct {
	gorm.Model
	Type            string         `xml:"type,attr" gorm:"uniqueIndex:health_entries_key"`
	SourceName      string         `xml:"sourceName,attr"`
	SourceVersion   string         `xml:"sourceVersion,attr"`
	Unit            string         `xml:"unit,attr"`
	CreationDate    types.DateTime `xml:"creationDate,attr" gorm:"uniqueIndex:health_entries_key"`
	StartDate       types.DateTime `xml:"startDate,attr"`
	EndDate         types.DateTime `xml:"endDate,attr"`
	Value           string         `xml:"value,attr"`
	DeviceID        *int
	Device          *Device             `xml:"device,attr"`
	MetadataEntries []EntryMetadataItem `xml:"MetadataEntry"`
	BeatsPerMinutes []BeatsPerMinute    `xml:"HeartRateVariabilityMetadataList"`
}

func (Entry) TableName() string {
	return tablePrefix + "entries"
}

func (e Entry) Conditions() map[string]interface{} {
	return map[string]interface{}{
		"type":          e.Type,
		"creation_date": e.CreationDate,
	}
}

type EntryMetadataItem struct {
	gorm.Model
	EntryID uint   `gorm:"uniqueIndex:health_entry_metadata_key"`
	Key     string `xml:"key,attr" gorm:"uniqueIndex:health_entry_metadata_key"`
	Value   string `xml:"value,attr"`
}

func (EntryMetadataItem) TableName() string {
	return tablePrefix + "entry_metadata"
}

func (m EntryMetadataItem) Conditions() map[string]interface{} {
	return map[string]interface{}{
		"entry_id": m.EntryID,
		"key":      m.Key,
	}
}

type BeatsPerMinute struct {
	gorm.Model
	EntryID uint   `gorm:"uniqueIndex:health_beats_per_minutes_key"`
	BPM     int    `xml:"bpm,attr"`
	Time    string `xml:"time,attr" gorm:"uniqueIndex:health_beats_per_minutes_key"`
}

func (BeatsPerMinute) TableName() string {
	return tablePrefix + "beats_per_minutes"
}

func (bpm BeatsPerMinute) Conditions() map[string]interface{} {
	return map[string]interface{}{
		"entry_id": bpm.EntryID,
		"time":     bpm.Time,
	}
}

func (e *Entry) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	type Alias Entry

	var data struct {
		Alias
		HeartRateVariabilityMetadataList struct {
			InstantaneousBeatsPerMinute []BeatsPerMinute `xml:"InstantaneousBeatsPerMinute"`
		} `xml:"HeartRateVariabilityMetadataList"`
	}

	if err := decoder.DecodeElement(&data, &start); err != nil {
		return err
	}

	*e = Entry(data.Alias)

	e.BeatsPerMinutes = data.HeartRateVariabilityMetadataList.InstantaneousBeatsPerMinute

	return nil
}

func (p *health) parseRecord(_ *DataExport, decoder *xml.Decoder, start *xml.StartElement) error {
	var entry Entry

	if err := decoder.DecodeElement(&entry, start); err != nil {
		return err
	}

	err := p.DB().
		Find(&entry, entry.Conditions()).
		Error
	if err != nil {
		return err
	}

	if entry.Device != nil {
		err = p.DB().
			FirstOrCreate(entry.Device, entry.Device.Conditions()).
			Error
		if err != nil {
			return err
		}
	}

	for i := range entry.MetadataEntries {
		metadataEntry := &entry.MetadataEntries[i]

		metadataEntry.EntryID = entry.ID

		err = p.DB().
			FirstOrCreate(metadataEntry, metadataEntry.Conditions()).
			Error
		if err != nil {
			return err
		}
	}

	for i := range entry.BeatsPerMinutes {
		beatsPerMinute := &entry.BeatsPerMinutes[i]

		beatsPerMinute.EntryID = entry.ID

		err = p.DB().
			FirstOrCreate(beatsPerMinute, beatsPerMinute.Conditions()).
			Error
		if err != nil {
			return err
		}
	}

	return p.DB().
		FirstOrCreate(&entry, entry.Conditions()).
		Error
}
