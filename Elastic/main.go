package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/InfiniteDevices/plugins/redisstream/consumer"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gomodule/redigo/redis"
)

const (
	envRedisAddr    = "REDIS_ADDR"
	envEsAddr       = "ES_ADDR"
	indexInfinimesh = "infinimesh"
)

func main() {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{os.Getenv(envEsAddr)},
	})
	if err != nil {
		log.Fatalf("failed to initialise es client: %v\n", err)
	}
	res, err := es.Ping()
	if err != nil {
		log.Fatalf("failed to ping es: %v\n", err)
	}
	log.Printf("successfully connected to es with result: %v\n", res)
	redisPool := &redis.Pool{Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", os.Getenv(envRedisAddr))
	}}
	c := consumer.New(redisPool)
	for event := range c.Consume() {
		if event == nil {
			log.Println("received nil event")
			continue
		}
		t, err := time.Parse(time.RFC3339Nano, event.State.Timestamp)
		if err != nil {
			log.Printf("failed to parse device timestamp: %v\n", err)
			continue
		}
		id := event.Object.UID + "-" + strconv.FormatInt(t.Unix(), 10)
		b, _ := json.Marshal(event.State)
		res, err := es.Index(indexInfinimesh, bytes.NewReader(b), es.Index.WithDocumentID(id))
		if err != nil {
			log.Printf("failed to index document: %v\n", err)
			continue
		}
		log.Printf("successfully indexed document: %v\n", res)
	}
}
