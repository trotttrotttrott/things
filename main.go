package main

import (
	"log"
	"os"
	"path"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	thingsDir string
)

func main() {

	thingsDir = os.Getenv("THINGS_DIR")

	if thingsDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalln("Error:", err)
		}
		thingsDir = path.Join(home, ".things")
	}

	m := model{
		sort:       "priority",
		filter:     "current",
		thingTypes: thingTypes(),
		modes:      []string{"thing", "type"},
	}

	m.search.input = textinput.New()
	m.search.input.Prompt = "  Search: "

	m.things = things(m.filter)

	m.sortThings()

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatalln("Error:", err)
	}
}
