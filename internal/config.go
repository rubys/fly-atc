package internal

import (
	"errors"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	KB = 1024
	MB = 1024 * KB

	ENV_PREFIX = "THRUSTER_"

	defaultTargetPort = 3000

	defaultCacheSize             = 64 * MB
	defaultMaxCacheItemSizeBytes = 1 * MB
	defaultMaxRequestBody        = 0

	defaultBadGatewayPage = "./public/502.html"

	defaultHttpPort         = 80
	defaultHttpIdleTimeout  = 60 * time.Second
	defaultHttpReadTimeout  = 30 * time.Second
	defaultHttpWriteTimeout = 30 * time.Second

	defaultLogLevel = slog.LevelInfo
)

type Config struct {
	TargetPort      int
	UpstreamCommand string
	UpstreamArgs    []string

	CacheSizeBytes        int
	MaxCacheItemSizeBytes int
	XSendfileEnabled      bool
	MaxRequestBody        int

	BadGatewayPage string

	HttpPort         int
	HttpIdleTimeout  time.Duration
	HttpReadTimeout  time.Duration
	HttpWriteTimeout time.Duration

	LogLevel slog.Level
}

func NewConfig() (*Config, error) {
	if len(os.Args) < 2 {
		return nil, errors.New("missing upstream command")
	}

	logLevel := defaultLogLevel
	if getEnvBool("DEBUG", false) {
		logLevel = slog.LevelDebug
	}

	config := &Config{
		TargetPort:      getEnvInt("TARGET_PORT", defaultTargetPort),
		UpstreamCommand: os.Args[1],
		UpstreamArgs:    os.Args[2:],

		CacheSizeBytes:        getEnvInt("CACHE_SIZE", defaultCacheSize),
		MaxCacheItemSizeBytes: getEnvInt("MAX_CACHE_ITEM_SIZE", defaultMaxCacheItemSizeBytes),
		XSendfileEnabled:      getEnvBool("X_SENDFILE_ENABLED", true),
		MaxRequestBody:        getEnvInt("MAX_REQUEST_BODY", defaultMaxRequestBody),

		BadGatewayPage: getEnvString("BAD_GATEWAY_PAGE", defaultBadGatewayPage),

		HttpPort:         getEnvInt("HTTP_PORT", defaultHttpPort),
		HttpIdleTimeout:  getEnvDuration("HTTP_IDLE_TIMEOUT", defaultHttpIdleTimeout),
		HttpReadTimeout:  getEnvDuration("HTTP_READ_TIMEOUT", defaultHttpReadTimeout),
		HttpWriteTimeout: getEnvDuration("HTTP_WRITE_TIMEOUT", defaultHttpWriteTimeout),

		LogLevel: logLevel,
	}

	return config, nil
}

func findEnv(key string) (string, bool) {
	value, ok := os.LookupEnv(ENV_PREFIX + key)
	if ok {
		return value, true
	}

	value, ok = os.LookupEnv(key)
	if ok {
		return value, true
	}

	return "", false
}

func getEnvString(key, defaultValue string) string {
	value, ok := findEnv(key)
	if ok {
		return value
	}

	return defaultValue
}

func getEnvStrings(key string, defaultValue []string) []string {
	value, ok := findEnv(key)
	if ok {
		items := strings.Split(value, ",")
		result := []string{}

		for _, item := range items {
			item = strings.TrimSpace(item)
			if item != "" {
				result = append(result, item)
			}
		}

		return result
	}

	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	value, ok := findEnv(key)
	if !ok {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	value, ok := findEnv(key)
	if !ok {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return time.Duration(intValue) * time.Second
}

func getEnvBool(key string, defaultValue bool) bool {
	value, ok := findEnv(key)
	if !ok {
		return defaultValue
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return boolValue
}
