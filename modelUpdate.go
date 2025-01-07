package main

import (
	"github.com/trotttrotttrott/things/things"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	m.errs = m.errs[:0]

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.viewport.width = msg.Width
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
			m.errs = append(m.errs, m.confirmDelete.Remove())
			m.confirmDelete = nil
			m.setCursorInBounds()
			return m, nil
		} else {
			m.confirmDelete = nil
		}

		var cursorLimit int
		switch m.modes[m.mode] {
		case "thing":
			cursorLimit = len(m.things.Things)
		case "type":
			cursorLimit = len(m.things.Types)
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
			m.things.Filter("current")
			m.searchDeactivate()
		case "D":
			m.things.Filter("done")
			m.searchDeactivate()
		case "A":
			m.things.Filter("")
			m.searchDeactivate()
		case "P":
			m.things.Filter("pause")
			m.searchDeactivate()
		case "T":
			m.things.Filter("today")
			m.searchDeactivate()

		// sort
		case "a":
			m.things.Sort("age")
		case "p":
			m.things.Sort("priority")
		case "t":
			m.things.Sort("type")

		// display
		case "#":
			m.lineNum = !m.lineNum

		// edit
		case "n":
			if m.modes[m.mode] == "thing" {
				t, err := m.things.NewThing(m.thingTypeKeys())
				if err != nil {
					m.errs = append(m.errs, err)
				} else {
					things.Start(t.TimePath)
					return m, editThing(t.Path)
				}
			}
		case "enter":
			switch m.modes[m.mode] {
			case "thing":
				t := m.things.Things[m.cursor]
				things.Start(t.TimePath)
				return m, editThing(t.Path)
			case "type":
				return m, editType(m.thingTypeKeys()[m.cursor])
			}
		case "ctrl+e":
			if m.modes[m.mode] == "thing" {
				t := m.things.Things[m.cursor]
				return m, editThingTime(t)
			}
		case "ctrl+x":
			if m.modes[m.mode] == "thing" {
				t := m.things.Things[m.cursor]
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
		m.errs = append(m.errs, things.Stop())
		m.errs = append(m.errs, msg.err)
		if m.search.active {
			m.searchThings()
		}

	case editThingTimeFinishedMsg:
		m.errs = append(m.errs, msg.err)

	case editTypeFinishedMsg:
		m.errs = append(m.errs, msg.err)
		m.things.ResetTypes()
	}

	m.setCursorInView()

	return m, nil
}
