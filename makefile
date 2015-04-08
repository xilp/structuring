.PHONY: all test travis

travis: all test

all:
	go install -v ./...

test:
	go test -v ./...
