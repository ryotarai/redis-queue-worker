package cli

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

func versionAction(c *cli.Context) error {
	fmt.Println("redis-queue-worker v0.0.1")
	return nil
}
