GOPATH:=$(shell go env GOPATH)

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
	@go build -o gamecoordinatorclient *.go

.PHONY: test
test:
	@go test -v ./... -cover

.PHONY: docker
docker:
	@docker build -t gamecoordinatorclient:latest .

.PHONY: mock
mock:
	@mockery --all --dir=internal/gamecoordinator
