package things

import (
	"bytes"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/adrg/frontmatter"
)

type Type struct {
	Description string
	Color       string
}

func (ts *Things) TypesPath() string {
	return path.Join(ts.Path, "types")
}

func (ts *Things) ResetTypes() error {

	thingTypes := make(map[string]Type)

	dir, err := os.ReadDir(ts.TypesPath())
	if err != nil {
		return err
	}

	for _, entry := range dir {

		t := Type{}

		data, err := os.ReadFile(path.Join(ts.TypesPath(), entry.Name()))
		if err != nil {
			continue
		}

		rest, err := frontmatter.Parse(bytes.NewReader(data), &t)
		if err != nil {
			continue
		}

		t.Description = string(rest)

		thingTypes[strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))] = t
	}

	ts.Types = thingTypes

	return nil
}
