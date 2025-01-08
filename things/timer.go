package things

import (
	"encoding/csv"
	"os"
	"time"
)

var timer struct {
	start    time.Time
	fileName string
}

func Start(fileName string) {
	timer.start = time.Now().UTC()
	timer.fileName = fileName
}

func Stop() error {
	end := time.Now().UTC()

	f, err := os.OpenFile(timer.fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	w := csv.NewWriter(f)
	w.WriteAll([][]string{
		{timer.start.Format(time.RFC3339), end.Format(time.RFC3339)},
	})

	return w.Error()
}
