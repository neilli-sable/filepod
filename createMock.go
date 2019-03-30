package main

//go:generate mockgen -source vendor/github.com/aws/aws-sdk-go/service/s3/s3iface/interface.go -destination mock/s3_mock.go -package mock

//go:generate mockgen -source vendor/github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface/interface.go -destination mock/s3manager_mock.go -package mock
