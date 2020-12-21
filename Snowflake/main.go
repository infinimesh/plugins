package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/infinimesh/plugins/redisstream/consumer"
	"github.com/gomodule/redigo/redis"
	sf "github.com/snowflakedb/gosnowflake"
)

const (
	envSnowflakeAccount  = "SNOWFLAKE_ACCOUNT"
	envSnowflakeUser     = "SNOWFLAKE_USER"
	envSnowflakePassword = "SNOWFLAKE_PASSWORD"
	envSnowflakeInitDB   = "SNOWFLAKE_INITDB"
	envRedisAddr         = "REDIS_ADDR"
)

func main() {
	account := os.Getenv(envSnowflakeAccount)
	user := os.Getenv(envSnowflakeUser)
	password := os.Getenv(envSnowflakePassword)
	dsn, err := sf.DSN(&sf.Config{
		Account:  account,
		User:     user,
		Password: password,
	})
	if err != nil {
		log.Fatalf("failed to create DSN, err: %v", err)
	}

	db, err := sql.Open("snowflake", dsn)
	if err != nil {
		log.Fatalf("failed to connect. %v, err: %v", dsn, err)
	}
	defer db.Close()

	shouldInitDB, _ := strconv.ParseBool(os.Getenv(envSnowflakeInitDB))
	if shouldInitDB {
		initDB(db)
	}

	_, err = db.Exec(`USE DATABASE infinimesh;`)
	if err != nil {
		log.Fatalf("failed to use infinimesh database, err: %v", err)
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
			log.Printf("failed to insert event to Snowflake: err=%v", err)
		} else {
			log.Println("inserted event to Snowflake:", event)
		}
	}
}

func initDB(db *sql.DB) {
	for _, stmt := range []string{
		`CREATE DATABASE IF NOT EXISTS infinimesh;`,
		`CREATE TABLE IF NOT EXISTS device_states (
			uid VARCHAR(1000),
			timestamp TIMESTAMP,
			version VARCHAR(1000),
			data VARIANT,
			PRIMARY KEY (uid, timestamp)
		);`,
	} {
		_, err := db.Exec(stmt)
		if err != nil {
			log.Printf("failed to execute db init statement :%s with err: %v\n", stmt, err)
		}
	}
}

func insertDeviceEvent(db *sql.DB, event *consumer.DeviceEvent) error {
	t, err := time.Parse(time.RFC3339Nano, event.State.Timestamp)
	if err != nil {
		return err
	}
	jsonStr, _ := json.Marshal(event.State.Data)
	// passing variant data as an argument in the query seems to be
	// unsupported, thus here we opt for a workaround where we pass the json
	// string into SELECT and PARSE_JSON() instead. This is assumed to be safe
	// from sql injection
	//
	// also see:
	// - https://community.snowflake.com/s/question/0D50Z00008DCk08SAD/
	// - https://github.com/snowflakedb/snowflake-connector-nodejs/issues/59
	_, err = db.Exec("INSERT INTO device_states (uid, timestamp, version, data) SELECT ?, ?, ?, PARSE_JSON('"+string(jsonStr)+"');", event.Object.UID, t, event.State.Version)
	return err
}
