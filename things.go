package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
)

type thing struct {
	Title    string
	Type     string
	Priority int
	Done     bool
	content  string
	path     string
}

func (t *thing) thingType() thingType {
	return thingTypes[t.Type]
}

func (t *thing) age() string {
	n := time.Now().UTC()
	b := filepath.Base(t.path)
	tm, err := time.Parse("20060102150405", strings.TrimSuffix(b, filepath.Ext(b)))
	if err != nil {
		log.Fatalln("Error:", err)
	}
	return fmt.Sprintf("%.2f", n.Sub(tm).Hours()/24)
}

type thingType struct {
	description string
	Color       string
}

func typesInit() {

	thingTypes = map[string]thingType{}

	dir, err := os.ReadDir(path.Join(thingsDir, "types"))
	if err != nil {
		log.Fatalln("Error:", err)
	}

	for _, entry := range dir {

		t := thingType{}

		data, err := os.ReadFile(path.Join(thingsDir, "types", entry.Name()))
		if err != nil {
			log.Fatalln("Error:", err)
		}

		rest, err := frontmatter.Parse(bytes.NewReader(data), &t)
		if err != nil {
			log.Fatalln("Error:", err)
		}

		t.description = string(rest)

		thingTypes[strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))] = t
	}
}

func things() (things []thing) {

	dir, err := os.ReadDir(path.Join(thingsDir, "things"))
	if err != nil {
		log.Fatalln("Error:", err)
	}

	for _, entry := range dir {

		t := thing{
			path: path.Join(thingsDir, "things", entry.Name()),
		}

		data, err := os.ReadFile(t.path)
		if err != nil {
			log.Fatalln("Error:", err)
		}

		rest, err := frontmatter.Parse(bytes.NewReader(data), &t)
		if err != nil {
			log.Fatalln("Error:", err)
		}

		if t.Done {
			continue
		}

		t.content = string(rest)

		things = append(things, t)
	}

	sort.Slice(things, func(i, j int) bool {
		return things[i].Priority < things[j].Priority
	})

	return things
}
