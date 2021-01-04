package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/araddon/dateparse"
	"time"
)

type DateTime time.Time

func (dt *DateTime) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}

	t, err := dateparse.ParseStrict(str)
	if err != nil {
		return err
	}

	*dt = DateTime(t)

	return nil
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
	return time.Time(dt), nil
}
