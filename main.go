package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"slices"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	thingsDir string
)

type model struct {
	cursor     int
	things     []thing
	sort       string
	filter     string
	thingTypes map[string]thingType
	lineNum    bool
	modes      []string
	mode       int
	viewport   struct {
		height  int
		startAt int
	}
	confirmDelete *thing
	errs          []error
}

func main() {

	thingsDir = os.Getenv("THINGS_DIR")

	if thingsDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalln("Error:", err)
		}
		thingsDir = path.Join(home, ".things")
	}

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalln("Error:", err)
	}
}

func initialModel() model {

	m := model{
		things:     things(""),
		sort:       "priority",
		thingTypes: thingTypes(),
		modes:      []string{"thing", "type"},
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

func (m *model) thingTypeKeys() (typeKeys []string) {
	for k := range m.thingTypes {
		typeKeys = append(typeKeys, k)
	}
	sort.Strings(typeKeys)
	return
}

func (m *model) maxTypeLen() (mx int) {
	for _, t := range m.thingTypeKeys() {
		if len(t) > mx {
			mx = len(t)
		}
	}
	return
}

func (m *model) setCursorInBounds() {
	if m.cursor+1 > len(m.things) {
		m.cursor = len(m.things) - 1
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	m.errs = m.errs[:0]

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.viewport.height = msg.Height - 2

	case tea.KeyMsg:

		if m.confirmDelete != nil && msg.String() == "enter" {
			m.errs = append(m.errs, m.confirmDelete.remove())
			m.confirmDelete = nil
			m.things = things(m.filter)
			m.sortThings()
			m.setCursorInBounds()
			return m, nil
		} else {
			m.confirmDelete = nil
		}

		var cursorLimit int
		switch m.modes[m.mode] {
		case "thing":
			cursorLimit = len(m.things)
		case "type":
			cursorLimit = len(m.thingTypes)
		}

		switch msg.String() {

		// navigation
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < cursorLimit-1 {
				m.cursor++
			}
		case "ctrl+u":
			if m.cursor-5 > 0 {
				m.cursor -= 5
			} else {
				m.cursor = 0
			}
		case "ctrl+d":
			if m.cursor+5 < cursorLimit-1 {
				m.cursor += 5
			} else {
				m.cursor = cursorLimit - 1
			}
		case "g":
			m.cursor = 0
		case "G":
			m.cursor = cursorLimit - 1

		// modes
		case ">":
			m.mode = (m.mode + 1) % len(m.modes)
			m.cursor = 0

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
			if m.modes[m.mode] == "thing" {
				t, err := thingNew(m.thingTypeKeys())
				if err != nil {
					m.errs = append(m.errs, err)
				} else {
					timeThing(t.timePath)
					return m, editThing(t.path)
				}
			}
		case "enter":
			switch m.modes[m.mode] {
			case "thing":
				t := m.things[m.cursor]
				timeThing(t.timePath)
				return m, editThing(t.path)
			case "type":
				return m, editType(m.thingTypeKeys()[m.cursor])
			}
		case "ctrl+e":
			if m.modes[m.mode] == "thing" {
				t := m.things[m.cursor]
				return m, editThingTime(t)
			}
		case "ctrl+x":
			if m.modes[m.mode] == "thing" {
				t := m.things[m.cursor]
				m.confirmDelete = &t
			}

		// quit
		case "ctrl+c", "q":
			return m, tea.Quit

		}

	case editThingFinishedMsg:
		m.errs = append(m.errs, stopThingTime())
		m.errs = append(m.errs, msg.err)
		m.things = things(m.filter)
		m.sortThings()
		m.setCursorInBounds()

	case editThingTimeFinishedMsg:
		m.errs = append(m.errs, msg.err)

	case editTypeFinishedMsg:
		m.errs = append(m.errs, msg.err)
		m.thingTypes = thingTypes()
	}

	// ensure cursor is in view
	if m.cursor > m.viewport.height+m.viewport.startAt {
		m.viewport.startAt = m.cursor - m.viewport.height
	} else if m.cursor < m.viewport.startAt {
		m.viewport.startAt = m.cursor
	}

	return m, nil
}

func (m model) View() string {

	m.errs = slices.DeleteFunc(
		m.errs,
		func(err error) bool {
			return err == nil
		},
	)
	if len(m.errs) > 0 {
		return m.errorView()
	}

	if m.confirmDelete != nil {
		return m.confirmDeleteView()
	}

	switch m.modes[m.mode] {
	case "thing":
		return m.thingView()
	case "type":
		return m.typeView()
	}

	m.errs = append(m.errs, fmt.Errorf("No view found for model state"))
	return m.errorView()
}

func (m model) errorView() string {
	return lipgloss.
		NewStyle().
		Foreground(lipgloss.Color("#ff0000")).
		Render(fmt.Sprint(m.errs))
}

func (m model) confirmDeleteView() string {
	return lipgloss.
		NewStyle().
		Foreground(lipgloss.Color("#ff0000")).
		Render(fmt.Sprintf("Delete %q? [press enter to confirm]", m.confirmDelete.Title))
}

func (m model) thingView() string {

	s := ""

	for i, t := range m.things {

		if i < m.viewport.startAt {
			continue
		}
		if i > m.viewport.height+m.viewport.startAt {
			return s
		}

		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s ", cursor)

		maxTitleLen, maxPriorityLen := 30, 5
		if m.lineNum {
			numWidth := len(fmt.Sprintf("%v", len(m.things)))
			maxTitleLen = maxTitleLen - numWidth - 1
			s += fmt.Sprintf("%*v ", numWidth, i+1)
		}

		ttt, ttp, tpr := t.Title, t.Type, fmt.Sprintf("%d ", t.Priority)
		if len(t.Title) > maxTitleLen {
			ttt = fmt.Sprintf("%s...", t.Title[0:maxTitleLen-3])
		}
		if len(tpr) > maxPriorityLen {
			tpr = fmt.Sprintf("%s+", tpr[0:maxPriorityLen-1])
		}

		s += lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.thingTypes[t.Type].Color)).
			Faint(t.Pause).
			Bold(t.Today).
			Render(fmt.Sprintf("%-*s | %-*v | %*v| %*sd | %s", maxTitleLen, ttt, m.maxTypeLen(), ttp, maxPriorityLen, tpr, 3, t.age(), t.time().String()))
		s += "\n"
	}

	return s
}

func (m model) typeView() string {

	s := ""

	for i, t := range m.thingTypeKeys() {

		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s ", cursor)

		if m.lineNum {
			numWidth := len(fmt.Sprintf("%v", len(m.things)))
			s += fmt.Sprintf("%*v ", numWidth, i+1)
		}

		description := regexp.MustCompile(`\n+`).ReplaceAllString(strings.TrimSpace(m.thingTypes[t].description), "...")
		if len(description) > 50 {
			description = fmt.Sprintf("%s...", description[0:50])
		}

		s += lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.thingTypes[t].Color)).
			Render(fmt.Sprintf("%-*s | %s", m.maxTypeLen(), t, description))
		s += "\n"

	}

	return s
}
