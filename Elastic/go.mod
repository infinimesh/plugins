module github.com/infinimesh/plugins/Elastic

go 1.15

replace (
	github.com/infinimesh/plugins/pkg => ../pkg
	github.com/infinimesh/plugins/redisstream => ../redisstream
)

require (
	github.com/infinimesh/plugins/redisstream v0.0.0
	github.com/elastic/go-elasticsearch/v7 v7.0.0
	github.com/gomodule/redigo v1.8.3
)
