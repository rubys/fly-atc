package internal

import (
	"context"
	"net/http"
	"strings"
	"sync"
)

var router_mutex sync.RWMutex

type Router struct {
	next   http.Handler
	config *Config
}

func NewRouter(config *Config, next http.Handler) *Router {
	return &Router{config: config, next: next}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	routes := Routes()
	for i := 0; i < len(routes); i++ {
		route := &routes[i]

		if strings.HasPrefix(r.URL.Path, route.Endpoint) && (r.URL.Path == route.Endpoint || r.URL.Path[len(route.Endpoint)] == '/') {
			if route.Monitor == nil {
				router_mutex.Lock()

				if route.Monitor == nil {
					route.Monitor = NewMonitor(route, router.config, router.next)
				}

				router_mutex.Unlock()
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, "fly_atc_scope", route.Endpoint)
			r = r.WithContext(ctx)

			route.Monitor.ServeHTTP(w, r)
			return
		}
	}

	w.WriteHeader(404)
	w.Write([]byte(`404 not found`))
}
