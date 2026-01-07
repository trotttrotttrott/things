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

// priorityGroup returns the group number for a priority value.
// Priorities 0-4 are individual groups, 5+ are all grouped together as group 5
func priorityGroup(p int) int {
	if p < 0 {
		return 5
	}
	if p <= 4 {
		return p
	}
	return 5
}

func (m model) thingView() string {

	s := ""

	if m.search.active {
		s += fmt.Sprintf("\n%s\n\n", m.search.input.View())
	}

	if len(m.things.Things) == 0 {
		s += lipgloss.NewStyle().Faint(true).Render("  No things to show")
	}

	// Calculate fixed widths for non-title fields
	// Format: "> TITLE | TYPE | PRIORITY| AGEd | TIME*"
	cursorWidth := 2 // "> "
	lineNumWidth := 0
	if m.lineNum {
		lineNumWidth = len(fmt.Sprintf("%v ", len(m.things.Things)))
	}

	maxPriorityLen := m.maxPriorityLen()
	maxTypeLen := m.maxTypeLen()

	// Calculate max time string width across all things
	maxTimeWidth := 0
	for _, t := range m.things.Things {
		timeStr := t.TimeString()
		if len(timeStr) > maxTimeWidth {
			maxTimeWidth = len(timeStr)
		}
	}

	// Calculate the fixed width for other fields
	// " | TYPE | PRIORITY| AGEd | TIME*"
	// " | " + TYPE + " | " + PRIORITY + "| " + AGE(3) + "d | " + TIME + "*"(maybe)
	otherFieldsWidth := 3 + maxTypeLen + 3 + maxPriorityLen + 2 + 3 + 4 + maxTimeWidth + 2

	// Calculate available space for title+note
	availableWidth := m.viewport.width - cursorWidth - lineNumWidth - otherFieldsWidth
	if availableWidth < 10 {
		availableWidth = 10 // Minimum reasonable width
	}
	maxTitleLen := availableWidth

	// First pass: determine if we need truncation mode or compact mode
	needsTruncation := false
	longestLen := 0
	for _, t := range m.things.Things {
		displayTitle := t.Title
		noteText := ""
		if t.Note != "" {
			noteText = " | " + t.Note
		}
		fullLen := len(displayTitle) + len(noteText)

		if fullLen > maxTitleLen {
			needsTruncation = true
			break
		}
		if fullLen > longestLen {
			longestLen = fullLen
		}
	}

	// Determine padding length: use maxTitleLen if truncation needed, otherwise longest actual length
	paddingLen := maxTitleLen
	if !needsTruncation {
		paddingLen = longestLen
	}

	linesRendered := 0
	lastRenderedGroup := -1

	for i, t := range m.things.Things {

		currentGroup := priorityGroup(t.Priority)

		if i < m.viewport.startAt {
			continue
		}

		// Check if we need a blank line separator
		needsSeparator := lastRenderedGroup != -1 && currentGroup != lastRenderedGroup

		// Check if we have room for both the separator (if needed) and the item
		linesNeeded := 1
		if needsSeparator {
			linesNeeded = 2
		}
		if linesRendered+linesNeeded > m.viewportHeight() {
			return s
		}

		// Show blank line when priority group changes (only if we rendered the previous group)
		if needsSeparator {
			s += "\n"
			linesRendered++
		}

		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s ", cursor)

		if m.lineNum {
			numWidth := len(fmt.Sprintf("%v", len(m.things.Things)))
			s += fmt.Sprintf("%*v ", numWidth, i+1)
		}

		// Calculate combined title + note for truncation
		displayTitle := t.Title
		noteText := ""
		if t.Note != "" {
			noteText = " | " + t.Note
		}
		fullLen := len(displayTitle) + len(noteText)

		// Handle truncation treating title + note as one string
		ttt := t.Title
		notePart := noteText
		if fullLen > maxTitleLen {
			// Truncate the combined string
			full := displayTitle + noteText
			truncated := full[0:maxTitleLen-3] + "..."

			// Check if the note delimiter position is within the truncated portion
			// Use the actual title length rather than searching for " | " to handle titles containing pipes
			titleLen := len(displayTitle)
			truncatedLen := maxTitleLen - 3 // Length before adding "..."

			if titleLen < truncatedLen {
				// Delimiter is in the truncated portion, split at actual title boundary
				ttt = truncated[0:titleLen]
				notePart = truncated[titleLen:]
			} else {
				// Truncation happened within the title, no note shown
				ttt = truncated
				notePart = ""
			}
		} else {
			// Pad to paddingLen (either maxTitleLen or longest actual length)
			padding := paddingLen - fullLen
			notePart = noteText + strings.Repeat(" ", padding)
		}

		// Render title with main style
		titleStyled := lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.things.Types[t.Type].Color)).
			Faint(t.Pause).
			Render(ttt)

		// Render note with faint style (always faint)
		noteStyled := ""
		if notePart != "" {
			noteStyled = lipgloss.NewStyle().
				Foreground(lipgloss.Color(m.things.Types[t.Type].Color)).
				Faint(true).
				Render(notePart)
		}

		titleFormatted := titleStyled + noteStyled

		ttp, tpr := t.Type, fmt.Sprintf("%d ", t.Priority)

		maxPriorityLen := m.maxPriorityLen()
		if len(tpr) > maxPriorityLen {
			tpr = fmt.Sprintf("%s+", tpr[0:maxPriorityLen-1])
		}

		deepIndicator := ""
		if m.hasDeepDir(t) {
			deepIndicator = " *"
		}

		// Apply the same styling to the rest of the row
		rowStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.things.Types[t.Type].Color)).
			Faint(t.Pause)

		restOfRow := rowStyle.Render(fmt.Sprintf(" | %-*v | %*v| %*sd | %s%s", m.maxTypeLen(), ttp, maxPriorityLen, tpr, 3, t.Age(), t.TimeString(), deepIndicator))

		s += titleFormatted + restOfRow
		s += "\n"
		linesRendered++
		lastRenderedGroup = currentGroup

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
	s += "  " + key.Render("P") + "       pause: true\n\n"

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
