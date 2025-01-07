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

func (t *Thing) Remove() error {
	err := os.Remove(t.Path)
	if err != nil {
		return err
	}
	return os.Remove(t.TimePath)
}
