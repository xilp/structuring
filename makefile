.PHONY: all test travis

travis: all test

all:
	go clean ./...
	go install -v ./...

test:
	go test -v ./...
