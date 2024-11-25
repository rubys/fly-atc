package internal

import (
	"errors"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	configFile = "config/atc.yml"

	KB = 1024
	MB = 1024 * KB

	ENV_PREFIX = "THRUSTER_"

	defaultTargetPort = 0

	defaultCacheSize             = 64 * MB
	defaultMaxCacheItemSizeBytes = 1 * MB
	defaultMaxRequestBody        = 0

	defaultBadGatewayPage = "./public/502.html"

	defaultHttpPort         = 3000
	defaultHttpIdleTimeout  = 60 * time.Second
	defaultHttpReadTimeout  = 30 * time.Second
	defaultHttpWriteTimeout = 30 * time.Second

	defaultLogLevel = slog.LevelInfo
)

type Config struct {
	TargetPort      int      `yaml:"target_port"`
	UpstreamCommand string   `yaml:"upstream_command"`
	UpstreamArgs    []string `yaml:"upstream_args"`

	CacheSizeBytes        int  `yaml:"cache_size_bytes"`
	MaxCacheItemSizeBytes int  `yaml:"max_cache_item_size_bytes"`
	XSendfileEnabled      bool `yaml:"x_sendfile_enabled"`
	MaxRequestBody        int  `yaml:"max_request_body"`

	BadGatewayPage string `yaml:"bad_gateway_page"`

	HttpPort         int           `yaml:"http_port"`
	HttpIdleTimeout  time.Duration `yaml:"http_idle_timeout"`
	HttpReadTimeout  time.Duration `yaml:"http_read_timeout"`
	HttpWriteTimeout time.Duration `yaml:"http_write_timeout"`

	LogLevel slog.Level `yaml:"log_level"`
}

type Route struct {
	Name     string `yaml:"name"`
	Endpoint string `yaml:"endpoint"`

	Database string `yaml:"database"`
	Region   string `yaml:"region"`
	Instance string `yaml:"instance"`

	Monitor *Monitor
}

type Settings struct {
	Server *Config `yaml:"server"`
	Routes []Route `yaml:"routes"`
}

var settings Settings

func Routes() []Route {
	return settings.Routes
}

func NewConfig() (*Config, error) {
	config := &Config{
		TargetPort: defaultTargetPort,

		CacheSizeBytes:        defaultCacheSize,
		MaxCacheItemSizeBytes: defaultMaxCacheItemSizeBytes,
		XSendfileEnabled:      true,
		MaxRequestBody:        defaultMaxRequestBody,

		BadGatewayPage: defaultBadGatewayPage,

		HttpPort:         defaultHttpPort,
		HttpIdleTimeout:  defaultHttpIdleTimeout,
		HttpReadTimeout:  defaultHttpReadTimeout,
		HttpWriteTimeout: defaultHttpWriteTimeout,

		LogLevel: defaultLogLevel,
	}

	settings.Server = config

	data, err := os.ReadFile(configFile)
	if err == nil {
		err = yaml.Unmarshal(data, &settings)
		if err != nil {
			return nil, err
		}
	}

	config.TargetPort = getEnvInt("TARGET_PORT", config.TargetPort)
	config.CacheSizeBytes = getEnvInt("CACHE_SIZE", config.CacheSizeBytes)
	config.MaxCacheItemSizeBytes = getEnvInt("MAX_CACHE_ITEM_SIZE", config.MaxCacheItemSizeBytes)
	config.XSendfileEnabled = getEnvBool("X_SENDFILE_ENABLED", config.XSendfileEnabled)
	config.MaxRequestBody = getEnvInt("MAX_REQUEST_BODY", config.MaxRequestBody)
	config.BadGatewayPage = getEnvString("BAD_GATEWAY_PAGE", config.BadGatewayPage)
	config.HttpPort = getEnvInt("HTTP_PORT", config.HttpPort)
	config.HttpIdleTimeout = getEnvDuration("HTTP_IDLE_TIMEOUT", config.HttpIdleTimeout)
	config.HttpReadTimeout = getEnvDuration("HTTP_READ_TIMEOUT", config.HttpReadTimeout)
	config.HttpWriteTimeout = getEnvDuration("HTTP_WRITE_TIMEOUT", config.HttpWriteTimeout)

	if getEnvBool("DEBUG", false) {
		config.LogLevel = slog.LevelDebug
	}

	if len(os.Args) >= 2 {
		config.UpstreamCommand = os.Args[1]
		config.UpstreamArgs = os.Args[2:]
	} else if config.UpstreamCommand == "" {
		return nil, errors.New("missing upstream command")
	}

	if len(settings.Routes) == 0 {
		settings.Routes = append(settings.Routes, Route{
			Name:     "",
			Endpoint: "",
		})
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
