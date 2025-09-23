package collector

import (
	"context"
	"time"

	"github.com/Arinashin3/ari-agent/cmd/unisphere_exporter/config"
	"github.com/Arinashin3/ari-agent/utils/convert"
	"github.com/Arinashin3/ari-agent/utils/provider"

	"go.opentelemetry.io/otel/metric"
	otlpmetric "go.opentelemetry.io/otel/sdk/metric"
)

func init() {
	InitClientList()
	for _, cl := range ClientList {
		registMeterProvider(newCapacityProvider(cl))

	}
}

type capacityProvider struct {
	moduleName    string
	interval      time.Duration
	meterProvider *otlpmetric.MeterProvider
	host          *ClientResourceStruct
}

func newCapacityProvider(cl *ClientResourceStruct) MeterProvider {
	moduleName := "capacity"
	if MeterExporter == nil {
		return nil
	}

	pvConf := config.GetProviderCapacity()
	if !pvConf.Enabled {
		logger.Debug("Disabled provider", "client", cl.endpoint, "provider", moduleName)
		return nil
	}
	logger.Info("Enabled provider", "client", cl.endpoint, "provider", moduleName)

	interval, err := time.ParseDuration(pvConf.Interval)
	if err != nil {
		logger.Error("Failed to parse interval", "provider", moduleName, "error", err)
	}

	mp := provider.NewMeterProvider(serviceName, interval, MeterExporter)

	return &capacityProvider{
		moduleName:    moduleName,
		interval:      interval,
		meterProvider: mp,
		host:          cl,
	}
}

var capacityMetricList = []*MetricDescriptor{
	{
		key:      "total_cap",
		name:     "unisphere_capacity_total_capacity",
		desc:     "Total capacity of unisphere capacity",
		unit:     "mb",
		typeName: "gauge",
	},
	{
		key:      "used_cap",
		name:     "unisphere_capacity_used_capacity",
		desc:     "Used capacity of unisphere capacity",
		unit:     "mb",
		typeName: "gauge",
	},
	{
		key:      "free_cap",
		name:     "unisphere_capacity_free_capacity",
		desc:     "Free capacity of unisphere capacity",
		unit:     "mb",
		typeName: "gauge",
	},
	{
		key:      "pre_cap",
		name:     "unisphere_capacity_preallocated_capacity",
		desc:     "Total provisioned capacity of unisphere capacity",
		unit:     "mb",
		typeName: "gauge",
	},
	{
		key:      "total_pv",
		name:     "unisphere_capacity_total_provision",
		desc:     "Total provisioned capacity of unisphere capacity",
		unit:     "mb",
		typeName: "gauge",
	},
	{
		key:      "used_pv",
		name:     "unisphere_capacity_used_provision",
		desc:     "Used provisioned capacity of unisphere capacity",
		unit:     "mb",
		typeName: "gauge",
	},
	{
		key:      "free_pv",
		name:     "unisphere_capacity_free_provision",
		desc:     "Free provisioned capacity of unisphere capacity",
		unit:     "mb",
		typeName: "gauge",
	},
}

func (pv *capacityProvider) RunMeter() {
	// Set Meter
	logger.Info("Starting provider", "endpoint", pv.host.endpoint, "provider", pv.moduleName)
	meter := pv.meterProvider.Meter(pv.moduleName)

	// Create Observable Metrics
	var observableMap map[string]metric.Float64Observable
	observableMap = mappingMetricDescriptor(meter, capacityMetricList)

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
		// systemCapacity
		var fields = []string{
			"sizeTotal",
			"sizeUsed",
			"sizeFree",
			"totalLogicalSize",
			"sizePreallocated",
		}
		data, err := uc.GetSystemCapacity(fields)
		if err != nil {
			logger.Error("Failed to get system capacity", "error", err)
		}
		for _, entry := range data.Entries {
			content := entry.Content

			// Capacity

			observer.ObserveFloat64(observableMap["total_cap"], convert.BytesConvert(float64(content.SizeTotal), convert.Bytes, convert.Megabytes), clientAttrs)
			observer.ObserveFloat64(observableMap["used_cap"], convert.BytesConvert(float64(content.SizeUsed), convert.Bytes, convert.Megabytes), clientAttrs)
			observer.ObserveFloat64(observableMap["free_cap"], convert.BytesConvert(float64(content.SizeFree), convert.Bytes, convert.Megabytes), clientAttrs)
			observer.ObserveFloat64(observableMap["pre_cap"], convert.BytesConvert(float64(content.SizePreallocated), convert.Bytes, convert.Megabytes), clientAttrs)

			// Provisioned Capacity
			contentFree := content.TotalLogicalSize - content.SizeUsed
			observer.ObserveFloat64(observableMap["total_pv"], convert.BytesConvert(float64(content.TotalLogicalSize), convert.Bytes, convert.Megabytes), clientAttrs)
			observer.ObserveFloat64(observableMap["used_pv"], convert.BytesConvert(float64(content.SizeUsed), convert.Bytes, convert.Megabytes), clientAttrs)
			observer.ObserveFloat64(observableMap["free_pv"], convert.BytesConvert(float64(contentFree), convert.Bytes, convert.Megabytes), clientAttrs)
		}

		return nil
	}, observables...)

}
