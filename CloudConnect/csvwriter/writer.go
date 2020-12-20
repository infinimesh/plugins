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

	// write to existing file
	if w.currWriter != nil && tn.Before(w.writeUntil) {
		w.currWriter.Write(record)
		return nil
	}

	// clean up previous file
	if w.currWriter != nil {
		w.currWriter.Flush()
	}
	if w.currFile != nil {
		if err := w.currFile.Close(); err != nil {
			return err
		}
	}

	// write to new file
	wu := writeUntil(tn)
	f, err := os.Create(w.writeDir + strconv.FormatInt(wu.Unix(), 10) + ".csv")
	if err != nil {
		return err
	}
	w.writeUntil = wu
	w.currFile = f
	writer := csv.NewWriter(f)
	w.currWriter = writer
	writer.Write(record)
	return nil
}

func writeUntil(t time.Time) time.Time {
	// this currently returns "the next minute", which should make sense for
	// most cases. We can make this configurable in future if needed
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute()+1, 0, 0, t.Location())
}
