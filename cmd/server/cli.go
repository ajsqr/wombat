package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ajsqr/wombat/server"
)

type cli struct {
	server  *server.Server
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
		fmt.Printf("wombat-server %s build %s\n", c.version, c.commit)
	case "run":
		c.server.Run()
	default:
		c.logger.Error("unexpected cli command", slog.String("cmd", baseCmd))
		os.Exit(1)
	}
}
