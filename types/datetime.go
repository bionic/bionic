package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/mattn/go-sqlite3"
	"time"
)

type DateTime time.Time

func (dt *DateTime) UnmarshalText(text []byte) error {
	t, err := dateparse.ParseStrict(string(text))
	if err != nil {
		return err
	}

	*dt = DateTime(t)

	return nil
}

func (dt *DateTime) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}

	return dt.UnmarshalText([]byte(str))
}

func (dt *DateTime) UnmarshalCSV(csv string) (err error) {
	if csv == "" {
		return nil
	}

	return dt.UnmarshalText([]byte(csv))
}

func (dt *DateTime) Scan(src interface{}) error {
	switch v := src.(type) {
	case time.Time:
		*dt = DateTime(v)
		return nil
	case string:
		return dt.UnmarshalText([]byte(v))
	case []byte:
		return dt.UnmarshalText(v)
	default:
		return fmt.Errorf("failed to scan value into DateTime: %+v", src)
	}
}

func (dt DateTime) Value() (driver.Value, error) {
	if time.Time(dt).IsZero() {
		return nil, nil
	}

	return time.Time(dt).Format(sqlite3.SQLiteTimestampFormats[0]), nil
}
