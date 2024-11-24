package internal

import (
	"log/slog"
	"net/http"

	"github.com/klauspost/compress/gzhttp"
)

type HandlerOptions struct {
	badGatewayPage           string
	cache                    Cache
	config                   *Config
	maxCacheableResponseBody int
	maxRequestBody           int
	xSendfileEnabled         bool
	forwardHeaders           bool
}

func NewHandler(options HandlerOptions) http.Handler {
	handler := NewProxyHandler(options.badGatewayPage, options.forwardHeaders)
	handler = NewCacheHandler(options.cache, options.maxCacheableResponseBody, handler)
	handler = NewSendfileHandler(options.xSendfileEnabled, handler)
	handler = gzhttp.GzipHandler(handler)

	if options.maxRequestBody > 0 {
		handler = http.MaxBytesHandler(handler, int64(options.maxRequestBody))
	}

	handler = NewLoggingMiddleware(slog.Default(), handler)

	handler = NewMonitor("__default__", options.config, handler)

	return handler
}
