package internal

import (
	"fmt"
	"net/http"
	"sync"
)

var mutex sync.RWMutex

func NewMonitor(config *Config, next http.Handler) *Monitor {
	monitor := &Monitor{config: config, next: next}
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
	started bool
}

var registry = make([]*Monitor, 0)

func (m *Monitor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Lock()
	if m.service == nil {
		m.service = NewService(m.config)
		m.service.Start()
		m.Unlock()
		m.service.HealthCheck(fmt.Sprintf("http://localhost:%d", m.config.TargetPort))
	} else {
		m.Unlock()
	}

	m.next.ServeHTTP(w, r)
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
