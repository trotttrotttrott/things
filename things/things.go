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

	// Ensure all required directories exist
	os.MkdirAll(path.Join(p, "things"), 0755)
	os.MkdirAll(path.Join(p, "types"), 0755)
	os.MkdirAll(path.Join(p, "time"), 0755)

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
			"priority: 5",
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

// priorityGroup returns the group number for a priority value.
// Priorities 0-4 are individual groups, 5+ are all grouped together as group 5
func priorityGroup(p int) int {
	if p < 0 {
		return 5
	}
	if p <= 4 {
		return p
	}
	return 5
}

func (ts *Things) Sort(s string) {

	if s != "" {
		ts.sort = s
	}

	// Sort by priority groups (0, 1, 2, 3, 4, 5+) and within each group apply the selected sort
	sort.Slice(ts.Things, func(i, j int) bool {
		iGroup := priorityGroup(ts.Things[i].Priority)
		jGroup := priorityGroup(ts.Things[j].Priority)

		// Different priority groups: sort by group number
		if iGroup != jGroup {
			return iGroup < jGroup
		}

		// Same priority group: apply selected sort mode
		switch ts.sort {
		case "age":
			return ts.Things[i].Path > ts.Things[j].Path
		case "priority":
			// Within same group, sort by actual priority value
			if ts.Things[i].Priority != ts.Things[j].Priority {
				return ts.Things[i].Priority < ts.Things[j].Priority
			}
			return ts.Things[i].Path > ts.Things[j].Path // Secondary sort by age
		case "type":
			if ts.Things[i].Type != ts.Things[j].Type {
				return ts.Things[i].Type < ts.Things[j].Type
			}
			// Secondary sort by priority within same type
			if ts.Things[i].Priority != ts.Things[j].Priority {
				return ts.Things[i].Priority < ts.Things[j].Priority
			}
			return ts.Things[i].Path > ts.Things[j].Path // Tertiary sort by age
		}
		return false
	})
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
