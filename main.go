package main

import (
	"fmt"
	"log"
	"os"
	"path"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	thingsDir  string
	thingTypes map[string]thingType
)

type model struct {
	cursor   int
	selected *int

	things []thing

	err error
}

func main() {

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("Error:", err)
	}

	thingsDir = path.Join(home, ".things")

	typesInit()

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalln("Error:", err)
	}
}

func initialModel() model {
	return model{
		things: things(),
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

		case "g":
			m.cursor = 0

		case "G":
			m.cursor = len(m.things) - 1

		case "enter", " ":
			m.selected = &m.cursor

		case "esc":
			m.selected = nil

		case "n":
			return m, newThing()

		case "e":
			return m, editThing(m.things[m.cursor])
		}

	case editorFinishedMsg:
		if msg.err != nil {
			m.err = msg.err
		}
		m.things = things()
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

		var style = lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.thingType().Color))

		s += fmt.Sprintf("%s [%s] ", cursor, checked)
		s += style.Render(fmt.Sprintf("%s %v %v %v", t.Title, t.Type, t.Priority, t.Done))
		s += "\n"
	}

	return s
}
