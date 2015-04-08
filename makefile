.PHONY: all test travis

travis: all test

all:
	go clean -v ./...
	go install -v ./...

test:
	go test -v ./...
