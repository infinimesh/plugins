package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/InfiniteDevices/plugins/grafana/api"
	"github.com/InfiniteDevices/plugins/grafana/proc"
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

	pool := &redis.Pool{Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", redisAddr)
	}}
	client := redistimeseries.NewClientFromPool(pool, "plugin")

	proc.New(
		api.NewHandler(
			api.NewTokenHandler(username, password, apiUrl, time.Duration(tokenRefreshIntervalSecs)*time.Second),
			apiUrl,
		),
		client,
	).Run(time.Duration(deviceRefreshIntervalSecs) * time.Second)
}
