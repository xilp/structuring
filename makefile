.PHONY: all build test

all: build test

build:
	go clean ./...
	go install -v ./...

client:
	./test.sh
	echo 

test:
	go test -v ./...
