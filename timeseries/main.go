package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/InfiniteDevices/plugins/pkg/api"
	"github.com/InfiniteDevices/plugins/pkg/wrappers"
	redistimeseries "github.com/RedisTimeSeries/redistimeseries-go"
	"github.com/gomodule/redigo/redis"
)

const (
	envUsername                  = "USERNAME"
	envPassword                  = "PASSWORD"
	envApiUrl                    = "API_URL"
	envTokenRefreshIntervalSecs  = "TOKEN_REFRESH_INTERVAL_SECS"
	envDeviceRefreshIntervalSecs = "DEVICE_REFRESH_INTERVAL_SECS"
	envRedisAddr                 = "REDIS_ADDR"
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
	deviceRefreshIntervalSecsStr := os.Getenv(envDeviceRefreshIntervalSecs)
	deviceRefreshIntervalSecs, err := strconv.Atoi(deviceRefreshIntervalSecsStr)
	if err != nil {
		log.Fatalf("invalid %s value: %s", envDeviceRefreshIntervalSecs, deviceRefreshIntervalSecsStr)
	}
	redisAddr := os.Getenv(envRedisAddr)
	redisPool := &redis.Pool{Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", redisAddr)
	}}
	redisClient := redistimeseries.NewClientFromPool(redisPool, "plugin")
	apiHandler := api.NewHandler(
		api.NewTokenHandler(username, password, apiUrl, time.Duration(tokenRefreshIntervalSecs)*time.Second),
		apiUrl,
	)
	wrappers.NewObjectManager(
		apiHandler,
		(&objectWorkerFactory{
			api:         apiHandler,
			redisClient: redisClient,
		}).NewObjectWorker,
		time.Duration(deviceRefreshIntervalSecs)*time.Second,
	).Start()
}
