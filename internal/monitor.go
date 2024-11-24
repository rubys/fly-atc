package internal

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"
)

var mutex sync.RWMutex
var registry []*Monitor

func NewMonitor(name string, config *Config, next http.Handler) *Monitor {
	monitor := &Monitor{config: config, next: next}

	if config.TargetPort == 0 {
		config.TargetPort = monitor.availablePort()
	}

	targetURL, err := url.Parse(fmt.Sprintf("http://localhost:%d", config.TargetPort))
	if err == nil {
		monitor.target = targetURL
	}

	mutex.Lock()
	registry = append(registry, monitor)
	mutex.Unlock()

	return monitor
}

type Monitor struct {
	sync.RWMutex
	config  *Config
	next    http.Handler
	service *Service
	target  *url.URL
	started bool
}

func (m *Monitor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if m.service == nil {
		m.Lock()

		if m.service == nil {
			service := NewService(m.config)
			service.Start()
			service.HealthCheck(m.target.String())
			m.service = service
		}

		m.Unlock()
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, "target_url", m.target)
	newReq := r.WithContext(ctx)

	m.next.ServeHTTP(w, newReq)
}

func Shutdown() int {
	mutex.Lock()
	defer mutex.Unlock()

	for _, monitor := range registry {
		monitor.Lock()
		if monitor.service != nil {
			monitor.service.Stop()
		}
		monitor.Unlock()
	}

	return 0
}

// private

func (m *Monitor) availablePort() int {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0
	}

	defer listener.Close()
	addr := listener.Addr().(*net.TCPAddr)
	return addr.Port
}
