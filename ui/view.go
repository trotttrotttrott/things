package ui

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

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

	if m.helpActive {
		return m.helpView()
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

	if m.search.active {
		s += fmt.Sprintf("\n%s\n\n", m.search.input.View())
	}

	if len(m.things.Things) == 0 {
		s += lipgloss.NewStyle().Faint(true).Render("  No things to show")
	}

	for i, t := range m.things.Things {

		if i < m.viewport.startAt {
			continue
		}
		if i > m.viewportHeight()+m.viewport.startAt {
			return s
		}

		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s ", cursor)

		maxTitleLen := 35
		switch {
		case m.viewport.width > 90:
			maxTitleLen = 55
		case m.viewport.width > 80:
			maxTitleLen = 45
		}
		if m.lineNum {
			numWidth := len(fmt.Sprintf("%v", len(m.things.Things)))
			maxTitleLen = maxTitleLen - numWidth - 1
			s += fmt.Sprintf("%*v ", numWidth, i+1)
		}
		ttt, ttp, tpr := t.Title, t.Type, fmt.Sprintf("%d ", t.Priority)
		if len(t.Title) > maxTitleLen {
			ttt = fmt.Sprintf("%s...", t.Title[0:maxTitleLen-3])
		}

		maxPriorityLen := 5
		if len(tpr) > maxPriorityLen {
			tpr = fmt.Sprintf("%s+", tpr[0:maxPriorityLen-1])
		}

		deepIndicator := ""
		if m.hasDeepDir(t) {
			deepIndicator = " *"
		}

		s += lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.things.Types[t.Type].Color)).
			Faint(t.Pause).
			Bold(t.Today).
			Render(fmt.Sprintf("%-*s | %-*v | %*v| %*sd | %s%s", maxTitleLen, ttt, m.maxTypeLen(), ttp, maxPriorityLen, tpr, 3, t.Age(), t.TimeString(), deepIndicator))
		s += "\n"

		if t.Pin && len(m.things.Things) > i-1 && !m.things.Things[i+1].Pin {
			s += "  ---\n"
		}

	}

	return s
}

func (m model) typeView() string {

	s := ""

	if m.newType.active {
		s += fmt.Sprintf("\n%s\n\n", m.newType.input.View())
	}

	for i, t := range m.thingTypeKeys() {

		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s ", cursor)

		if m.lineNum {
			numWidth := len(fmt.Sprintf("%v", len(m.things.Things)))
			s += fmt.Sprintf("%*v ", numWidth, i+1)
		}

		description := regexp.MustCompile(`\n+`).ReplaceAllString(strings.TrimSpace(m.things.Types[t].Description), "...")
		if len(description) > 50 {
			description = fmt.Sprintf("%s...", description[0:50])
		}

		s += lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.things.Types[t].Color)).
			Render(fmt.Sprintf("%-*s | %s", m.maxTypeLen(), t, description))
		s += "\n"

	}

	return s
}

func (m model) helpView() string {
	section := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffff00")).
		Bold(true)

	key := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00ff00"))

	s := section.Render("  Mode") + "\n"
	s += "  " + key.Render(">") + "       switch between \"thing\" and \"type\" modes\n"
	s += "  " + key.Render("/") + "       search (\"thing\" mode only)\n"
	s += "  " + key.Render("?") + "       toggle help\n\n"

	s += section.Render("  Navigation") + "\n"
	s += "  " + key.Render("k") + "       cursor up\n"
	s += "  " + key.Render("j") + "       cursor down\n"
	s += "  " + key.Render("ctrl+u") + "  cursor up 5\n"
	s += "  " + key.Render("ctrl+d") + "  cursor down 5\n"
	s += "  " + key.Render("g") + "       set cursor to first\n"
	s += "  " + key.Render("G") + "       set cursor to last\n\n"

	s += section.Render("  Filter") + " (\"thing\" mode only)\n"
	s += "  " + key.Render("C") + "       current, done: false (default)\n"
	s += "  " + key.Render("D") + "       done: true\n"
	s += "  " + key.Render("A") + "       all, no filter\n"
	s += "  " + key.Render("P") + "       pause: true\n"
	s += "  " + key.Render("T") + "       today: true\n\n"

	s += section.Render("  Sort") + " (\"thing\" mode only)\n"
	s += "  " + key.Render("a") + "       sort things by age\n"
	s += "  " + key.Render("p") + "       sort things by priority (default)\n"
	s += "  " + key.Render("t") + "       sort things by type and priority\n\n"

	s += section.Render("  Display") + "\n"
	s += "  " + key.Render("#") + "       toggle line numbers\n\n"

	s += section.Render("  Edit") + "\n"
	s += "  " + key.Render("n") + "       create new thing or type\n"
	s += "  " + key.Render("enter") + "   open thing or type in $EDITOR\n"
	s += "  " + key.Render("E") + "       open thing directory in $EDITOR (\"thing\" mode only)\n"
	s += "  " + key.Render("ctrl+e") + "  open thing time file in $EDITOR (\"thing\" mode only)\n"
	s += "  " + key.Render("ctrl+x") + "  delete thing (\"thing\" mode only)\n\n"

	s += section.Render("  Quit") + "\n"
	s += "  " + key.Render("ctrl+c") + "  quit\n"
	s += "  " + key.Render("q") + "       quit\n\n"

	s += lipgloss.NewStyle().Faint(true).Render("  Press ? again to close help\n")

	return s
}
