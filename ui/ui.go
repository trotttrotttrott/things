package ui

import (
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/trotttrotttrott/things/things"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	cursor int

	things things.Things

	lineNum bool
	modes   []string
	mode    int

	search struct {
		active bool
		input  textinput.Model
	}
	newType struct {
		active bool
		input  textinput.Model
	}
	viewport struct {
		width   int
		height  int
		startAt int
	}
	confirmDelete *things.Thing
	errs          []error
	helpActive    bool
}

func Start(thingsDir string) {

	m := model{
		modes: []string{"thing", "type"},
	}

	m.search.input = textinput.New()
	m.search.input.Prompt = "  Search: "

	m.newType.input = textinput.New()
	m.newType.input.Prompt = "  Type name: "

	m.things = things.New(thingsDir)

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatalln("Error:", err)
	}
}

func (m *model) searchThings() {
	err := m.things.Search(m.search.input.Value())
	m.errs = append(m.errs, err)
	m.setCursorInBounds()
}

func (m *model) searchDeactivate() {
	m.things.ResetThings()
	m.search.active = false
	m.search.input.Blur()
	m.search.input.Reset()
}

func (m *model) newTypeDeactivate() {
	m.newType.active = false
	m.newType.input.Blur()
	m.newType.input.Reset()
}

func (m *model) thingTypeKeys() (typeKeys []string) {
	for k := range m.things.Types {
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
	if m.cursor+1 > len(m.things.Things) {
		m.cursor = len(m.things.Things) - 1
	}
}

func (m *model) setCursorInView() {
	if m.cursor > m.viewportHeight()+m.viewport.startAt {
		m.viewport.startAt = m.cursor - m.viewportHeight()
	} else if m.cursor < m.viewport.startAt {
		m.viewport.startAt = m.cursor
	}
}

func (m *model) viewportHeight() int {
	h := m.viewport.height
	if m.search.active {
		h -= 3
	}
	if m.newType.active {
		h -= 3
	}
	if len(m.things.Things) > 0 && m.things.Things[0].Pin {
		h -= 1
	}
	return h
}

func (m *model) hasDeepDir(t things.Thing) bool {
	basename := filepath.Base(t.Path)
	thingID := basename[:len(basename)-3]
	deepDir := filepath.Join(m.things.Path, "things-deep", thingID)

	info, err := os.Stat(deepDir)
	if err != nil {
		return false
	}
	return info.IsDir()
}
