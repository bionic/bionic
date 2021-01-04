package netflix

import (
	"strconv"
	"strings"
	"time"
)

type duration time.Duration

func (d *duration) UnmarshalCSV(csv string) error {
	parts := strings.Split(csv, ":")
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

	*d = duration(time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second)

	return nil
}
