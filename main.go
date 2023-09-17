package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/adrg/frontmatter"
	tea "github.com/charmbracelet/bubbletea"
)

var thingsDir string

type thing struct {
	Title    string
	Type     string
	Priority int
	Done     bool
	content  string
}

type thingType struct {
	description string
	path        string
	color       string
}

type model struct {
	things []thing

	cursor   int
	selected *int

	err error
}

func main() {

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("Error:", err)
	}

	thingsDir = path.Join(home, ".things")

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalln("Error:", err)
	}
}

func initialModel() model {

	dir, err := os.ReadDir(path.Join(thingsDir, "things"))
	if err != nil {
		log.Fatalln("Error:", err)
	}

	var things []thing

	for _, entry := range dir {

		t := thing{}

		data, err := os.ReadFile(path.Join(thingsDir, "things", entry.Name()))
		if err != nil {
			log.Fatalln("Error:", err)
		}

		rest, err := frontmatter.Parse(bytes.NewReader(data), &t)
		if err != nil {
			log.Fatalln("Error:", err)
		}

		t.content = string(rest)

		things = append(things, t)
	}

	sort.Slice(things, func(i, j int) bool {
		return things[i].Priority < things[j].Priority
	})

	return model{
		things: things,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.things)-1 {
				m.cursor++
			}

		case "enter", " ":
			m.selected = &m.cursor

		case "esc":
			m.selected = nil

		}
	}

	return m, nil
}

func (m model) View() string {

	s := ""

	for i, t := range m.things {

		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if m.selected != nil && *m.selected == i {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s %v %v %v %s\n", cursor, checked, t.Title, t.Type, t.Priority, t.Done, strings.TrimSpace(t.content))
	}

	return s
}
