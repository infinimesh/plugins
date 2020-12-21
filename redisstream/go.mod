module github.com/infinimesh/plugins/redisstream

go 1.15

require (
	github.com/infinimesh/plugins/pkg v0.0.0
	github.com/gomodule/redigo v1.8.3
	github.com/stretchr/testify v1.6.1 // indirect
)

replace github.com/infinimesh/plugins/pkg => ../pkg
