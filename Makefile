test:
	@go test

build:
	@go build -o bin/lgtm cmd/lgtm/main.go

run:
	@go run cmd/lgtm/main.go

install:
	@go get ./...
