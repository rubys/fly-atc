package internal

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"
)

var monitor_mutex sync.RWMutex
var registry []*Monitor

func NewMonitor(route *Route, config *Config, next http.Handler) *Monitor {
	monitor := &Monitor{route: route, config: config, next: next}

	if config.TargetPort == 0 {
		monitor.port = availablePort()
	} else {
		monitor.port = config.TargetPort
	}

	targetURL, err := url.Parse(fmt.Sprintf("http://localhost:%d", monitor.port))
	if err == nil {
		monitor.target = targetURL
	}

	monitor_mutex.Lock()
	registry = append(registry, monitor)
	monitor_mutex.Unlock()

	return monitor
}

type Monitor struct {
	sync.RWMutex
	route   *Route
	config  *Config
	next    http.Handler
	service *Service
	target  *url.URL
	port    int
}

func (m *Monitor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if m.service == nil {
		m.Lock()

		if m.service == nil {
			service := NewService(m.config)
			service.Start(m.route)
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
	monitor_mutex.Lock()
	defer monitor_mutex.Unlock()

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

func availablePort() int {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0
	}

	defer listener.Close()
	addr := listener.Addr().(*net.TCPAddr)
	return addr.Port
}
