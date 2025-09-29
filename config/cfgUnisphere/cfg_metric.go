package cfgUnisphere

import (
	"strconv"
	"time"
)

type UnisphereProviderMetric struct {
	Enabled  string   `yaml:"enabled,omitempty"`
	Paths    []string `yaml:"paths,omitempty"`
	Interval string   `yaml:"interval,omitempty"`
}

func (pv *UnisphereProviderMetric) GetEnabled(defaults bool) bool {
	if pv.Enabled == "" {
		return defaults
	}
	enabled, _ := strconv.ParseBool(pv.Enabled)
	return enabled
}

func (pv *UnisphereProviderMetric) GetInterval() time.Duration {
	interval, _ := time.ParseDuration(pv.Interval)
	return interval
}
