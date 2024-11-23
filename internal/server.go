package internal

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	config     *Config
	handler    http.Handler
	httpServer *http.Server
}

func NewServer(config *Config) *Server {
	server := &Server{
		config: config,
	}

	handlerOptions := HandlerOptions{
		cache:                    server.cache(),
		config:                   config,
		xSendfileEnabled:         config.XSendfileEnabled,
		maxCacheableResponseBody: config.MaxCacheItemSizeBytes,
		maxRequestBody:           config.MaxRequestBody,
		badGatewayPage:           config.BadGatewayPage,
	}

	server.handler = NewHandler(handlerOptions)

	return server
}

func (s *Server) Start() {
	httpAddress := fmt.Sprintf(":%d", s.config.HttpPort)

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
}

func (s *Server) defaultHttpServer(addr string) *http.Server {
	return &http.Server{
		Addr:         addr,
		IdleTimeout:  s.config.HttpIdleTimeout,
		ReadTimeout:  s.config.HttpReadTimeout,
		WriteTimeout: s.config.HttpWriteTimeout,
	}
}

// Private

func (s *Server) cache() Cache {
	return NewMemoryCache(s.config.CacheSizeBytes, s.config.MaxCacheItemSizeBytes)
}
