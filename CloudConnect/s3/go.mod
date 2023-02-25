module github.com/infinimesh/plugins/CloudConnect/s3

go 1.15

require (
	github.com/aws/aws-sdk-go v1.36.12
	github.com/infinimesh/plugins/CloudConnect/csvprocessor v0.0.0
	golang.org/x/net v0.7.0 // indirect
)

replace github.com/infinimesh/plugins/CloudConnect/csvprocessor => ../csvprocessor
