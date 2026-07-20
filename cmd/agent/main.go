package main

import (
	"log/slog"
	"os"

	"github.com/ajsqr/wombat/agent"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	wa, err := agent.NewAgent(logger)
	if err != nil {
		logger.Error("error initializing agent", slog.Any("error", err))
		os.Exit(1)
	}

	wa.Run()
}
