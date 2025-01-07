package main

import (
	"log"
	"os"
	"path"

	"github.com/trotttrotttrott/things/ui"
)

func main() {

	thingsDir := os.Getenv("THINGS_DIR")

	if thingsDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalln("Error:", err)
		}
		thingsDir = path.Join(home, ".things")
	}

	ui.Start(thingsDir)
}
