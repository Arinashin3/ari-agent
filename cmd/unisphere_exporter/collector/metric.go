package collector

import (
	"context"
	"reflect"
	"strconv"
	"strings"
	"time"

	config2 "github.com/Arinashin3/ari-agent/cmd/unisphere_exporter/config"
	"github.com/Arinashin3/ari-agent/utils/provider"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	otlpmetric "go.opentelemetry.io/otel/sdk/metric"
)

func init() {
	InitClientList()
	for _, cl := range ClientList {
		registMeterProvider(newMetricRealTimeQueryProvider(cl, "metric_a"))
		registMeterProvider(newMetricRealTimeQueryProvider(cl, "metric_b"))
		registMeterProvider(newMetricRealTimeQueryProvider(cl, "metric_c"))
	}
}

type MetricRealTimeQueryProvider struct {
	moduleName    string
	interval      time.Duration
	queryId       string
	paths         []string
	meterProvider *otlpmetric.MeterProvider
	host          *ClientResourceStruct
}

func newMetricRealTimeQueryProvider(cl *ClientResourceStruct, moduleName string) MeterProvider {
	var m *config2.ProviderMetric

	if MeterExporter == nil {
		return nil
	}

	switch moduleName {
	case "metric_a":
		m = config2.GetProviderMetricA()
	case "metric_b":
		m = config2.GetProviderMetricB()
	case "metric_c":
		m = config2.GetProviderMetricC()

	}
	if !m.Enabled {
		return nil
	}
	if len(m.Paths) == 0 {
		return nil
	}
	interval, err := time.ParseDuration(m.Interval)
	if err != nil {
		logger.Error("Failed to parse interval", "provider", moduleName, "error", err)
		return nil
	}

	mp := provider.NewMeterProvider(serviceName, interval, MeterExporter)

	return &MetricRealTimeQueryProvider{
		moduleName:    moduleName,
		interval:      interval,
		paths:         m.Paths,
		meterProvider: mp,
		host:          cl,
	}
}

func (pv *MetricRealTimeQueryProvider) RunMeter() {
	// Set Meter
	logger.Info("Starting provider", "endpoint", pv.host.endpoint, "provider", pv.moduleName)
	meter := pv.meterProvider.Meter(pv.moduleName)

	// Get Metric Descriptions...
	var err error
	uc := pv.host.client

	metricData, err := uc.GetMetric([]string{"name", "path", "type", "unitDisplayString", "description"}, "realtime")
	if err != nil {
		logger.Error("Failed to get metric instances", "provider", pv.moduleName, "error", err)
		return
	}

	paths := pv.paths
	var metricDescList []*MetricDescriptor
	var metricPaths []string
	for _, entry := range metricData.Entries {
		content := entry.Content
		for _, path := range paths {
			var match bool
			if string(path[len(path)-1]) == "%" {
				pattern := strings.Replace(path, "%", "", -1)
				match = strings.Contains(content.Path, pattern)
			} else {
				if path == content.Path {
					match = true
				}
			}
			if match {
				var mType string
				switch content.Type {
				case 2:
					mType = "counter"
				case 3:
					mType = "counter"
				case 4:
					mType = "gauge"
				case 5:
					mType = "gauge"
				case 6:
					logger.Info("SKIP METRIC: this metric's value is not number", "provider", pv.moduleName, "path", content.Path)
					continue
				case 7:
					mType = "counter"
				case 8:
					mType = "counter"
				}
				tmp := "unisphere_" + strings.Replace(strings.ToLower(content.Path), ".*.", "_", -1)

				metricPaths = append(metricPaths, content.Path)
				metricDescList = append(metricDescList, &MetricDescriptor{
					key:      content.Path,
					name:     strings.Replace(tmp, ".", "_", -1),
					desc:     content.Description,
					unit:     strings.ToLower(content.UnitDisplayString),
					typeName: mType,
				})
			}
		}
	}

	var observableMap map[string]metric.Float64Observable
	observableMap = mappingMetricDescriptor(meter, metricDescList)

	var observables []metric.Observable
	for _, obs := range observableMap {
		observables = append(observables, obs)
	}

	if len(metricPaths) > 48 {
		logger.Error("Too Many Paths", "provider", pv.moduleName, "path_count", len(metricPaths))

	}
	logger.Info("Create Metric Query", "provider", pv.moduleName, "path_count", len(metricPaths))

	pv.queryId, err = uc.PostMetricRealTimeQuery(metricPaths, pv.interval)
	if err != nil {
		logger.Error("Failed to post metric query", "provider", pv.moduleName, "error", err)
	}

	// Callback
	meter.RegisterCallback(func(ctx context.Context, observer metric.Observer) error {
		if pv.queryId == "" {
			pv.queryId, err = uc.PostMetricRealTimeQuery(metricPaths, pv.interval)
			if err != nil {
				logger.Warn("Failed to post metric", "provider", pv.moduleName, "error", err)
				return nil
			}
		}

		if pv.host.attributes == nil {
			return nil
		}

		clientAttrs := metric.WithAttributes(pv.host.attributes...)

		// Metric (MetricRealTimeQuery Performance)
		data, err := uc.GetMetricQueryResult(pv.queryId)
		if err != nil {
			logger.Error("Failed to get metric query", "provider", pv.moduleName, "error", err)
			return nil
		}

		for _, entry := range data.Entries {
			content := entry.Content

			// Create Label Name
			var labels []string
			var preString string
			for _, v := range strings.Split(content.Path, ".") {
				if v == "*" {
					labels = append(labels, preString)
				}
				preString = v
			}

			// Get Metric
			for k1, v1 := range content.Values.(map[string]interface{}) {
				if reflect.TypeOf(v1).Kind().String() != "map" {
					var f float64
					f, err = strconv.ParseFloat(v1.(string), 64)
					if err != nil {
						logger.Error("Failed to parse metric value", "provider", pv.moduleName, "path", content.Path, "error", err)
					}
					observer.ObserveFloat64(observableMap[strings.ToLower(content.Path)], f, clientAttrs, metric.WithAttributes(attribute.String(labels[0], k1)))
					continue
				}
				for k2, v2 := range v1.(map[string]interface{}) {
					if reflect.TypeOf(v2).Kind().String() != "map" {
						var f float64
						f, err = strconv.ParseFloat(v2.(string), 64)
						if err != nil {
							logger.Error("Failed to parse metric value", "provider", pv.moduleName, "path", content.Path, "error", err)
						}
						observer.ObserveFloat64(observableMap[strings.ToLower(content.Path)], f, clientAttrs, metric.WithAttributes(attribute.String(labels[0], k1), attribute.String(labels[1], k2)))
						continue
					}
					for k3, v3 := range v2.(map[string]interface{}) {
						if reflect.TypeOf(v3).Kind().String() != "map" {
							var f float64
							f, err = strconv.ParseFloat(v3.(string), 64)
							if err != nil {
								logger.Error("Failed to parse metric value", "provider", pv.moduleName, "path", content.Path, "error", err)
							}
							observer.ObserveFloat64(observableMap[strings.ToLower(content.Path)], f, clientAttrs, metric.WithAttributes(attribute.String(labels[0], k1), attribute.String(labels[1], k2), attribute.String(labels[2], k3)))
							continue
						}
					}
				}

			}
		}

		return nil
	}, observables...)
}
