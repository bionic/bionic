package health

import "encoding/xml"

func (p *health) parseHealthData(export *DataExport, _ *xml.Decoder, start *xml.StartElement) error {
	export.Locale = start.Attr[0].Value

	return nil
}
