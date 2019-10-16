SHELL := /bin/bash

test: lint vet build

lint:
	golint ./...

vet:
	go vet ./...

build: clean
	go build -v -o ./bin/executor ./cmd/executor
	go build -v -o ./bin/ordermatch ./cmd/ordermatch
	go build -v -o ./bin/tradeclient ./cmd/tradeclient

clean:
	rm -rf ./bin
	rm -rf ./tmp