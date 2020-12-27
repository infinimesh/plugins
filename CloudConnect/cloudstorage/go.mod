module github.com/infinimesh/plugins/CloudConnect/cloudstorage

go 1.15

require (
	cloud.google.com/go/storage v1.12.0
	github.com/infinimesh/plugins/CloudConnect/csvprocessor v0.0.0
)

replace github.com/infinimesh/plugins/CloudConnect/csvprocessor => ../csvprocessor
