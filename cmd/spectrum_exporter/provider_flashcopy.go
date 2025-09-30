package main

import (
	"context"
	"time"

	"github.com/Arinashin3/ari-agent/utils/provider"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	sdkMetric "go.opentelemetry.io/otel/sdk/metric"
)

type flashcopyProvider struct {
	moduleName    string
	interval      time.Duration
	meterProvider *sdkMetric.MeterProvider
	clientDesc    *ClientDesc
}

func init() {
	moduleName := "flashcopy"
	registProvider(moduleName, &flashcopyProvider{moduleName: moduleName})
}

func (pv *flashcopyProvider) IsDefaultEnabled() bool {
	return true
}

func (pv *flashcopyProvider) NewProvider(moduleName string, cl *ClientDesc) Provider {
	pvConf := cfg.Providers.Flashcopy
	enabled := pvConf.GetEnabled(pv.IsDefaultEnabled())
	interval := pvConf.GetInterval()

	if !enabled {
		return nil
	}
	if MetricExporter == nil {
		return nil
	}
	mp := provider.NewMeterProvider(serviceName, interval, MetricExporter)
	return &flashcopyProvider{
		moduleName:    moduleName,
		interval:      interval,
		meterProvider: mp,
		clientDesc:    cl,
	}
}

var FlashcopyMetricDescs = []*provider.MetricDescriptor{
	{
		Key:      "progress",
		Name:     "spectrum_flashcopy_progress",
		Desc:     "Information about the flashcopy",
		Unit:     "%",
		TypeName: "gauge",
	},
	{
		Key:      "status",
		Name:     "spectrum_flashcopy_status",
		Desc:     "Information about the flashcopy",
		Unit:     "",
		TypeName: "gauge",
	},
	{
		Key:      "copy_rate",
		Name:     "spectrum_flashcopy_copy",
		Desc:     "Information about the flashcopy",
		Unit:     "rate",
		TypeName: "gauge",
	},
	{
		Key:      "clean_progress",
		Name:     "spectrum_flashcopy_clean_progress",
		Desc:     "Information about the flashcopy",
		Unit:     "%",
		TypeName: "gauge",
	},
	{
		Key:      "start_time",
		Name:     "spectrum_flashcopy_start",
		Desc:     "Information about the flashcopy",
		Unit:     "timestamp",
		TypeName: "gauge",
	},
}

func (pv *flashcopyProvider) Run() {
	logger.Info("Starting provider", "endpoint", pv.clientDesc.endpoint, "provider", pv.moduleName)
	meter := pv.meterProvider.Meter(pv.moduleName)

	// Register Metrics...
	var observableMap map[string]metric.Float64Observable
	observableMap = provider.CreateMapMetricDescriptor(meter, FlashcopyMetricDescs, logger)

	// Register Metrics for Observables...
	var observableArray []metric.Observable
	for _, observable := range observableMap {
		observableArray = append(observableArray, observable)
	}

	// ==============================
	// Callback
	// ==============================
	meter.RegisterCallback(func(ctx context.Context, observer metric.Observer) error {
		// Client Attributes
		clientAttrs := metric.WithAttributes(pv.clientDesc.hostLabels...)

		// Request Data
		c := pv.clientDesc.client
		data := c.PostLsFcMap()
		if data == nil {
			logger.Warn("data is nil", "provider", pv.moduleName, "endpoint", pv.clientDesc.endpoint)
			return nil
		}
		for _, v := range data {
			fcAttrs := metric.WithAttributes(
				attribute.String("source.vdisk", v.SourceVdiskName),
				attribute.String("target.vdisk", v.TargetVdiskName),
				attribute.String("fc.group", v.GroupName),
			)
			observer.ObserveFloat64(observableMap["status"], v.Status.Float64(), clientAttrs, fcAttrs)
			observer.ObserveFloat64(observableMap["progress"], v.Progress.Float64(), clientAttrs, fcAttrs)
			observer.ObserveFloat64(observableMap["copy_rate"], v.CopyRate.Float64(), clientAttrs, fcAttrs)
			observer.ObserveFloat64(observableMap["clean_progress"], v.CleanProgress.Float64(), clientAttrs, fcAttrs)
			observer.ObserveFloat64(observableMap["start_time"], v.StartTime.Time2Float64(), clientAttrs, fcAttrs)

		}

		// Info Attributes

		return nil
	}, observableArray...)

}
