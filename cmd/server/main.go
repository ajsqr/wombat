package main

import (
	"encoding/json"
	"flag"
	"log/slog"
	"os"

	"github.com/ajsqr/wombat/server"
)

func main() {
	configFilePath := flag.String("config", "", "path to wombat-server config file")
	flag.Parse()

	content, err := os.ReadFile(*configFilePath)
	if err != nil {
		panic(err)
	}

	var serverConfig server.ServerConfig
	err = json.Unmarshal(content, &serverConfig)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	wa := server.NewServer(&serverConfig, logger)
	wa.Run()
}
