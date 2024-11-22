package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/rubys/fly-atc/internal"
)

func setLogger(level slog.Level) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})))
}

func main() {
	config, err := internal.NewConfig()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}

	setLogger(config.LogLevel)

	server := internal.NewServer(config)
	server.Start()
	defer server.Stop()

	service := internal.NewService(config)
	service.Start()
	os.Exit(service.Stop())
}
