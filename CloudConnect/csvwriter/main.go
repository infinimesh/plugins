package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/InfiniteDevices/plugins/pkg/api"
	"github.com/InfiniteDevices/plugins/pkg/wrappers"
)

const (
	envApiUrl                    = "API_URL"
	envUsername                  = "USERNAME"
	envPassword                  = "PASSWORD"
	envTokenRefreshIntervalSecs  = "TOKEN_REFRESH_INTERVAL_SECS"
	envObjectRefreshIntervalSecs = "OBJECT_REFRESH_INTERVAL_SECS"
	envWriteDir                  = "WRITE_DIR"
)

func main() {
	username := os.Getenv(envUsername)
	password := os.Getenv(envPassword)
	apiUrl := os.Getenv(envApiUrl)
	tokenRefreshIntervalSecsStr := os.Getenv(envTokenRefreshIntervalSecs)
	tokenRefreshIntervalSecs, err := strconv.Atoi(tokenRefreshIntervalSecsStr)
	if err != nil {
		log.Fatalf("invalid %s value: %s", envTokenRefreshIntervalSecs, tokenRefreshIntervalSecsStr)
	}
	objectRefreshIntervalSecsStr := os.Getenv(envObjectRefreshIntervalSecs)
	objectRefreshIntervalSecs, err := strconv.Atoi(objectRefreshIntervalSecsStr)
	if err != nil {
		log.Fatalf("invalid %s value: %s", envObjectRefreshIntervalSecs, objectRefreshIntervalSecsStr)
	}
	apiHandler := api.NewHandler(
		api.NewTokenHandler(username, password, apiUrl, time.Duration(tokenRefreshIntervalSecs)*time.Second),
		apiUrl,
	)
	writeDir := os.Getenv(envWriteDir)
	err = os.MkdirAll(writeDir, os.ModePerm)
	if err != nil {
		log.Fatalf("failed to mkdir: %v", err)
	}
	wrappers.NewObjectManager(
		apiHandler,
		(&objectWorkerFactory{
			api: apiHandler,
			writer: &writer{
				writeDir: writeDir,
			},
		}).NewObjectWorker,
		time.Duration(objectRefreshIntervalSecs)*time.Second,
	).Start()
}

type objectWorkerFactory struct {
	api    api.Handler
	writer *writer
}

func (f *objectWorkerFactory) NewObjectWorker(obj api.Object) wrappers.Process {
	return &objectWorker{
		obj:    obj,
		api:    f.api,
		writer: f.writer,
		done:   make(chan struct{}),
	}
}

type objectWorker struct {
	obj    api.Object
	api    api.Handler
	writer *writer
	done   chan struct{}
}

func (w *objectWorker) Start() {
	ch, err := w.api.GetDevicesStateStream(w.obj.UID)
	if err != nil {
		log.Printf("error on get devices state stream: %s\n", err)
		return
	}

	for {
		select {
		case <-w.done:
			return
		case state := <-ch:
			if state == nil {
				log.Printf("received nil state for object %v", w.obj)
				continue
			}
			j, _ := json.Marshal(state.Result.ReportedState.Data)
			record := []string{
				w.obj.UID,
				w.obj.Name,
				w.obj.Kind,
				state.Result.ReportedState.Timestamp,
				state.Result.ReportedState.Version,
				string(j),
			}
			err := w.writer.Write(record)
			if err != nil {
				log.Printf("failed to write record to csv: %v\n", err)
			} else {
				log.Printf("successfully wrote record to csv: %v\n", record)
			}
		}
	}
}

func (w *objectWorker) Stop() {
	close(w.done)
}
