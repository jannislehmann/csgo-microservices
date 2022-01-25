GOPATH:=$(shell go env GOPATH)

.PHONY: init
init:
	@go get -u google.golang.org/protobuf/proto
	@go get -u github.com/vektra/mockery/v2/.../
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
