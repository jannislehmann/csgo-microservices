GOPATH:=$(shell go env GOPATH)

.PHONY: init
init:
	@go get -u github.com/vektra/mockery/v2/.../

.PHONY: update
update:
	@go get -u

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: test
test:
	@go test -v ./... -cover

.PHONY: mock
mock:
	@mockery --all --dir=pkg/queue
