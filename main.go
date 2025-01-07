package main

import (
	"log"
	"os"
	"path"

	"github.com/trotttrotttrott/things/things"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	// TODO: remove this
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
		modes: []string{"thing", "type"},
	}

	m.search.input = textinput.New()
	m.search.input.Prompt = "  Search: "

	m.things = things.New(thingsDir)

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatalln("Error:", err)
	}
}
