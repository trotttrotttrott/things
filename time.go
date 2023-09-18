package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"
)

var thingTime struct {
	start    time.Time
	fileName string
}

func timeThing(fileName string) {
	thingTime.start = time.Now().UTC()
	thingTime.fileName = fileName
}

func stopThingTime() {
	end := time.Now().UTC()

	fpath := path.Join(thingsDir, "time")
	fname := filepath.Join(fpath, fmt.Sprintf("%s.csv", thingTime.fileName))

	f, err := os.OpenFile(fname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	w := csv.NewWriter(f)
	w.WriteAll([][]string{
		{thingTime.start.String(), end.String()},
	})
	if err := w.Error(); err != nil {
		log.Fatalln("Error:", err)
	}
}
