GOPATH:=$(shell go env GOPATH)

.PHONY: init
init:
	@go get -u google.golang.org/protobuf/proto
	@go get -u github.com/vektra/mockery/v2/.../
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1

.PHONY: proto
proto:
	@protoc --go_out=. --go-grpc_out=require_unimplemented_servers=false:. proto/*.proto

.PHONY: update
update:
	@go get -u

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: build
build:
	@go build -o valveapiclient *.go

.PHONY: test
test:
	@go test -v ./... -cover

.PHONY: docker
docker:
	@docker build -t valveapiclient:latest .

.PHONY: mock
mock:
	@mockery --all --dir=internal/domain
	@mockery --all --dir=pkg/
