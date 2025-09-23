package collector

import (
	"context"
	"time"

	"github.com/Arinashin3/ari-agent/cmd/unisphere_exporter/config"
	"github.com/Arinashin3/ari-agent/utils/provider"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	otlpmetric "go.opentelemetry.io/otel/sdk/metric"
)

func init() {
	InitClientList()
	for _, cl := range ClientList {
		registMeterProvider(newSystemProvider(cl))

	}
}

type systemProvider struct {
	moduleName    string
	interval      time.Duration
	meterProvider *otlpmetric.MeterProvider
	host          *ClientResourceStruct
}

func newSystemProvider(cl *ClientResourceStruct) MeterProvider {
	moduleName := "system"
	if MeterExporter == nil {
		return nil
	}
	pvConf := config.GetProviderSystem()
	if !pvConf.Enabled {
		return nil
	}

	interval, err := time.ParseDuration(pvConf.Interval)
	if err != nil {
		logger.Error("Failed to parse interval", "provider", moduleName, "error", err)
	}

	mp := provider.NewMeterProvider(serviceName, interval, MeterExporter)

	return &systemProvider{
		moduleName:    moduleName,
		interval:      interval,
		meterProvider: mp,
		host:          cl,
	}
}

var systemMetricList = []*MetricDescriptor{
	{
		key:      "info",
		name:     "unisphere_system_info",
		desc:     "Information about unisphere system",
		unit:     "",
		typeName: "gauge",
	},
}

func (pv *systemProvider) RunMeter() {
	// Set Meter
	logger.Info("Starting provider", "endpoint", pv.host.endpoint, "provider", pv.moduleName)
	//var err error
	meter := pv.meterProvider.Meter(pv.moduleName)

	// Create Observable Metrics
	var observableMap map[string]metric.Float64Observable
	observableMap = mappingMetricDescriptor(meter, systemMetricList)

	// Instruments...
	var observables []metric.Observable
	for _, obs := range observableMap {
		observables = append(observables, obs)
	}

	// Callback
	meter.RegisterCallback(func(ctx context.Context, observer metric.Observer) error {

		if pv.host.attributes == nil {
			return nil
		}

		clientAttrs := metric.WithAttributes(pv.host.attributes...)

		uc := pv.host.client
		//var data []byte

		// BasicSystemInfo
		data, err := uc.GetBasicSystemInfo([]string{"model", "softwareFullVersion"})
		if err != nil {
			return err
		}

		for _, entry := range data.Entries {
			content := entry.Content
			infoAttrs := metric.WithAttributes(
				attribute.String("product.name", content.Model),
				attribute.String("firmware.version", content.SoftwareFullVersion),
			)
			observer.ObserveFloat64(observableMap["info"], 1, clientAttrs, infoAttrs)
		}

		return nil
	}, observables...)

}
