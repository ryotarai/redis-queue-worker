export GO111MODULE=on

.PHONY: build
build:
	go build -o bin/redis-queue-worker .
