package config

import (
	"strconv"
	"time"
)

type Config interface {
	LoadFile(file *string) error
	GetMetricsMode() string
	GetMetricsEndpoint() string
	GetMetricsInsecure() bool
	GetLogsMode() string
	GetLogsEndpoint() string
	GetLogsInsecure() bool
	GetClientList() []*ClientConfig
	GetProviderSystem() any
}

type GlobalConfig struct {
	Server   *GlobalServerConfig   `yaml: "server,omitempty"`
	Client   *GlobalClientConfig   `yaml: "client,omitempty"`
	Provider *GlobalProviderConfig `yaml: "provider,omitempty"`
}

type GlobalServerConfig struct {
	Endpoint string `yaml: "endpoint"`
	Api_Path string `yaml: "api_path"`
	Mode     string `yaml: "mode,omitempty"`
	Insecure bool   `yaml: "insecure"`
}

type GlobalClientConfig struct {
	Auth     string            `yaml: "auth"`
	Insecure bool              `yaml: "insecure,omitempty"`
	Labels   map[string]string `yaml: "labels,omitempty"`
}

type GlobalProviderConfig struct {
	Interval string `yaml: "interval"`
}

type ServerConfig struct {
	Metrics *ServerMetricConfig `yaml: "metrics,omitempty"`
	Logs    *ServerLogConfig    `yaml: "logs,omitempty"`
	Traces  *ServerTraceConfig  `yaml: "traces,omitempty"`
}

type ServerMetricConfig struct {
	Endpoint string `yaml: "endpoint,omitempty"`
	Api_Path string `yaml:"api_path,omitempty"`
	Mode     string `yaml: "mode,omitempty"`
	Insecure string `yaml: "insecure,omitempty"`
	Enabled  bool   `yaml: "enabled,omitempty"`
}

type ServerLogConfig struct {
	Endpoint string `yaml: "endpoint,omitempty"`
	Api_Path string `yaml:"api_path,omitempty"`
	Mode     string `yaml: "mode,omitempty"`
	Insecure string `yaml: "insecure,omitempty"`
	Enabled  bool   `yaml: "enabled,omitempty"`
}

type ServerTraceConfig struct {
	Endpoint string `yaml: "endpoint,omitempty"`
	Api_Path string `yaml:"api_path,omitempty"`
	Mode     string `yaml: "mode,omitempty"`
	Insecure string `yaml: "insecure,omitempty"`
	Enabled  bool   `yaml: "enabled,omitempty"`
}

type ClientConfig struct {
	Endpoint string            `yaml: "endpoint"`
	Auth     string            `yaml: "auth,omitempty"`
	Insecure string            `yaml: "insecure,omitempty"`
	Labels   map[string]string `yaml: "labels,omitempty"`
}

type AuthConfig struct {
	Name     string `yaml: "name"`
	User     string `yaml: "user"`
	Password string `yaml: "password"`
}

// Providers...
type ProviderSystem struct {
	Enabled  bool   `yaml: "enabled,omitempty"`
	Interval string `yaml: "interval,omitempty"`
}

type ProviderCapacity struct {
	Enabled  bool   `yaml: "enabled,omitempty"`
	Interval string `yaml: "interval,omitempty"`
}

type ProviderLun struct {
	Enabled  bool   `yaml: "enabled,omitempty"`
	Interval string `yaml: "interval,omitempty"`
}

type ProviderDefaults struct {
	Enabled  string `yaml:"enabled,omitempty"`
	Interval string `yaml:"interval,omitempty"`
}

func (pv *ProviderDefaults) GetEnabled(defaults bool) bool {
	if pv.Enabled == "" {
		return defaults
	}
	enabled, _ := strconv.ParseBool(pv.Enabled)
	return enabled
}

func (pv *ProviderDefaults) GetInterval() time.Duration {
	interval, _ := time.ParseDuration(pv.Interval)
	return interval
}
