package health

import "encoding/xml"

func (p *health) parseMe(data *Data, decoder *xml.Decoder, start *xml.StartElement) error {
	if err := decoder.DecodeElement(&data.Me, start); err != nil {
		return err
	}

	return p.DB().
		FirstOrCreate(&data.Me, data.Me.Constraints()).
		Error
}
