module github.com/InfiniteDevices/plugins/Elastic

go 1.15

replace (
	github.com/InfiniteDevices/plugins/pkg => ../pkg
	github.com/InfiniteDevices/plugins/redisstream => ../redisstream
)

require (
	github.com/InfiniteDevices/plugins/redisstream v0.0.0
	github.com/elastic/go-elasticsearch/v7 v7.0.0
	github.com/gomodule/redigo v1.8.3
)
