package types

import (
	"database/sql"
	"strconv"
)

type NullableBool struct{ sql.NullBool }

func (b *NullableBool) UnmarshalCSV(csv string) (err error) {
	value, err := strconv.ParseBool(csv)

	if err == nil {
		b.Valid = true
		b.Bool = value
	} else {
		b.Valid = false
	}

	return nil
}
