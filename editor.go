package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	tea "github.com/charmbracelet/bubbletea"
)

type editThingFinishedMsg struct{ err error }
type editThingTimeFinishedMsg struct{ err error }
type editTypeFinishedMsg struct{ err error }

func editThing(thingPath string) tea.Cmd {

	e := os.Getenv("EDITOR")

	if e == "" {
		e = "vim"
	}

	cmd := exec.Command(e, thingPath)

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return editThingFinishedMsg{err}
	})
}

func editThingTime(t thing) tea.Cmd {

	e := os.Getenv("EDITOR")

	if e == "" {
		e = "vim"
	}

	cmd := exec.Command(e, t.timePath)

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return editThingTimeFinishedMsg{err}
	})
}

func editType(key string) tea.Cmd {

	p := path.Join(thingsDir, "types", fmt.Sprintf("%s.md", key))

	e := os.Getenv("EDITOR")

	if e == "" {
		e = "vim"
	}

	cmd := exec.Command(e, p)

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return editTypeFinishedMsg{err}
	})
}
