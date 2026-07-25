package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ajsqr/wombat/agent"
)

type cli struct {
	agent   *agent.Agent
	logger  *slog.Logger
	version string
	commit  string
}

func (c *cli) run() {
	if len(os.Args) < 2 {
		os.Exit(1)
	}

	baseCmd := os.Args[1]
	switch baseCmd {
	case "version":
		fmt.Printf("wombat-agent %s build %s\n", c.version, c.commit)
	case "run":
		c.agent.Run()
	default:
		c.logger.Error("unexpected cli command", slog.String("cmd", baseCmd))
		os.Exit(1)
	}
}
