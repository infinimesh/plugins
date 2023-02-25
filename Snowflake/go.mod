module github.com/infinimesh/plugins/Snowflake

go 1.15

replace (
	github.com/infinimesh/plugins/pkg => ../pkg
	github.com/infinimesh/plugins/redisstream => ../redisstream
)

require (
	github.com/gomodule/redigo v1.8.3
	github.com/infinimesh/plugins/redisstream v0.0.0
	github.com/snowflakedb/gosnowflake v1.3.12
	golang.org/x/crypto v0.1.0 // indirect
)
