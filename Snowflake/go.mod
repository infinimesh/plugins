module github.com/InfiniteDevices/plugins/Snowflake

go 1.15

replace (
	github.com/InfiniteDevices/plugins/pkg => ../pkg
	github.com/InfiniteDevices/plugins/redisstream => ../redisstream
)

require (
	github.com/InfiniteDevices/plugins/redisstream v0.0.0
	github.com/gomodule/redigo v1.8.3
	github.com/snowflakedb/gosnowflake v1.3.12
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9 // indirect
)
