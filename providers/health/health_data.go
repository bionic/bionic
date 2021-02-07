package health

import "encoding/xml"

func (p *health) parseHealthData(data *Data, _ *xml.Decoder, start *xml.StartElement) error {
	data.Locale = start.Attr[0].Value

	return nil
}
