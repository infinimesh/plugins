package consumer

import (
	"encoding/json"
	"log"
	"os"
	"runtime/debug"
	"time"

	"github.com/InfiniteDevices/plugins/pkg/api"
	"github.com/gomodule/redigo/redis"
)

type RedisStreamEvent struct {
	ID          string
	DeviceEvent *DeviceEvent
}

type DeviceEvent struct {
	Object api.Object
	State  api.DeviceState
}

type Consumer interface {
	Consume() <-chan *DeviceEvent
}

type consumerImpl struct {
	redis *redis.Pool
	name  string
}

func New(pool *redis.Pool) Consumer {
	hn, _ := os.Hostname()
	return &consumerImpl{
		redis: pool,
		name:  hn,
	}
}

func (i *consumerImpl) Consume() <-chan *DeviceEvent {
	ret := make(chan *DeviceEvent)
	go func() {
		for {
			shouldCooldown := i.loop(ret)
			if shouldCooldown {
				time.Sleep(time.Second)
			}
		}
	}()
	return ret
}

func (i *consumerImpl) loop(ch chan<- *DeviceEvent) bool {
	conn := i.redis.Get()
	defer conn.Close()
	reply, err := conn.Do("XREADGROUP", "GROUP", "group", i.name, "STREAMS", "objects", ">")
	if err != nil {
		log.Printf("error on XREADGROUP: %v\n", err)
		// XREADGROUP errors if the consumer group does not yet exist in the
		// stream. Most of the time however, the consumer group should have
		// already been created, thus this is usually not a problem. Placing
		// the group creation here feels more self-healing and slightly more
		// efficient than creating the group prior to consuming from the stream
		i.createGroupIfNotExists()
		return true
	}
	if reply == nil {
		return true
	}
	events := parseReply(reply)
	for _, e := range events {
		if _, err := conn.Do("XACK", "objects", "group", e.ID); err != nil {
			log.Printf("error on XACK: %v\n", err)
		}
		ch <- e.DeviceEvent
	}
	return false
}

func (i *consumerImpl) createGroupIfNotExists() {
	conn := i.redis.Get()
	defer conn.Close()
	_, err := conn.Do("XGROUP", "CREATE", "objects", "group", "$", "MKSTREAM")
	if err != nil {
		log.Printf("error on XGROUP CREATE: %v\n", err)
	}
}

func parseReply(reply interface{}) []RedisStreamEvent {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("unexpected error when parsing reply: %v\n", err)
			debug.PrintStack()
		}
	}()

	ret := []RedisStreamEvent{}
	streamAndEventss := reply.([]interface{})
	for _, se := range streamAndEventss {
		streamAndEvents := se.([]interface{})
		events := streamAndEvents[1].([]interface{})
		for _, e := range events {
			event := e.([]interface{})
			eventID := string(event[0].([]byte))
			eventData := event[1].([]interface{})
			ret = append(ret, RedisStreamEvent{
				ID:          eventID,
				DeviceEvent: parseEventData(eventData),
			})
		}
	}
	return ret
}

func parseEventData(eventData []interface{}) *DeviceEvent {
	if eventData == nil {
		return nil
	}
	ret := &DeviceEvent{}
	for i := 0; i < len(eventData); i += 2 {
		k := string(eventData[i].([]byte))
		v := eventData[i+1].([]byte)
		switch k {
		case "object":
			obj := api.Object{}
			err := json.Unmarshal(v, &obj)
			if err != nil {
				log.Printf("unexpected error when unmarshaling object: %v\n", err)
				continue
			}
			ret.Object = obj
		case "state":
			state := api.DeviceState{}
			err := json.Unmarshal(v, &state)
			if err != nil {
				log.Printf("unexpected error when unmarshaling object: %v\n", err)
			}
			ret.State = state
		default:
			log.Printf("unexpected key when parsing event data: %v\n", k)
		}
	}
	return ret
}
