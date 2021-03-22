package markdown

import (
	"fmt"
	"strings"
	"time"
)

func formatDuration(d time.Duration, round time.Duration) string {
	formatUnit := func(d time.Duration, unit string) string {
		result := fmt.Sprintf("%d %s", d, unit)
		if d != 1 {
			result += "s"
		}

		return result
	}

	d = d.Round(round)

	h := d / time.Hour
	d -= h * time.Hour

	m := d / time.Minute
	d -= m * time.Minute

	s := d / time.Second

	var parts []string
	if h > 0 {
		parts = append(parts, formatUnit(h, "hour"))
	}
	if m > 0 && round < time.Hour {
		parts = append(parts, formatUnit(m, "minute"))
	}
	if s > 0 && round < time.Minute {
		parts = append(parts, formatUnit(s, "second"))
	}

	return strings.Join(parts, " ")
}
