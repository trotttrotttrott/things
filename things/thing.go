package things

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Thing struct {
	Title    string
	Type     string
	Priority int
	Done     bool
	Pause    bool
	Today    bool
	Pin      bool

	Content  string
	Path     string
	TimePath string
}

func (t *Thing) Age() string {
	n := time.Now().UTC()
	b := filepath.Base(t.Path)
	tm, err := time.Parse("20060102150405", strings.TrimSuffix(b, filepath.Ext(b)))
	if err != nil {
		return ""
	}

	dur := n.Sub(tm).Hours() / 24
	switch {
	case dur >= 1000:
		return "+k"
	default:
		return fmt.Sprintf("%d", int(dur))
	}
}

func (t *Thing) Time() (timeSpent time.Duration) {

	if _, err := os.Stat(t.TimePath); errors.Is(err, os.ErrNotExist) {
		return
	}
	data, err := os.ReadFile(t.TimePath)
	if err != nil {
		return
	}

	rdr := csv.NewReader(bytes.NewReader(data))
	records, err := rdr.ReadAll()
	if err != nil {
		return
	}

	for _, r := range records {
		start, err := time.Parse(time.RFC3339, r[0])
		if err != nil {
			continue
		}
		end, err := time.Parse(time.RFC3339, r[1])
		if err != nil {
			continue
		}
		timeSpent += end.Sub(start)
	}

	return
}

// TimeString formats the time spent in a compact, human-friendly way.
// Examples: "45s", "25m", "2h30m", "3d2h", "14d"
func (t *Thing) TimeString() string {
	d := t.Time()
	if d == 0 {
		return "0s"
	}

	// For durations less than a minute, show seconds
	if d < time.Minute {
		seconds := int(d.Round(time.Second).Seconds())
		return fmt.Sprintf("%ds", seconds)
	}

	// Round to the nearest minute for display
	minutes := int(d.Round(time.Minute).Minutes())

	days := minutes / (60 * 24)
	hours := (minutes % (60 * 24)) / 60
	mins := minutes % 60

	if days > 0 {
		if hours > 0 {
			return fmt.Sprintf("%dd%dh", days, hours)
		}
		return fmt.Sprintf("%dd", days)
	}

	if hours > 0 {
		if mins > 0 {
			return fmt.Sprintf("%dh%dm", hours, mins)
		}
		return fmt.Sprintf("%dh", hours)
	}

	return fmt.Sprintf("%dm", mins)
}

func (t *Thing) Remove() error {
	err := os.Remove(t.Path)
	if err != nil {
		return err
	}
	return os.Remove(t.TimePath)
}
