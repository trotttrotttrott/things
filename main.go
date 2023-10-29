package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
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
	cursor  int
	things  []thing
	sort    string
	filter  string
	lineNum bool
	err     error
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
	m := model{
		things: things(""),
		sort:   "priority",
	}
	m.sortThings()
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) sortThings() {
	switch m.sort {
	case "age":
		sort.Slice(m.things, func(i, j int) bool {
			return m.things[i].path > m.things[j].path
		})
	case "priority":
		sort.Slice(m.things, func(i, j int) bool {
			return m.things[i].Priority < m.things[j].Priority
		})
	case "type":
		sort.Slice(m.things, func(i, j int) bool {
			if m.things[i].Type != m.things[j].Type {
				return m.things[i].Type < m.things[j].Type
			}
			return m.things[i].Priority < m.things[j].Priority
		})
	}
}

func (m *model) filterThings() {
	m.things = things(m.filter)
	m.sortThings()
	m.cursor = 0
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		// navigation
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.things)-1 {
				m.cursor++
			}
		case "ctrl+u":
			if m.cursor-5 > 0 {
				m.cursor -= 5
			} else {
				m.cursor = 0
			}
		case "ctrl+d":
			if m.cursor+5 < len(m.things)-1 {
				m.cursor += 5
			} else {
				m.cursor = len(m.things) - 1
			}
		case "g":
			m.cursor = 0
		case "G":
			m.cursor = len(m.things) - 1

		// filter
		case "A":
			m.filter = ""
			m.filterThings()
		case "D":
			m.filter = "done"
			m.filterThings()
		case "P":
			m.filter = "pause"
			m.filterThings()
		case "T":
			m.filter = "today"
			m.filterThings()

		// sort
		case "a":
			m.sort = "age"
			m.sortThings()
		case "p":
			m.sort = "priority"
			m.sortThings()
		case "t":
			m.sort = "type"
			m.sortThings()

		// display
		case "#":
			m.lineNum = !m.lineNum

		// edit
		case "n":
			t := time.Now().UTC()
			fileName := t.Format("20060102150405")
			timeThing(fileName)
			return m, newThing(fileName)
		case "enter":
			t := m.things[m.cursor]
			b := filepath.Base(t.path)
			timeThing(strings.TrimSuffix(b, filepath.Ext(b)))
			return m, editThing(t)

		// quit
		case "ctrl+c", "q":
			return m, tea.Quit

		}

	case editorFinishedMsg:
		stopThingTime()
		if msg.err != nil {
			m.err = msg.err
		}
		m.things = things(m.filter)
		m.sortThings()
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

		s += fmt.Sprintf("%s ", cursor)

		maxTitleLen, maxTypeLen, maxPriorityLen := 50, 15, 5
		numWidth := len(fmt.Sprintf("%v", len(m.things)))
		if m.lineNum {
			maxTitleLen = maxTitleLen - numWidth - 1
			s += fmt.Sprintf("%*v ", numWidth, i+1)
		}

		ttt, ttp, tpr := t.Title, t.Type, fmt.Sprintf("%d ", t.Priority)
		if len(t.Title) > maxTitleLen {
			ttt = fmt.Sprintf("%s...", t.Title[0:maxTitleLen-3])
		}
		if len(t.Type) > maxTypeLen {
			ttp = fmt.Sprintf("%s...", t.Type[0:maxTypeLen-3])
		}
		if len(tpr) > maxPriorityLen {
			tpr = fmt.Sprintf("%s+", tpr[0:maxPriorityLen-1])
		}

		s += lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.thingType().Color)).
			Faint(t.Pause).
			Bold(t.Today).
			Render(fmt.Sprintf("%-*s | %-*v | %*v| %*sd | %s", maxTitleLen, ttt, maxTypeLen, ttp, maxPriorityLen, tpr, 3, t.age(), timeSpentOnThing(t.path)))
		s += "\n"
	}

	return s
}
