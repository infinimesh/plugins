module github.com/infinimesh/plugins/SAPHana

go 1.15

replace (
	github.com/infinimesh/plugins/pkg => ../pkg
	github.com/infinimesh/plugins/redisstream => ../redisstream
)

require (
	github.com/infinimesh/plugins/redisstream v0.0.0
	github.com/SAP/go-hdb v0.102.6
	github.com/gomodule/redigo v1.8.3
)
