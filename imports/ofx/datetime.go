package ofx

import (
	"database/sql/driver"
	"fmt"
	"github.com/aclindsa/ofxgo"
	"github.com/aclindsa/xml"
	"github.com/mattn/go-sqlite3"
	"time"
)

type DateTime ofxgo.Date

func (dt *DateTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var ofxDate ofxgo.Date

	err := ofxDate.UnmarshalXML(d, start)
	if err != nil {
		return err
	}
	*dt = DateTime(ofxDate)

	return nil
}

func (dt *DateTime) Scan(src interface{}) error {
	switch v := src.(type) {
	case time.Time:
		dt.Time = v
		return nil
	default:
		return fmt.Errorf("failed to scan value into DateTime: %+v", src)
	}
}

func (dt DateTime) Value() (driver.Value, error) {
	if dt.Time.IsZero() {
		return nil, nil
	}

	return dt.Time.Format(sqlite3.SQLiteTimestampFormats[0]), nil
}
