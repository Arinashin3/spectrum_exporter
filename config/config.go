package config

import (
	"errors"
	"os"
	"reflect"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

type SpectrumConfig struct {
	Global    *GlobalConfig      `yaml: "global,omitempty"`
	Server    *ServerConfig      `yaml: "server,omitempty"`
	Clients   []*ClientConfig    `yaml: "targets,omitempty"`
	Auths     []*AuthConfig      `yaml: "auths,omitempty"`
	Providers *SpectrumProviders `yaml: "providers,omitempty"`
}

type SpectrumProviders struct {
	System       *ProviderDefaults `yaml: "system,omitempty"`
	Systemstats  *ProviderDefaults `yaml: "systemstats,omitempty"`
	Eventlog     *ProviderDefaults `yaml: "event,omitempty"`
	Flashcopy    *ProviderDefaults `yaml: "flashcopy,omitempty"`
	Enclosure    *ProviderDefaults `yaml: "enclosure,omitempty"`
	Canister     *ProviderDefaults `yaml: "canister,omitempty"`
	Mdisk        *ProviderDefaults `yaml: "mdisk,omitempty"`
	Vdisk        *ProviderDefaults `yaml: "vdisk,omitempty"`
	Array        *ProviderDefaults `yaml: "array,omitempty"`
	Drive        *ProviderDefaults `yaml: "drive,omitempty"`
	Host         *ProviderDefaults `yaml: "host,omitempty"`
	Node         *ProviderDefaults `yaml: "node,omitempty"`
	Hostvdiskmap *ProviderDefaults `yaml: "hostvdiskmap,omitempty"`
}

func NewConfiguration() *SpectrumConfig {
	return &SpectrumConfig{
		Global: &GlobalConfig{
			Server: &GlobalServerConfig{
				Endpoint: "http://127.0.0.1:8080",
				Api_Path: "",
				Insecure: false,
				Mode:     "http",
			},
			Client: &GlobalClientConfig{
				Auth:     "",
				Insecure: false,
			},
			Provider: &GlobalProviderConfig{
				Interval: "1m",
			},
		},
		Server: &ServerConfig{
			Metrics: &ServerMetricConfig{
				Enabled: true,
			},
			Logs: &ServerLogConfig{
				Enabled: true,
			},
			Traces: &ServerTraceConfig{
				Enabled: true,
			},
		},
		Clients: nil,
		Auths:   nil,
		Providers: &SpectrumProviders{
			System:       &ProviderDefaults{},
			Systemstats:  &ProviderDefaults{},
			Eventlog:     &ProviderDefaults{},
			Flashcopy:    &ProviderDefaults{},
			Enclosure:    &ProviderDefaults{},
			Canister:     &ProviderDefaults{},
			Mdisk:        &ProviderDefaults{},
			Array:        &ProviderDefaults{},
			Drive:        &ProviderDefaults{},
			Host:         &ProviderDefaults{},
			Node:         &ProviderDefaults{},
			Hostvdiskmap: &ProviderDefaults{},
			Vdisk:        &ProviderDefaults{},
		},
	}
}

func (cfg *SpectrumConfig) LoadFile(file *string) error {
	ymlContents, err := os.ReadFile(*file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(ymlContents, cfg)
	if err != nil {
		return err
	}

	err = cfg.applyGlobal()
	if err != nil {
		return err
	}

	return err
}

// applyGlobal
// Section 내용이 비어있을 경우,
// Global 설정을 각각의 Section에 적용
func (cfg *SpectrumConfig) applyGlobal() error {
	// Set Client
	g := cfg.Global
	if cfg.Clients == nil {
		return errors.New("no clients configured")
	}
	for _, c := range cfg.Clients {
		if c.Endpoint == "" {
			return errors.New("client endpoint is required")
		}
		if c.Auth == "" {
			c.Auth = g.Client.Auth
		}
		if c.Insecure == "" {
			c.Insecure = strconv.FormatBool(g.Client.Insecure)
		}
		for k, v := range g.Client.Labels {
			if c.Labels[k] == "" {
				c.Labels[k] = v
			}
		}
	}
	// Set Global config at Servers
	svNum := reflect.ValueOf(cfg.Server).Elem().NumField()
	for i := 0; i < svNum; i++ {
		sv := reflect.ValueOf(cfg.Server).Elem().Field(i).Elem()
		endpoint := sv.FieldByName("Endpoint")
		if endpoint.String() == "" {
			endpoint.SetString(g.Server.Endpoint)
		}
		apiPath := sv.FieldByName("Api_Path")
		if apiPath.String() == "" {
			apiPath.SetString(g.Server.Api_Path)
		}
		insecure := sv.FieldByName("Insecure")
		if insecure.String() == "" {
			insecure.SetString(strconv.FormatBool(g.Server.Insecure))
		}
		mode := sv.FieldByName("Mode")
		if mode.String() == "" {
			mode.SetString(g.Server.Mode)
		}
	}

	//Check Error to parse global provider interval
	_, err := time.ParseDuration(g.Provider.Interval)
	if err != nil {
		return err
	}
	//At the Providers, if value of field is null, then Apply Global
	pvNum := reflect.ValueOf(cfg.Providers).Elem().NumField()
	for i := 0; i < pvNum; i++ {
		pv := reflect.ValueOf(cfg.Providers).Elem().Field(i).Elem()

		//Apply Global Interval
		interval := pv.FieldByName("Interval")
		if interval.String() == "" {
			interval.SetString(cfg.Global.Provider.Interval)
		} else {
			_, err = time.ParseDuration(interval.String())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// SearchAuth
// 인증정보를 찾아, base64로 인코딩하여 리턴합니다.
func (cfg *SpectrumConfig) SearchAuth(name string) (string, string) {
	for _, auth := range cfg.Auths {
		if auth.Name == name {
			return auth.User, auth.Password
		}
	}
	return "", ""
}

func (cfg *SpectrumConfig) GetConfig() *SpectrumConfig {
	return cfg
}
