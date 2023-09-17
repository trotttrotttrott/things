package main

import (
	"bytes"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

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

func thingsLoad() {

	dir, err := os.ReadDir(path.Join(thingsDir, "things"))
	if err != nil {
		log.Fatalln("Error:", err)
	}

	var thngs []thing

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

		t.content = string(rest)

		thngs = append(thngs, t)
	}

	sort.Slice(thngs, func(i, j int) bool {
		return thngs[i].Priority < thngs[j].Priority
	})

	things = thngs
}
