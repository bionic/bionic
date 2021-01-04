package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
)

type IntString int

func (is *IntString) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		value, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		*is = IntString(value)
		return nil
	}

	return json.Unmarshal(b, (*int)(is))
}

func (is *IntString) Scan(src interface{}) error {
	i, ok := src.(int)
	if !ok {
		return fmt.Errorf("failed to scan value into int: %+v", src)
	}

	*is = IntString(i)

	return nil
}

func (is IntString) Value() (driver.Value, error) {
	return json.Marshal(int(is))
}

type IntStringSlice []IntString

func (iss *IntStringSlice) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan value into []byte: %+v", src)
	}

	return json.Unmarshal(b, iss)
}

func (iss IntStringSlice) Value() (driver.Value, error) {
	return json.Marshal(iss)
}
