package worker

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"os/exec"
	"time"
)

var backoffMax = time.Second * 256

type Worker struct {
	redisClient *redis.Client
	redisKey    string
	command     []string
	logger      zerolog.Logger
}

func New(logger zerolog.Logger, redisClient *redis.Client, redisKey string, command []string) (*Worker, error) {
	return &Worker{
		logger:      logger,
		redisClient: redisClient,
		redisKey:    redisKey,
		command:     command,
	}, nil
}

func (w *Worker) Start() error {
	myUUID := uuid.New()
	processingQueueKey := fmt.Sprintf("%s:worker:%s:processing", w.redisKey, myUUID.String())

	backoff := time.Second
	for {
		job, err := w.redisClient.RPopLPush(w.redisKey, processingQueueKey).Result()
		if err == redis.Nil {
			w.logger.Info().Msg("Exiting as no job is in queue")
			break
		}

		logger := w.logger.With().Str("job", job).Logger()

		for {
			if err := w.runCommand(logger, job); err == nil {
				break
			}
			logger.Info().Dur("backoff", backoff).Msg("Retrying...")
			time.Sleep(backoff)
			backoff = backoff * 2
			if backoff > backoffMax {
				backoff = backoffMax
			}
		}

		logger.Info().Msg("Succeeded to run a job")

		_, err = w.redisClient.LRem(processingQueueKey, 1, job).Result()
		if err != nil {
			logger.Warn().Str("queueKey", processingQueueKey).
				Msg("Failed to delete a job from processing queue")
		}
	}

	return nil
}

func (w *Worker) runCommand(logger zerolog.Logger, job string) error {
	logger.Info().Strs("command", w.command).Msg("Executing a command")

	cmd := exec.Command(w.command[0], w.command[1:]...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	fmt.Fprintln(stdin, job)
	stdin.Close()

	output, err := cmd.CombinedOutput()
	if err != nil {
		if exerr, ok := err.(*exec.ExitError); ok {
			logger.Warn().Int("exitCode", exerr.ExitCode()).Bytes("output", output).Msg("Failed to execute a commmand")
		}
		return err
	} else {
		logger.Debug().Bytes("output", output).Msg("Succeeded to execute a command")
	}

	return nil
}
