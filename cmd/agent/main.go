package main

import (
	"encoding/json"
	"flag"
	"log/slog"
	"os"

	"github.com/ajsqr/wombat/agent"
)

func main() {
	configFilePath := flag.String("config", "", "path to wombat-agent config file")
	flag.Parse()

	content, err := os.ReadFile(*configFilePath)
	if err != nil {
		panic(err)
	}

	var agentConfig agent.AgentConfig
	err = json.Unmarshal(content, &agentConfig)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	wa := agent.NewAgent(&agentConfig, logger)
	wa.Run()
}
