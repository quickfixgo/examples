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

# Commands for docker images.
# ----------------------------
build-linux:
	GOOS=linux GOARCH=amd64 go build -v -o ./bin/executor ./cmd/executor
	GOOS=linux GOARCH=amd64 go build -v -o ./bin/ordermatch ./cmd/ordermatch
	GOOS=linux GOARCH=amd64 go build -v -o ./bin/tradeclient ./cmd/tradeclient

build-docker: clean build-linux
	docker build -t quickfixgo/executor:latest -f ./cmd/executor/Dockerfile .
	docker build -t quickfixgo/ordermatch:latest -f ./cmd/ordermatch/Dockerfile .
	docker build -t quickfixgo/tradeclient:latest -f ./cmd/tradeclient/Dockerfile .