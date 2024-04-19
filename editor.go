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

func editThingTime(t thing) tea.Cmd {

	cmd := exec.Command(editor(), t.timePath)

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return editThingTimeFinishedMsg{err}
	})
}

func editType(key string) tea.Cmd {

	p := path.Join(thingsDir, "types", fmt.Sprintf("%s.md", key))

	cmd := exec.Command(editor(), p)

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return editTypeFinishedMsg{err}
	})
}
