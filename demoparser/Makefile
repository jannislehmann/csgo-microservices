GOPATH:=$(shell go env GOPATH)

.PHONY: init
init:
	@go get -u github.com/vektra/mockery/v2/.../
	@export PATH="$PATH:$(go env GOPATH)/bin"

.PHONY: update
update:
	@go get -u

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: build
build:
	@go build -o demoparser *.go

.PHONY: test
test:
	@go test -v ./... -cover

.PHONY: docker
docker:
	@docker build -t demoparser:latest .

.PHONY: mock
mock:
	@mockery --all --dir=internal/domain
	@mockery --all --dir=pkg/
