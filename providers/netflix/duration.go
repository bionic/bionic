package netflix

import (
	"fmt"
	"strconv"
	"strings"
)

type Duration int

func (d *Duration) UnmarshalCSV(csv string) error {
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

	*d = Duration(hours*3600 + minutes*60 + seconds)

	return nil
}
