package config

import (
	"encoding/base64"
	"errors"
	"log/slog"
	"os"
	"reflect"
	"strconv"

	"github.com/Arinashin3/ari-agent/flag"

	"gopkg.in/yaml.v3"
)

var cfg *ConfigYaml
var logger *slog.Logger

func init() {
	var err error
	cfg = newConfiguration()
	logger = flag.Logger

	logger.Debug("Load configs...")

	// Read Config file and set
	err = cfg.LoadFile(flag.ConfigFile)
	if err != nil {
		logger.Error("Failed to read config file", "error", err.Error())
		os.Exit(1)
	}

	// Apply Global Settings...
	err = cfg.applyGlobal()
	if err != nil {
		logger.Error("Failed to set global configs", "error", err.Error())
		os.Exit(1)
	}

	// New Exporter
}

func newConfiguration() *ConfigYaml {
	return &ConfigYaml{
		Global: &GlobalYaml{
			Server: &GlobalServerYaml{
				Endpoint: "http://127.0.0.1:8080",
				Api_Path: "",
				Insecure: false,
				Mode:     "http",
			},
			Client: &GlobalClientYaml{
				Auth:     "",
				Insecure: false,
			},
			Provider: &GlobalProviderYaml{
				Interval: "1m",
			},
		},
		Server: &ServerYaml{
			Metrics: &ServerMetricYaml{
				Enabled: true,
			},
			Logs: &ServerLogYaml{
				Enabled: true,
			},
			Traces: &ServerTraceYaml{
				Enabled: true,
			},
		},
		Clients: nil,
		Auths:   nil,
		Providers: &ProvidersYaml{
			System: &ProviderSystem{
				Enabled: true,
			},
			Capacity: &ProviderCapacity{
				Enabled: true,
			},
			Metric_A: &ProviderMetric{
				Enabled: false,
			},
			Metric_B: &ProviderMetric{
				Enabled: false,
			},
			Metric_C: &ProviderMetric{
				Enabled: false,
			},
			Event: &ProviderEvent{
				Enabled: true,
				Level:   5,
			},
			Lun: &ProviderLun{
				Enabled: false,
			},
		},
	}
}

func (cfg *ConfigYaml) LoadFile(file *string) error {
	ymlContents, err := os.ReadFile(*file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(ymlContents, &cfg)
	return err
}

// applyGlobal
// Section 내용이 비어있을 경우,
// Global 설정을 각각의 Section에 적용
func (cfg *ConfigYaml) applyGlobal() error {
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

	// Set Global config at Providers
	pvNum := reflect.ValueOf(cfg.Providers).Elem().NumField()
	for i := 0; i < pvNum; i++ {
		pv := reflect.ValueOf(cfg.Providers).Elem().Field(i).Elem()

		// Check Enabled
		enabled := pv.FieldByName("Enabled").Bool()
		if !enabled {
			continue
		}
		interval := pv.FieldByName("Interval")
		if interval.String() == "" {
			interval.SetString(cfg.Global.Provider.Interval)
		}
	}
	return nil
}

// SearchAuth
// 인증정보를 찾아, base64로 인코딩하여 리턴합니다.
func (cfg *ConfigYaml) SearchAuth(name string) string {
	for _, auth := range cfg.Auths {
		if auth.Name == name {
			return base64.StdEncoding.EncodeToString([]byte(auth.User + ":" + auth.Password))
		}
	}
	return ""
}

func GetConfig() *ConfigYaml {
	return cfg
}

func GetProviderSystem() *ProviderSystem {
	return cfg.Providers.System
}

func GetProviderCapacity() *ProviderCapacity {
	return cfg.Providers.Capacity
}

func GetProviderMetricA() *ProviderMetric {
	return cfg.Providers.Metric_A
}
func GetProviderMetricB() *ProviderMetric {
	return cfg.Providers.Metric_B
}
func GetProviderMetricC() *ProviderMetric {
	return cfg.Providers.Metric_C
}
func GetProviderEvent() *ProviderEvent {
	return cfg.Providers.Event
}
func GetProviderLun() *ProviderLun {
	return cfg.Providers.Lun
}

func GetMetricsEndpoint() string {
	if cfg.Server.Metrics.Enabled {
		return cfg.Server.Metrics.Endpoint + cfg.Server.Metrics.Api_Path
	}
	return ""
}

func GetMetricsMode() string {
	if cfg.Server.Metrics.Enabled {
		return cfg.Server.Metrics.Mode
	}
	return ""
}

func GetMetricsInsecure() bool {
	insecure, _ := strconv.ParseBool(cfg.Server.Metrics.Insecure)
	return insecure
}

func GetLogsEndpoint() string {
	if cfg.Server.Logs.Enabled {
		return cfg.Server.Logs.Endpoint + cfg.Server.Logs.Api_Path
	}
	return ""
}
func GetLogsMode() string {
	if cfg.Server.Logs.Enabled {
		return cfg.Server.Logs.Mode
	}
	return ""
}

func GetLogsInsecure() bool {
	insecure, _ := strconv.ParseBool(cfg.Server.Logs.Insecure)
	return insecure
}
