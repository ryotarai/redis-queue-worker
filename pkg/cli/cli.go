package cli

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
	"os"
	"time"
)

var logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Timestamp().Logger()
var globalFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "log-level",
		EnvVars: []string{"LOG_LEVEL"},
		Value:   "info",
	},
}

func Run(args []string) error {
	return BuildApp().Run(args)
}

func BuildApp() *cli.App {
	startFlags := []cli.Flag{
		&cli.StringFlag{
			Name:     "redis-addr",
			EnvVars:  []string{"REDIS_ADDR"},
			Required: true,
		},
		&cli.StringFlag{
			Name:     "redis-queue-key",
			EnvVars:  []string{"REDIS_QUEUE_KEY"},
			Required: true,
		},
		&cli.StringFlag{
			Name:    "redis-password",
			EnvVars: []string{"REDIS_PASSWORD"},
		},
		&cli.IntFlag{
			Name:    "redis-db",
			Value:   0,
			EnvVars: []string{"REDIS_DB"},
		},
	}
	startFlags = append(startFlags, globalFlags...)

	return &cli.App{
		Commands: []*cli.Command{
			{
				Name:   "version",
				Action: versionAction,
			},
			{
				Name:   "start",
				Action: startAction,
				Flags: startFlags,
			},
		},
	}
}

func handleGlobalFlags(c *cli.Context) error {
	logLevel := c.String("log-level")
	switch logLevel {
	case "trace":
		logger = logger.Level(zerolog.TraceLevel)
	case "debug":
		logger = logger.Level(zerolog.DebugLevel)
	case "info":
		logger = logger.Level(zerolog.InfoLevel)
	case "warn":
		logger = logger.Level(zerolog.WarnLevel)
	case "error":
		logger = logger.Level(zerolog.ErrorLevel)
	default:
		return fmt.Errorf(`"%s" is invalid log level`, logLevel)
	}
	return nil
}
