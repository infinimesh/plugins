package csvprocessor

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

const envReadDir = "READ_DIR"

func WalkLoop(f filepath.WalkFunc) {
	readDir := os.Getenv(envReadDir)
	for range time.Tick(5 * time.Second) {
		log.Printf("walking %s for new files...", readDir)
		err := filepath.Walk(readDir, f)
		if err != nil {
			log.Printf("failed to walk directory: %v\n", err)
		}
	}
}

func WalkFunc(handleFile func(*os.File) error) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		// only process files that have not been modified recently
		if time.Now().Add(-1 * time.Minute).After(info.ModTime()) {
			log.Println("processing file:", path)
			f, err := os.Open(path)
			if err != nil {
				log.Printf("failed to open %s: %v\n", path, err)
				return err
			}
			defer f.Close()
			err = handleFile(f)
			if err != nil {
				log.Printf("failed to handle file: %v\n", err)
				return err
			} else {
				log.Printf("successfully handled file: %s\n", path)
			}
			err = os.Remove(path)
			if err != nil {
				log.Printf("failed to remove object from path: %v\n", err)
				return err
			} else {
				log.Printf("successfully deleted %s\n", path)
			}
		}
		return nil
	}
}
