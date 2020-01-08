package cli

import (
	"github.com/go-redis/redis/v7"
	"github.com/ryotarai/redis-queue-worker/pkg/worker"
	"github.com/urfave/cli/v2"
)

func startAction(c *cli.Context) error {
	if err := handleGlobalFlags(c); err != nil {
		return err
	}

	redisAddr := c.String("redis-addr")
	redisPassword := c.String("redis-password")
	redisDB := c.Int("redis-db")
	redisQueueKey := c.String("redis-queue-key")
	args := c.Args().Slice()

	logger.Info().Str("redisAddr", redisAddr).Int("redisDB", redisDB).
		Str("redisQueueKey", redisQueueKey).Msg("Starting redis-queue-worker...")

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	worker, err := worker.New(logger, client, redisQueueKey, args)
	if err != nil {
		return err
	}

	if err := worker.Start(); err != nil {
		return err
	}

	return nil
}
