package main

import (
	"encoding/csv"
	"os"
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

func stopThingTime() error {
	end := time.Now().UTC()

	f, err := os.OpenFile(thingTime.fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	w := csv.NewWriter(f)
	w.WriteAll([][]string{
		{thingTime.start.Format(time.RFC3339), end.Format(time.RFC3339)},
	})

	return w.Error()
}
