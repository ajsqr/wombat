package main

import (
	"log/slog"
	"os"

	"github.com/ajsqr/wombat/server"
)

var (
	commit  = "commit"
	version = "version"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	wa, err := server.NewServer(logger)
	if err != nil {
		logger.Error("error initializing server", slog.Any("error", err))
		os.Exit(1)
	}

	serverCli := cli{
		logger:  logger,
		server:  wa,
		commit:  commit,
		version: version,
	}

	serverCli.run()
}
