package internal

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"
)

type Server struct {
	config      *Config
	handler     http.Handler
	httpServer  *http.Server
	httpsServer *http.Server
}

func NewServer(config *Config, handler http.Handler) *Server {
	return &Server{
		handler: handler,
		config:  config,
	}
}

func (s *Server) Start() {
	httpAddress := fmt.Sprintf(":%d", s.config.HttpPort)

	s.httpsServer = nil
	s.httpServer = s.defaultHttpServer(httpAddress)
	s.httpServer.Handler = s.handler

	go s.httpServer.ListenAndServe()

	slog.Info("Server started", "http", httpAddress)
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	defer slog.Info("Server stopped")

	slog.Info("Server stopping")

	s.httpServer.Shutdown(ctx)
	if s.httpsServer != nil {
		s.httpsServer.Shutdown(ctx)
	}
}

func (s *Server) defaultHttpServer(addr string) *http.Server {
	return &http.Server{
		Addr:         addr,
		IdleTimeout:  s.config.HttpIdleTimeout,
		ReadTimeout:  s.config.HttpReadTimeout,
		WriteTimeout: s.config.HttpWriteTimeout,
	}
}

func httpRedirectHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Connection", "close")

	host, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		host = r.Host
	}

	url := "https://" + host + r.URL.RequestURI()
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}
