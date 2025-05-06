pwd = $(shell basename $(shell pwd))

.PHONY: build
build:
	go build -o ./bin/weather ./cmd/weather