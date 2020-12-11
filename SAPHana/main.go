package main

import (
	"log"
	"os"

	"github.com/InfiniteDevices/plugins/redisstream/consumer"
	"github.com/gomodule/redigo/redis"
)

const (
	envRedisAddr = "REDIS_ADDR"
)

func main() {
	redisPool := &redis.Pool{Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", os.Getenv(envRedisAddr))
	}}
	c := consumer.New(redisPool)
	for event := range c.Consume() {
		log.Println("received event:", event)
	}
}
