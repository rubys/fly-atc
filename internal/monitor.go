package internal

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"
)

var mutex sync.RWMutex

func NewMonitor(name string, config *Config, next http.Handler) *Monitor {
	monitor := &Monitor{config: config, next: next}

	mutex.Lock()
	new_registry := make(map[string]*Monitor)
	for k, v := range registry {
		new_registry[k] = v
	}
	new_registry[name] = monitor
	registry = new_registry
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

var registry = make(map[string]*Monitor)

func (m *Monitor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Lock()

	if m.config.TargetPort == 0 {
		m.config.TargetPort = m.availablePort()
	}

	targetURL, err := url.Parse(fmt.Sprintf("http://localhost:%d", m.config.TargetPort))
	if err != nil {
		m.Unlock()
		http.Error(w, "Invalid target URL", http.StatusInternalServerError)
		return
	}
	m.target = targetURL

	if m.service == nil {
		m.service = NewService(m.config)
		m.service.Start()
		m.Unlock()
		m.service.HealthCheck(m.target.String())
	} else {
		m.Unlock()
	}

	m.next.ServeHTTP(w, r)
}

func TargetForMonitor(name string) *url.URL {
	if monitor, ok := registry[name]; ok {
		return monitor.target
	}

	return nil
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
