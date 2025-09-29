package cfgUnisphere

import (
	"strconv"
	"time"
)

type UnisphereProviderEvent struct {
	Enabled  string `yaml:"enabled,omitempty"`
	Level    int    `yaml:"level,omitempty"`
	Interval string `yaml:"interval,omitempty"`
}

func (pv *UnisphereProviderEvent) GetEnabled(defaults bool) bool {
	if pv.Enabled == "" {
		return defaults
	}
	enabled, _ := strconv.ParseBool(pv.Enabled)
	return enabled
}

func (pv *UnisphereProviderEvent) GetInterval() time.Duration {
	interval, _ := time.ParseDuration(pv.Interval)
	return interval
}
