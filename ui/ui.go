package ui

import (
	"fmt"
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
	m.things.ResetThings()
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

func (m *model) maxPriorityLen() int {
	maxPriority := 0
	for _, t := range m.things.Things {
		if t.Priority > maxPriority {
			maxPriority = t.Priority
		}
	}
	return len(fmt.Sprintf("%d ", maxPriority))
}

func (m *model) setCursorInBounds() {
	if m.cursor+1 > len(m.things.Things) {
		m.cursor = len(m.things.Things) - 1
	}
}

// countSeparatorLines counts how many blank separator lines exist between two item indices
// This matches the rendering logic: no separator before the first rendered item
func (m *model) countSeparatorLines(start, end int) int {
	if start >= end || len(m.things.Things) == 0 {
		return 0
	}

	count := 0
	lastRenderedGroup := -1

	for i := start; i <= end && i < len(m.things.Things); i++ {
		currentGroup := priorityGroup(m.things.Things[i].Priority)
		// Only count separator if we've rendered a previous item (lastRenderedGroup != -1)
		// and the group changed
		if lastRenderedGroup != -1 && currentGroup != lastRenderedGroup {
			count++
		}
		lastRenderedGroup = currentGroup
	}

	return count
}

func (m *model) setCursorInView() {
	if m.cursor < m.viewport.startAt {
		// Scrolling up
		m.viewport.startAt = m.cursor
		return
	}

	// Check if cursor is within visible range (accounting for separators)
	for {
		separators := m.countSeparatorLines(m.viewport.startAt, m.cursor)
		visibleItems := m.viewportHeight() - separators

		if m.cursor <= m.viewport.startAt+visibleItems-1 {
			// Cursor is visible
			break
		}

		// Need to scroll down
		m.viewport.startAt++
		if m.viewport.startAt > m.cursor {
			break
		}
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
