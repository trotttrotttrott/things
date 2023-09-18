package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	thingsDir  string
	thingTypes map[string]thingType
)

type model struct {
	cursor int

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

		case "enter":
			t := m.things[m.cursor]
			b := filepath.Base(t.path)
			timeThing(strings.TrimSuffix(b, filepath.Ext(b)))
			return m, editThing(t)

		case "n":
			t := time.Now().UTC()
			fileName := t.Format("20060102150405")
			timeThing(fileName)
			return m, newThing(fileName)

		}

	case editorFinishedMsg:
		stopThingTime()
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

		ttt, ttp, tpr := t.Title, t.Type, fmt.Sprintf("%d ", t.Priority)
		if len(t.Title) > 50 {
			ttt = fmt.Sprintf("%s...", t.Title[0:47])
		}
		if len(t.Type) > 15 {
			ttp = fmt.Sprintf("%s...", t.Type[0:12])
		}
		if len(tpr) > 5 {
			tpr = fmt.Sprintf("%s+", tpr[0:4])
		}

		s += fmt.Sprintf("%s ", cursor)
		s += lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.thingType().Color)).
			Render(fmt.Sprintf("%-50s | %-15v | %5v| %s", ttt, ttp, tpr, timeSpentOnThing(t.path)))
		s += "\n"
	}

	return s
}
