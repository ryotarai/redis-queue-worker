package main

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/ryotarai/redis-queue-worker/pkg/cli"
	"log"
	"os"
	"strconv"
)

func main() {
	if err := cli.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

//func

func run() error {
	redisDBStr := os.Getenv("REDIS_DB")
	if redisDBStr == "" {
		redisDBStr = "0"
	}
	redisDB, err := strconv.Atoi(redisDBStr)
	if err != nil {
		return err
	}

	redisKey := os.Getenv("REDIS_KEY")
	if redisKey == "" {
		return fmt.Errorf("REDIS_KEY environment variable is required")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       redisDB,
	})

	log.Printf("Connecting %s/%d (key: %s)")

	item, err := client.LPop(redisKey).Result()
	if err == redis.Nil {
		log.Println(item)
	}

	return nil
}
