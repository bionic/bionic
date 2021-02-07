package health

import "encoding/xml"

func (p *health) parseExportDate(data *Data, _ *xml.Decoder, start *xml.StartElement) error {
	if err := data.ExportDate.UnmarshalText([]byte(start.Attr[0].Value)); err != nil {
		return err
	}

	err := p.DB().
		Find(&data, data.Constraints()).
		Error
	if err != nil {
		return err
	}

	data.Me.DataID = data.ID

	return nil
}
