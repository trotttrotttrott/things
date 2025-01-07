package things

import (
	"encoding/csv"
	"os"
	"time"
)

var Timer struct {
	start    time.Time
	fileName string
}

func Start(fileName string) {
	Timer.start = time.Now().UTC()
	Timer.fileName = fileName
}

func Stop() error {
	end := time.Now().UTC()

	f, err := os.OpenFile(Timer.fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	w := csv.NewWriter(f)
	w.WriteAll([][]string{
		{Timer.start.Format(time.RFC3339), end.Format(time.RFC3339)},
	})

	return w.Error()
}
