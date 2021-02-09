package health

import (
	"encoding/xml"
	"github.com/BionicTeam/bionic/types"
	"gorm.io/gorm"
)

type MeRecord struct {
	gorm.Model
	DataExportID                uint           `gorm:"unique"`
	DateOfBirth                 types.DateTime `xml:"HKCharacteristicTypeIdentifierDateOfBirth,attr"`
	BiologicalSex               string         `xml:"HKCharacteristicTypeIdentifierBiologicalSex,attr"`
	BloodType                   string         `xml:"HKCharacteristicTypeIdentifierBloodType,attr"`
	FitzpatrickSkinType         string         `xml:"HKCharacteristicTypeIdentifierFitzpatrickSkinType,attr"`
	CardioFitnessMedicationsUse string         `xml:"HKCharacteristicTypeIdentifierCardioFitnessMedicationsUse,attr"`
}

func (MeRecord) TableName() string {
	return tablePrefix + "me_records"
}

func (m MeRecord) Constraints() map[string]interface{} {
	return map[string]interface{}{
		"data_export_id": m.DataExportID,
	}
}

func (p *health) parseMe(export *DataExport, decoder *xml.Decoder, start *xml.StartElement) error {
	if err := decoder.DecodeElement(&export.Me, start); err != nil {
		return err
	}

	return p.DB().
		FirstOrCreate(&export.Me, export.Me.Constraints()).
		Error
}
