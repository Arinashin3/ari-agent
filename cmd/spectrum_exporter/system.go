package main

import (
	"time"

	sdkMetric "go.opentelemetry.io/otel/sdk/metric"
)

func init() {
	registProvider("system", &systemProvider{})
}

func (pv *systemProvider) IsDefaultEnabled() bool {
	return true
}

func (pv *systemProvider) NewProvider(moduleName string, interval time.Duration, cl *ClientDesc) Provider {
	return &systemProvider{
		moduleName: moduleName,
		interval:   interval,
		//meterProvider: exp,
		clientDesc: cl,
	}
}

type systemProvider struct {
	moduleName    string
	interval      time.Duration
	meterProvider *sdkMetric.MeterProvider
	clientDesc    *ClientDesc
}

func registSystemProvider() (string, Provider) {
	moduleName := "system"

	return moduleName, &systemProvider{}
}
