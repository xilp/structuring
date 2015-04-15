.PHONY: all build test

all: build test

build:
	go clean ./...
	go install -v ./...

test:
	go test -v ./...
