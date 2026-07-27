package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ajsqr/wombat/agent"
)

type cli struct {
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
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		wa, err := agent.NewAgent(logger)
		if err != nil {
			logger.Error("error initializing agent", slog.Any("error", err))
			os.Exit(1)
		}
		wa.Run()
	default:
		os.Exit(1)
	}
}
