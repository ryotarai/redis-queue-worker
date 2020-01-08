FROM golang:1.13.5 as builder

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . /workspace

RUN make

FROM ubuntu:18.04
#RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /
COPY --from=builder /workspace/bin/redis-queue-worker /usr/local/bin/redis-queue-worker
ENTRYPOINT ["/usr/local/bin/redis-queue-worker"]
