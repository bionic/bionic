package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/araddon/dateparse"
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
	t, ok := src.(time.Time)
	if !ok {
		return fmt.Errorf("failed to scan value into time.Time: %+v", src)
	}

	*dt = DateTime(t)

	return nil
}

func (dt DateTime) Value() (driver.Value, error) {
	if time.Time(dt).IsZero() {
		return nil, nil
	}

	return time.Time(dt), nil
}
