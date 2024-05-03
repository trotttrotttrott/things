package main

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
)

type thing struct {
	Title    string
	Type     string
	Priority int
	Done     bool
	Pause    bool
	Today    bool
	content  string
	path     string
	timePath string
}

func (t *thing) age() string {
	n := time.Now().UTC()
	b := filepath.Base(t.path)
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

func (t *thing) remove() error {
	err := os.Remove(t.path)
	if err != nil {
		return err
	}
	return os.Remove(t.timePath)
}

func (t *thing) time() (timeSpent time.Duration) {

	if _, err := os.Stat(t.timePath); errors.Is(err, os.ErrNotExist) {
		return
	}
	data, err := os.ReadFile(t.timePath)
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

func things(filter string) (things []thing) {

	dir, err := os.ReadDir(path.Join(thingsDir, "things"))
	if err != nil {
		return
	}

	for _, entry := range dir {

		t := thing{
			path: path.Join(thingsDir, "things", entry.Name()),
		}
		t.timePath = path.Join(thingsDir, "time", fmt.Sprintf("%s.csv", strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))))

		data, err := os.ReadFile(t.path)
		if err != nil {
			continue
		}

		rest, err := frontmatter.Parse(bytes.NewReader(data), &t)
		if err != nil {
			continue
		}

		switch filter {
		case "current":
			if t.Done {
				continue
			}
		case "done":
			if !t.Done {
				continue
			}
		case "pause":
			if !t.Pause {
				continue
			}
		case "today":
			if !t.Today {
				continue
			}
		}

		t.content = string(rest)

		things = append(things, t)
	}

	return
}

func thingNew(thingTypeKeys []string) (t thing, err error) {

	now := time.Now().UTC().Format("20060102150405")

	t.path = path.Join(thingsDir, "things", fmt.Sprintf("%s.md", now))
	t.timePath = path.Join(thingsDir, "time", fmt.Sprintf("%s.csv", now))

	f, err := os.Create(t.path)
	if err != nil {
		return
	}

	defer f.Close()

	_, err = f.WriteString(strings.Join(
		[]string{
			"---",
			"title: Thing",
			fmt.Sprintf("type: # %s", strings.Join(thingTypeKeys, " ")),
			"priority: 0",
			"---",
			"",
		}, "\n"))

	if err != nil {
		return
	}

	err = f.Sync()

	return
}
