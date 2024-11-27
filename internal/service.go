package internal

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type Service struct {
	config   *Config
	upstream *UpstreamProcess
}

func NewService(config *Config) *Service {
	return &Service{
		config:   config,
		upstream: nil,
	}
}

func (s *Service) Start(route *Route) error {
	database_url := os.Getenv("DATABASE_URL")

	if database_url == "" {
		environment := os.Getenv("RAILS_ENV")

		if environment == "" {
			environment = "development"
		}

		database_url = fmt.Sprintf("sqlite3:./storage/%s.sqlite3", environment)
	}

	if route.Database != "" {
		dir := filepath.Dir(database_url)
		database_url = filepath.Join(dir, fmt.Sprintf("%s.sqlite3", route.Database))
	}

	litestream_config := fmt.Sprintf("tmp/litestream_%s.yml", route.Name)

	cmd := exec.Command("bin/rails", "atc:prepare")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("LITESTREAM_CONFIG=%s", litestream_config),
		fmt.Sprintf("DATABASE_URL=%s", database_url),
	)
	err := cmd.Run()
	if err != nil {
		slog.Error("Failed to prepare atc environment", "error", err)
		return err
	}

	s.upstream = NewUpstreamProcess(s.config.UpstreamCommand, s.config.UpstreamArgs...)

	s.upstream.setEnvironment("PORT", fmt.Sprintf("%d", route.Monitor.port))
	s.upstream.setEnvironment("FLY_ATC_SCOPE", route.Endpoint)
	s.upstream.setEnvironment("LITESTREAM_CONFIG", litestream_config)
	s.upstream.setEnvironment("PIDFILE", fmt.Sprintf("tmp/pids/%s.pid", route.Name))
	s.upstream.setEnvironment("DATABASE_URL", database_url)

	err = s.upstream.Start()
	if err != nil {
		slog.Error("Failed to start wrapped process", "command", s.config.UpstreamCommand, "args", s.config.UpstreamArgs, "error", err)
		return err
	}

	return nil
}

func (s *Service) HealthCheck(endpoint string) error {
	alive := make(chan error)

	go func() {
		stop := time.Now().Add(s.config.HttpIdleTimeout)

		for {
			time.Sleep(250 * time.Millisecond)

			response, err := http.Get(endpoint)
			if err == nil && response != nil && response.StatusCode == 200 {
				alive <- nil
				return
			}

			if time.Now().After(stop) {
				if err != nil {
					alive <- err
				} else if response == nil {
					alive <- fmt.Errorf("no response")
				} else {
					alive <- fmt.Errorf("unexpected status code: %d", response.StatusCode)
				}

				break
			}
		}
	}()

	return <-alive
}

func (s *Service) Stop() int {
	exitCode, err := s.upstream.Stop()
	if err != nil {
		slog.Error("Failed to stop wrapped process", "command", s.config.UpstreamCommand, "args", s.config.UpstreamArgs, "error", err)
		return 1
	}

	return exitCode
}
