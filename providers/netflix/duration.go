package netflix

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

type Duration struct {
	sql.NullInt64
}

func (d *Duration) UnmarshalCSV(csv string) error {
	if csv == "Not latest view" {
		d.Valid = false
		return nil
	}

	parts := strings.Split(csv, ":")
	if len(parts) != 3 {
		return fmt.Errorf("incorrect duration format: %s", csv)
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return err
	}

	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}

	seconds, err := strconv.Atoi(parts[2])
	if err != nil {
		return err
	}

	d.Valid = true
	d.Int64 = int64(hours*3600 + minutes*60 + seconds)

	return nil
}
