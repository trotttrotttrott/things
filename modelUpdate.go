package main

import tea "github.com/charmbracelet/bubbletea"

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	m.errs = m.errs[:0]

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.viewport.height = msg.Height - 2

	case tea.KeyMsg:

		switch msg.String() {

		// quit
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			if !m.search.input.Focused() {
				return m, tea.Quit
			}
		}

		if m.search.active {

			switch msg.String() {

			case "esc":
				m.searchDeactivate()
				m.things = things(m.filter)

			case "enter":
				if m.search.input.Focused() {
					m.search.input.Blur()
					m.searchThings()
					m.cursor = 0
					return m, nil
				}

			default:
				if m.search.input.Focused() {
					var cmd tea.Cmd
					m.search.input, cmd = m.search.input.Update(msg)
					return m, cmd
				}
			}
		}

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
		case "C":
			m.filter = "current"
			m.filterThings()
			m.searchDeactivate()
		case "D":
			m.filter = "done"
			m.filterThings()
			m.searchDeactivate()
		case "A":
			m.filter = ""
			m.filterThings()
			m.searchDeactivate()
		case "P":
			m.filter = "pause"
			m.filterThings()
			m.searchDeactivate()
		case "T":
			m.filter = "today"
			m.filterThings()
			m.searchDeactivate()

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

		// search
		case "/":
			if m.modes[m.mode] == "thing" {
				m.search.active = true
				m.search.input.Focus()

			}
		}

	case editThingFinishedMsg:
		m.errs = append(m.errs, stopThingTime())
		m.errs = append(m.errs, msg.err)
		m.things = things(m.filter)
		m.sortThings()
		if m.search.active {
			m.searchThings()
		}

	case editThingTimeFinishedMsg:
		m.errs = append(m.errs, msg.err)

	case editTypeFinishedMsg:
		m.errs = append(m.errs, msg.err)
		m.thingTypes = thingTypes()
	}

	// ensure cursor is in view
	if m.cursor > m.viewportHeight()+m.viewport.startAt {
		m.viewport.startAt = m.cursor - m.viewportHeight()
	} else if m.cursor < m.viewport.startAt {
		m.viewport.startAt = m.cursor
	}

	return m, nil
}
