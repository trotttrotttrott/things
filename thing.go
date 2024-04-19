package main

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
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
		log.Fatalln("Error:", err)
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

func (t *thing) time() time.Duration {

	var timeSpent time.Duration

	if _, err := os.Stat(t.timePath); errors.Is(err, os.ErrNotExist) {
		return timeSpent
	}
	data, err := os.ReadFile(t.timePath)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	rdr := csv.NewReader(bytes.NewReader(data))
	records, err := rdr.ReadAll()
	if err != nil {
		log.Fatalln("Error:", err)
	}

	for _, r := range records {
		start, err := time.Parse(time.RFC3339, r[0])
		if err != nil {
			log.Fatalln("Error:", err)
		}
		end, err := time.Parse(time.RFC3339, r[1])
		if err != nil {
			log.Fatalln("Error:", err)
		}
		timeSpent += end.Sub(start)
	}

	return timeSpent
}

func things(filter string) (things []thing) {

	dir, err := os.ReadDir(path.Join(thingsDir, "things"))
	if err != nil {
		log.Fatalln("Error:", err)
	}

	for _, entry := range dir {

		t := thing{
			path: path.Join(thingsDir, "things", entry.Name()),
		}
		t.timePath = path.Join(thingsDir, "time", fmt.Sprintf("%s.csv", strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))))

		data, err := os.ReadFile(t.path)
		if err != nil {
			log.Fatalln("Error:", err)
		}

		rest, err := frontmatter.Parse(bytes.NewReader(data), &t)
		if err != nil {
			log.Fatalln("Error:", err)
		}

		switch filter {
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
		default:
			if t.Done {
				continue
			}
		}

		t.content = string(rest)

		things = append(things, t)
	}

	return
}
