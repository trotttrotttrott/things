package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type editorFinishedMsg struct{ err error }

func newThing(fileName string, thingTypeKeys []string) tea.Cmd {

	fpath := path.Join(thingsDir, "things")
	fname := filepath.Join(fpath, fmt.Sprintf("%s.md", fileName))

	f, err := os.Create(fname)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	_, err = f.WriteString(strings.Join(
		[]string{
			"---",
			"title: Thing",
			fmt.Sprintf("type: # %s", strings.Join(thingTypeKeys, " ")),
			"priority: 0",
			"---",
			"",
		}, "\n"))
	if err != nil {
		log.Fatalln("Error:", err)
	}

	f.Sync()

	defer f.Close()

	e := os.Getenv("EDITOR")

	if e == "" {
		e = "vim"
	}

	cmd := exec.Command(e, fname)

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return editorFinishedMsg{err}
	})
}

func editThing(t thing) tea.Cmd {

	e := os.Getenv("EDITOR")

	if e == "" {
		e = "vim"
	}

	cmd := exec.Command(e, t.path)

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return editorFinishedMsg{err}
	})
}
