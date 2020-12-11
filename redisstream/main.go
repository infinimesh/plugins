package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/InfiniteDevices/plugins/pkg/api"
	"github.com/InfiniteDevices/plugins/pkg/wrappers"
	"github.com/gomodule/redigo/redis"
)

const (
	envApiUrl                    = "API_URL"
	envUsername                  = "USERNAME"
	envPassword                  = "PASSWORD"
	envTokenRefreshIntervalSecs  = "TOKEN_REFRESH_INTERVAL_SECS"
	envObjectRefreshIntervalSecs = "OBJECT_REFRESH_INTERVAL_SECS"
	envRedisAddr                 = "REDIS_ADDR"
)

func main() {
	redisPool := &redis.Pool{Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", os.Getenv(envRedisAddr))
	}}
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

	wrappers.NewObjectManager(
		apiHandler,
		(&objectWorkerFactory{
			api:   apiHandler,
			redis: redisPool,
		}).NewObjectWorker,
		time.Duration(objectRefreshIntervalSecs)*time.Second,
	).Start()
}

type objectWorkerFactory struct {
	api   api.Handler
	redis *redis.Pool
}

func (f *objectWorkerFactory) NewObjectWorker(obj api.Object) wrappers.Process {
	return &objectWorker{
		obj:   obj,
		api:   f.api,
		redis: f.redis,
		done:  make(chan struct{}),
	}
}

type objectWorker struct {
	obj   api.Object
	api   api.Handler
	redis *redis.Pool
	done  chan struct{}
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
			objJson, _ := json.Marshal(w.obj)
			stateJson, _ := json.Marshal(state.Result.ReportedState)
			conn := w.redis.Get()
			reply, err := conn.Do("XADD", "objects", "*", "object", string(objJson), "state", string(stateJson))
			conn.Close()
			if err != nil {
				log.Printf("failed to stream object: object=%v err=%v\n", w.obj, err)
			} else {
				log.Printf("successfully streamed object: object=%v reply=%v\n", w.obj, string(reply.([]byte)))
			}
		}
	}
}

func (w *objectWorker) Stop() {
	close(w.done)
}
