module github.com/InfiniteDevices/plugins/timeseries

go 1.15

require (
	github.com/InfiniteDevices/plugins/pkg v0.0.0
	github.com/RedisTimeSeries/redistimeseries-go v1.4.3
	github.com/gomodule/redigo v1.8.3
)

replace github.com/InfiniteDevices/plugins/pkg => ../pkg
