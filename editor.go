package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/trotttrotttrott/things/things"

	tea "github.com/charmbracelet/bubbletea"
)

type editThingFinishedMsg struct{ err error }
type editThingTimeFinishedMsg struct{ err error }
type editTypeFinishedMsg struct{ err error }

func editor() string {

	e := os.Getenv("EDITOR")

	if e == "" {
		e = "vim"
	}

	return e
}

func editThing(thingPath string) tea.Cmd {

	cmd := exec.Command(editor(), thingPath)

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return editThingFinishedMsg{err}
	})
}

func editThingTime(t things.Thing) tea.Cmd {

	cmd := exec.Command(editor(), t.TimePath)

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return editThingTimeFinishedMsg{err}
	})
}

func editType(typePath, key string) tea.Cmd {

	p := path.Join(typePath, fmt.Sprintf("%s.md", key))

	cmd := exec.Command(editor(), p)

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return editTypeFinishedMsg{err}
	})
}
