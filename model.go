package main

import (
	"sort"

	"github.com/charmbracelet/bubbles/textinput"
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
	search     struct {
		active bool
		input  textinput.Model
	}
	viewport struct {
		width   int
		height  int
		startAt int
	}
	confirmDelete *thing
	errs          []error
}

func (m *model) searchThings() {
	things, err := thingSearch(things(m.filter), m.search.input.Value())
	m.errs = append(m.errs, err)
	m.things = things
	m.setCursorInBounds()
}

func (m *model) searchDeactivate() {
	m.search.active = false
	m.search.input.Blur()
	m.search.input.Reset()
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
	return h
}
