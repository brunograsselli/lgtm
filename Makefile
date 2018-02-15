version := $(shell cat VERSION)

test:
	@go test ./lgtm

build:
	@go build -o ./bin/lgtm ./main.go

build_all:
	env GOOS=darwin GOARCH=amd64 go build -o ./bin/lgtm-darwin-amd64-$(version) ./cmd/lgtm/main.go
	env GOOS=linux GOARCH=amd64 go build -o ./bin/lgtm-linux-amd64-$(version) ./cmd/lgtm/main.go

run:
	@go run cmd/lgtm/main.go

install:
	@go get ./...
