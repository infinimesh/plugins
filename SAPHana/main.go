package main

import (
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/InfiniteDevices/plugins/redisstream/consumer"
	"github.com/SAP/go-hdb/driver"
	"github.com/gomodule/redigo/redis"
)

const (
	envSapHanaInitDB   = "SAPHANA_INITDB"
	envSapHanaUsername = "SAPHANA_USERNAME"
	envSapHanaPassword = "SAPHANA_PASSWORD"
	envSapHanaHost     = "SAPHANA_HOST"
	envSapHanaPort     = "SAPHANA_PORT"
	envRedisAddr       = "REDIS_ADDR"
)

func main() {
	db := newSapHanaDB(
		os.Getenv(envSapHanaHost),
		os.Getenv(envSapHanaPort),
		os.Getenv(envSapHanaUsername),
		os.Getenv(envSapHanaPassword),
	)
	shouldInitDB, _ := strconv.ParseBool(os.Getenv(envSapHanaInitDB))
	if shouldInitDB {
		if err := initDB(db); err != nil {
			log.Printf("failed to initialise db: err=%v\n", err)
		}
	}
	redisPool := &redis.Pool{Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", os.Getenv(envRedisAddr))
	}}
	c := consumer.New(redisPool)
	for event := range c.Consume() {
		if event == nil {
			log.Println("received nil event")
			continue
		}
		err := insertDeviceEvent(db, event)
		if err != nil {
			log.Printf("failed to insert event to SAP HANA: err=%v", err)
		} else {
			log.Println("inserted event to SAP HANA:", event)
		}
	}
}

// for some reason, connecting to the db via connection string as mentioned in
// the go-hdb driver documentation doesn't work, thus we opt for this
// alternative connection method
//
// also see https://stackoverflow.com/questions/58698188
func newSapHanaDB(host, port, username, password string) *sql.DB {
	c := driver.NewBasicAuthConnector(host+":"+port, username, password)
	tlsConfig := tls.Config{
		InsecureSkipVerify: false,
		ServerName:         host,
	}
	c.SetTLSConfig(&tlsConfig)
	return sql.OpenDB(c)
}

// initDB initialises the required relations in SAP HANA for storage
//
// ideally we would want to use the SAP HANA JSON Docstore, but the cloud
// offering doesn't support it as of time of writing. Here we opt for an
// alternative where we store retrieved data as a JSON string in a table
//
// also see https://answers.sap.com/questions/13207475/
func initDB(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE devices (uid varchar(1000), timestamp timestamp, version varchar(1000), data varchar(5000), PRIMARY KEY (uid, timestamp));`)
	return err
}

func insertDeviceEvent(db *sql.DB, event *consumer.DeviceEvent) error {
	t, err := time.Parse(time.RFC3339Nano, event.State.Timestamp)
	if err != nil {
		return err
	}
	jsonData, _ := json.Marshal(event.State.Data)
	_, err = db.Exec("INSERT INTO devices (uid, timestamp, version, data) VALUES (?, ?, ?, ?)", event.Object.UID, t, event.State.Version, string(jsonData))
	return err
}
