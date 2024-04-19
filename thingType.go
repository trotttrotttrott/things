package main

import (
	"bytes"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/adrg/frontmatter"
)

type thingType struct {
	description string
	Color       string
}

func thingTypes() map[string]thingType {

	thingTypes := make(map[string]thingType)

	dir, err := os.ReadDir(path.Join(thingsDir, "types"))
	if err != nil {
		return thingTypes
	}

	for _, entry := range dir {

		t := thingType{}

		data, err := os.ReadFile(path.Join(thingsDir, "types", entry.Name()))
		if err != nil {
			continue
		}

		rest, err := frontmatter.Parse(bytes.NewReader(data), &t)
		if err != nil {
			continue
		}

		t.description = string(rest)

		thingTypes[strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))] = t
	}

	return thingTypes
}
