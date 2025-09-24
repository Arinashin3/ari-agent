package cfgSpectrum

import (
	"strconv"
	"time"
)

type SpectrumProviderLsSystem struct {
	Enabled  string `yaml: "enabled,omitempty"`
	Interval string `yaml: "interval,omitempty"`
}

func (sp *SpectrumProviderLsSystem) GetEnabled(defaults bool) bool {
	if sp.Enabled == "" {
		return defaults
	}
	enabled, _ := strconv.ParseBool(sp.Enabled)
	return enabled
}

func (sp *SpectrumProviderLsSystem) GetInterval() time.Duration {
	interval, _ := time.ParseDuration(sp.Interval)
	return interval
}
