# redis-queue-worker

## Usage

Put jobs into a queue:

```
$ redis-cli lpush examplequeue foo
(integer) 1
$ redis-cli lpush examplequeue bar
(integer) 2
$ redis-cli lpush examplequeue baz
(integer) 3
```

The job string will be passed to a worker command via stdin.

Start a worker:

```
$ mkdir /tmp/redis-queue-worker-example
$ redis-queue-worker start --redis-addr localhost:6379 --redis-queue-key examplequeue -- sh -c 'input=$(cat); touch /tmp/redis-queue-worker-example/$input'
2020-01-08T09:57:59+09:00 INF Starting redis-queue-worker... redisAddr=localhost:6379 redisDB=0 redisQueueKey=examplequeue
2020-01-08T09:57:59+09:00 INF Executing a command command=["sh","-c","input=$(cat); touch /tmp/$input"] job=foo
2020-01-08T09:57:59+09:00 INF Succeeded to run a job job=foo
2020-01-08T09:57:59+09:00 INF Executing a command command=["sh","-c","input=$(cat); touch /tmp/$input"] job=bar
2020-01-08T09:57:59+09:00 INF Succeeded to run a job job=bar
2020-01-08T09:57:59+09:00 INF Executing a command command=["sh","-c","input=$(cat); touch /tmp/$input"] job=baz
2020-01-08T09:57:59+09:00 INF Succeeded to run a job job=baz
2020-01-08T09:57:59+09:00 INF Exiting as no job is in queue
```

```
$ ls /tmp/redis-queue-worker-example
foo bar baz
```

## Build

```
make build
```
