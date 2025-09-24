package collector

import (
	"context"
	"reflect"
	"strings"
	"time"

	"github.com/Arinashin3/ari-agent/cmd/unisphere_exporter/config"
	"github.com/Arinashin3/ari-agent/utils/convert"
	"github.com/Arinashin3/ari-agent/utils/provider"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	otlpmetric "go.opentelemetry.io/otel/sdk/metric"
)

func init() {
	InitClientList()
	for _, cl := range ClientList {
		registMeterProvider(newLunProvider(cl))

	}
}

type lunProvider struct {
	moduleName    string
	interval      time.Duration
	meterProvider *otlpmetric.MeterProvider
	host          *ClientResourceStruct
}

func newLunProvider(cl *ClientResourceStruct) MeterProvider {
	moduleName := "lun"
	if MeterExporter == nil {
		return nil
	}

	pvConf := config.GetProviderLun()
	if !pvConf.Enabled {
		return nil
	}

	interval, err := time.ParseDuration(pvConf.Interval)
	if err != nil {
		logger.Error("Failed to parse interval", "provider", moduleName, "error", err)
	}

	mp := provider.NewMeterProvider(serviceName, interval, MeterExporter)

	return &lunProvider{
		moduleName:    moduleName,
		interval:      interval,
		meterProvider: mp,
		host:          cl,
	}
}

var lunMetricList = []*MetricDescriptor{
	{
		key:      "sizeTotal",
		name:     "unisphere_lun_total_size",
		desc:     "Total Size lun of unisphere",
		unit:     "mb",
		typeName: "gauge",
	},
	{
		key:      "sizeUsed",
		name:     "unisphere_lun_used_size",
		desc:     "Used Size lun of unisphere",
		unit:     "mb",
		typeName: "gauge",
	},
	{
		key:      "sizeAllocated",
		name:     "unisphere_lun_allocated_size",
		desc:     "Size of space actually allocated in the pool for the LUN.",
		unit:     "mb",
		typeName: "gauge",
	},
	{
		key:      "sizePreallocated",
		name:     "unisphere_lun_preallocated_size",
		desc:     "Total provisioned lun of unisphere lun",
		unit:     "mb",
		typeName: "gauge",
	},
}

func (pv *lunProvider) RunMeter() {
	// Set Meter
	logger.Info("Starting provider", "endpoint", pv.host.endpoint, "provider", pv.moduleName)
	meter := pv.meterProvider.Meter(pv.moduleName)
	uc := pv.host.client

	// Create Observable Metrics
	var observableMap map[string]metric.Float64Observable
	observableMap = mappingMetricDescriptor(meter, lunMetricList)

	// Observables...
	var observables []metric.Observable
	for _, obs := range observableMap {
		observables = append(observables, obs)
	}

	// Request Fields
	var paramsFields []string
	for _, v := range lunMetricList {
		paramsFields = append(paramsFields, v.key)
	}
	paramsFields = append(paramsFields, "name", "wwn")

	// Callback
	meter.RegisterCallback(func(ctx context.Context, observer metric.Observer) error {

		if pv.host.attributes == nil {
			return nil
		}

		clientAttrs := metric.WithAttributes(pv.host.attributes...)

		// Output Entry
		data, err := uc.GetLun(paramsFields)
		if err != nil {
			logger.Error("Failed to get lun", "error", err)
		}
		for _, entry := range data.Entries {
			content := entry.Content
			for _, field := range reflect.VisibleFields(reflect.TypeOf(content)) {
				key := strings.ToLower(field.Name)
				if observableMap[key] == nil {
					continue
				} else {
					fieldValue := reflect.ValueOf(content).FieldByName(field.Name).Interface()
					var f float64
					f = convert.InterfaceToFloat64(fieldValue)
					if err != nil {
						logger.Error("This type is not supported yet.", "client", pv.host.endpoint, "provider", pv.moduleName, "field", field.Name, "error", err)
					}
					observer.ObserveFloat64(observableMap[key], convert.BytesConvert(f, convert.Bytes, convert.Megabytes), clientAttrs, metric.WithAttributes(attribute.String("lun.name", content.Name), attribute.String("lun.wwn", content.Wwn)))
				}
			}
		}

		return nil
	}, observables...)

}
