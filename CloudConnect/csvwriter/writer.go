package main

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"
)

type writer struct {
	writeDir   string
	currFile   *os.File
	currWriter *csv.Writer
	writeUntil time.Time
}

func (w *writer) Write(record []string) error {
	tn := time.Now()

	if w.currWriter != nil && tn.Before(w.writeUntil) {
		w.currWriter.Write(record)
		w.currWriter.Flush() // this ensures that the file modtime is continously updated
		return nil
	}

	if w.currFile != nil {
		if err := w.currFile.Close(); err != nil {
			return err
		}
	}
	wu := writeUntil(tn)
	f, err := os.Create(w.writeDir + strconv.FormatInt(wu.Unix(), 10) + ".csv")
	if err != nil {
		return err
	}
	w.writeUntil = wu
	w.currFile = f
	writer := csv.NewWriter(f)
	w.currWriter = writer
	w.currWriter.Write(record)
	w.currWriter.Flush()
	return nil
}

func writeUntil(t time.Time) time.Time {
	// this returns "the next minute", which should make sense for most cases.
	// We can make this configurable in future if needed
	rounded := t.Round(time.Minute)
	if rounded.Before(t) {
		rounded = rounded.Add(time.Minute)
	}
	return rounded
}
