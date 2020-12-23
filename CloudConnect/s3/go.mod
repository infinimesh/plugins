module github.com/infinimesh/plugins/CloudConnect/s3

go 1.15

require (
	github.com/aws/aws-sdk-go v1.36.12
	github.com/infinimesh/plugins/CloudConnect/csvprocessor v0.0.0
	golang.org/x/net v0.0.0-20201216054612-986b41b23924 // indirect
	golang.org/x/text v0.3.4 // indirect
)

replace github.com/infinimesh/plugins/CloudConnect/csvprocessor => ../csvprocessor
