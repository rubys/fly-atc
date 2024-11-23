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
		for i := 0; i < 240; i++ {
			time.Sleep(250 * time.Millisecond)
			response, err := http.Get(fmt.Sprintf("http://localhost:%d", s.config.TargetPort))
			if err == nil && response != nil && response.StatusCode == 200 {
				alive <- nil
				return
			}
		}

		response, err := http.Get(fmt.Sprintf("http://localhost:%d", s.config.TargetPort))
		if err != nil {
			alive <- err
		}

		if response.StatusCode != 200 {
			alive <- fmt.Errorf("unexpected status code: %d", response.StatusCode)
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
