SHELL = /bin/bash

GOOS    ?= $(shell go env GOOS)
GOARCH  ?= $(shell go env GOARCH)
GOPATH  ?= $(shell go env GOPATH)
VERSION := $(shell git describe --tags --abbrev=5 --always)

export GOPATH

.PHONY = environment build dependency tests
.RECIPEPREFIX = >

environment:
> @echo -e "system environment:\n"
> @env | sort -u;
> @echo -e "\ngo environment:\n"
> @go env

dependency:
> go get -v -u ./...
> go mod tidy

lint:
> go vet -x ./...
> gofmt -d -s .

build: dependency
> @echo -e "\nUsing $(GOPATH) as GOPATH";
> go build -v -o bin/redis-scalor \
>  -ldflags='-s -w -X main.version=$(VERSION)' ./cmd;

test:
> @echo -e "\nTesting using $(GOPATH) as GOPATH";
> go test -v ./...

run-redis-cluster:
> @mage;
> @mage listRedis;

stop-redis-cluster:
> @mage stopRedis;
