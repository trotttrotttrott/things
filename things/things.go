package things

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
)

type Things struct {
	Things []Thing
	Types  map[string]Type
	Path   string
	sort   string
	filter string
}

func New(p string) Things {
	t := Things{
		Path:   p,
		sort:   "priority",
		filter: "current",
	}

	t.ResetThings()
	t.ResetTypes()

	return t
}

func (ts *Things) NewThing(thingTypeKeys []string) (t Thing, err error) {

	now := time.Now().UTC().Format("20060102150405")

	t.Path = path.Join(ts.Path, "things", fmt.Sprintf("%s.md", now))
	t.TimePath = path.Join(ts.Path, "time", fmt.Sprintf("%s.csv", now))

	f, err := os.Create(t.Path)
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

func (ts *Things) Filter(filter string) error {
	ts.filter = filter
	return ts.ResetThings()
}

func (ts *Things) ResetThings() error {

	dir, err := os.ReadDir(path.Join(ts.Path, "things"))
	if err != nil {
		return err
	}

	var list []Thing

	for _, entry := range dir {

		t := Thing{
			Path: path.Join(ts.Path, "things", entry.Name()),
		}
		t.TimePath = path.Join(ts.Path, "time", fmt.Sprintf("%s.csv", strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))))

		data, err := os.ReadFile(t.Path)
		if err != nil {
			continue
		}

		rest, err := frontmatter.Parse(bytes.NewReader(data), &t)
		if err != nil {
			continue
		}

		switch ts.filter {
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

		t.Content = string(rest)

		list = append(list, t)
	}

	ts.Things = list
	ts.Sort("")

	return nil
}

func (ts *Things) Sort(s string) {

	if s != "" {
		ts.sort = s
	}

	switch ts.sort {
	case "age":
		sort.Slice(ts.Things, func(i, j int) bool {
			return ts.Things[i].Path > ts.Things[j].Path
		})
	case "priority":
		sort.Slice(ts.Things, func(i, j int) bool {
			return ts.Things[i].Priority < ts.Things[j].Priority
		})
	case "type":
		sort.Slice(ts.Things, func(i, j int) bool {
			if ts.Things[i].Type != ts.Things[j].Type {
				return ts.Things[i].Type < ts.Things[j].Type
			}
			return ts.Things[i].Priority < ts.Things[j].Priority
		})
	}

	var pinned []Thing
	var unpinned []Thing

	for _, t := range ts.Things {
		if t.Pin {
			pinned = append(pinned, t)
		} else {
			unpinned = append(unpinned, t)
		}
	}

	ts.Things = append(pinned, unpinned...)
}

func (ts *Things) Search(s string) error {

	var result []Thing

	for _, t := range ts.Things {
		contents, err := os.ReadFile(t.Path)
		if err != nil {
			return err
		}

		if strings.Contains(strings.ToLower(string(contents)), strings.ToLower(s)) {
			result = append(result, t)
		}
	}
	ts.Things = result
	ts.Sort("")

	return nil
}
