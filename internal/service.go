package internal

import (
	"fmt"
	"log/slog"
	"net/http"
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

func (s *Service) Start() error {
	s.upstream = NewUpstreamProcess(s.config.UpstreamCommand, s.config.UpstreamArgs...)

	s.upstream.setEnvironment("PORT", fmt.Sprintf("%d", s.config.TargetPort))
	fmt.Printf("PORT: %d\n", s.config.TargetPort)

	err := s.upstream.Start()
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
