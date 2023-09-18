package main

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
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
		{thingTime.start.Format(time.RFC3339), end.Format(time.RFC3339)},
	})
	if err := w.Error(); err != nil {
		log.Fatalln("Error:", err)
	}
}

func timeSpentOnThing(thingPath string) string {

	var timeSpent time.Duration

	b := filepath.Base(thingPath)
	fpath := filepath.Join(thingsDir, "time", fmt.Sprintf("%s.csv", strings.TrimSuffix(b, filepath.Ext(b))))
	if _, err := os.Stat(fpath); errors.Is(err, os.ErrNotExist) {
		return timeSpent.String()
	}
	data, err := os.ReadFile(fpath)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	rdr := csv.NewReader(bytes.NewReader(data))
	records, err := rdr.ReadAll()
	if err != nil {
		log.Fatalln("Error:", err)
	}

	for _, r := range records {
		start, err := time.Parse(time.RFC3339, r[0])
		if err != nil {
			log.Fatalln("Error:", err)
		}
		end, err := time.Parse(time.RFC3339, r[1])
		if err != nil {
			log.Fatalln("Error:", err)
		}
		timeSpent += end.Sub(start)
	}

	return timeSpent.String()
}
