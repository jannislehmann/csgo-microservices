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
	@go build -ldflags "-X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=warn" -o gamecoordinatorclient *.go

.PHONY: test
test:
	@GOLANG_PROTOBUF_REGISTRATION_CONFLICT=warn go test -v ./... -cover

.PHONY: docker
docker:
	@docker build -t gamecoordinatorclient:latest .

.PHONY: mock
mock:
	@mockery --all --dir=internal/gamecoordinator
