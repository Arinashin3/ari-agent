package config

type ConfigYaml struct {
	Global    *GlobalYaml    `yaml: "global,omitempty"`
	Server    *ServerYaml    `yaml: "server,omitempty"`
	Clients   []*ClientYaml  `yaml: "targets,omitempty"`
	Auths     []*AuthYaml    `yaml: "auths,omitempty"`
	Providers *ProvidersYaml `yaml: "providers,omitempty"`
}

type GlobalYaml struct {
	Server   *GlobalServerYaml   `yaml: "server,omitempty"`
	Client   *GlobalClientYaml   `yaml: "client,omitempty"`
	Provider *GlobalProviderYaml `yaml: "provider,omitempty"`
}

type GlobalServerYaml struct {
	Endpoint string `yaml: "endpoint"`
	Api_Path string `yaml: "api_path"`
	Mode     string `yaml: "mode,omitempty"`
	Insecure bool   `yaml: "insecure"`
}

type GlobalClientYaml struct {
	Auth     string            `yaml: "auth"`
	Insecure bool              `yaml: "insecure,omitempty"`
	Labels   map[string]string `yaml: "labels,omitempty"`
}

type GlobalProviderYaml struct {
	Interval string `yaml: "interval"`
}

type ServerYaml struct {
	Metrics *ServerMetricYaml `yaml: "metrics,omitempty"`
	Logs    *ServerLogYaml    `yaml: "logs,omitempty"`
	Traces  *ServerTraceYaml  `yaml: "traces,omitempty"`
}

type ServerMetricYaml struct {
	Endpoint string `yaml: "endpoint,omitempty"`
	Api_Path string `yaml: "api_path,omitempty"`
	Mode     string `yaml: "mode,omitempty"`
	Insecure string `yaml: "insecure,omitempty"`
	Enabled  bool   `yaml: "enabled,omitempty"`
}

type ServerLogYaml struct {
	Endpoint string `yaml: "endpoint,omitempty"`
	Api_Path string `yaml: "api_path,omitempty"`
	Mode     string `yaml: "mode,omitempty"`
	Insecure string `yaml: "insecure,omitempty"`
	Enabled  bool   `yaml: "enabled,omitempty"`
}

type ServerTraceYaml struct {
	Endpoint string `yaml: "endpoint,omitempty"`
	Api_Path string `yaml: "api_path,omitempty"`
	Mode     string `yaml: "mode,omitempty"`
	Insecure string `yaml: "insecure,omitempty"`
	Enabled  bool   `yaml: "enabled,omitempty"`
}

type ClientYaml struct {
	Endpoint string            `yaml: "endpoint"`
	Auth     string            `yaml: "auth,omitempty"`
	Insecure string            `yaml: "insecure,omitempty"`
	Labels   map[string]string `yaml: "labels,omitempty"`
}

type AuthYaml struct {
	Name     string `yaml: "name"`
	User     string `yaml: "user"`
	Password string `yaml: "password"`
}
type ProviderYaml struct {
	Module   string `yaml: "module,omitempty"`
	Interval string `yaml: "interval,omitempty"`
}
type ProvidersYaml struct {
	System   *ProviderSystem   `yaml: "system,omitempty"`
	Capacity *ProviderCapacity `yaml: "capacity,omitempty"`
	Metric_A *ProviderMetric   `yaml: "metric_a,omitempty"`
	Metric_B *ProviderMetric   `yaml: "metric_b,omitempty"`
	Metric_C *ProviderMetric   `yaml: "metric_c,omitempty"`
	Event    *ProviderEvent    `yaml: "event,omitempty"`
	Lun      *ProviderLun      `yaml: "lun,omitempty"`
}

type ProviderSystem struct {
	Enabled  bool   `yaml: "enabled,omitempty"`
	Interval string `yaml: "interval,omitempty"`
}

type ProviderCapacity struct {
	Enabled  bool   `yaml: "enabled,omitempty"`
	Interval string `yaml: "interval,omitempty"`
}

type ProviderMetric struct {
	Enabled  bool     `yaml: "enabled,omitempty"`
	Interval string   `yaml: "interval,omitempty"`
	Paths    []string `yaml: "paths,omitempty"`
}
type ProviderEvent struct {
	Enabled  bool   `yaml: "enabled,omitempty"`
	Interval string `yaml: "interval,omitempty"`
	Level    int    `yaml: "level,omitempty"`
}
type ProviderLun struct {
	Enabled  bool   `yaml: "enabled,omitempty"`
	Interval string `yaml: "interval,omitempty"`
}
