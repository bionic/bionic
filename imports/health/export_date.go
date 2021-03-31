package health

import "encoding/xml"

func (p *health) parseExportDate(export *DataExport, _ *xml.Decoder, start *xml.StartElement) error {
	if err := export.ExportDate.UnmarshalText([]byte(start.Attr[0].Value)); err != nil {
		return err
	}

	err := p.DB().
		Select("ID").
		Find(&export, export.Conditions()).
		Error
	if err != nil {
		return err
	}

	export.Me.DataExportID = export.ID

	return nil
}
