package internal

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
)

var router_mutex sync.RWMutex

var Instance = os.Getenv("FLY_MACHINE_ID")
var Region = os.Getenv("FLY_REGION")

type Router struct {
	next   http.Handler
	config *Config
}

func NewRouter(config *Config, next http.Handler) *Router {
	return &Router{
		config: config,
		next:   next,
	}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	routes := Routes()
	for i := 0; i < len(routes); i++ {
		route := &routes[i]

		if strings.HasPrefix(r.URL.Path, route.Endpoint) && (r.URL.Path == route.Endpoint || r.URL.Path[len(route.Endpoint)] == '/') {
			if route.Instance != "" && route.Instance != Instance {
				route.replay(w, r, "instance", route.Instance)
				return
			} else if route.Region != "" && route.Region != Region {
				route.replay(w, r, "region", route.Region)
				return
			}

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

// Note: this is a method of the Route struct, not Router
func (route *Route) replay(w http.ResponseWriter, _ *http.Request, field string, value string) {
	w.Header().Set("Fly-Replay", fmt.Sprintf("%s=%s", field, value))
	w.WriteHeader(http.StatusTemporaryRedirect)
}
