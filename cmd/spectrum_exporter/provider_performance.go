package main

import (
	"context"
	"strconv"
	"time"

	"github.com/Arinashin3/ari-agent/utils/provider"
	"go.opentelemetry.io/otel/metric"
	sdkMetric "go.opentelemetry.io/otel/sdk/metric"
)

type systemStatsProvider struct {
	moduleName    string
	interval      time.Duration
	meterProvider *sdkMetric.MeterProvider
	clientDesc    *ClientDesc
}

func init() {
	moduleName := "performance"
	registProvider(moduleName, &systemStatsProvider{moduleName: moduleName})
}

func (pv *systemStatsProvider) IsDefaultEnabled() bool {
	return true
}

func (pv *systemStatsProvider) NewProvider(moduleName string, cl *ClientDesc) Provider {
	pvConf := cfg.Providers.Performance
	enabled := pvConf.GetEnabled(pv.IsDefaultEnabled())

	//enabled := pvConf.GetEnabled(pv.IsDefaultEnabled())
	interval := pvConf.GetInterval()

	if !enabled {
		return nil
	}
	if MetricExporter == nil {
		return nil
	}
	mp := provider.NewMeterProvider(serviceName, interval, MetricExporter)
	return &systemStatsProvider{
		moduleName:    moduleName,
		interval:      interval,
		meterProvider: mp,
		clientDesc:    cl,
	}
}

var lsSystemStatsMetricDescs = []*provider.MetricDescriptor{
	{
		Key:      "cpu_pc",
		Name:     "spectrum_system_stats_cpu",
		Desc:     "Information about the system",
		Unit:     "%",
		TypeName: "gauge",
	},
	{
		Key:      "compression_cpu_pc",
		Name:     "spectrum_system_stats_compression",
		Desc:     "Information about the system",
		Unit:     "%",
		TypeName: "gauge",
	},
	{
		Key:      "fc_mb",
		Name:     "spectrum_system_stats_fibrechannel",
		Desc:     "Information about the system",
		Unit:     "mbps",
		TypeName: "gauge",
	},
	{
		Key:      "fc_io",
		Name:     "spectrum_system_stats_fibrechannel",
		Desc:     "Information about the system",
		Unit:     "iops",
		TypeName: "gauge",
	},
	{
		Key:      "sas_mb",
		Name:     "spectrum_system_stats_sas",
		Desc:     "Information about the system",
		Unit:     "mbps",
		TypeName: "gauge",
	},
	{
		Key:      "sas_io",
		Name:     "spectrum_system_stats_sas",
		Desc:     "Information about the system",
		Unit:     "iops",
		TypeName: "gauge",
	},
	{
		Key:      "iscsi_mb",
		Name:     "spectrum_system_stats_iscsi",
		Desc:     "Information about the system",
		Unit:     "mbps",
		TypeName: "gauge",
	},
	{
		Key:      "iscsi_io",
		Name:     "spectrum_system_stats_iscsi",
		Desc:     "Information about the system",
		Unit:     "iops",
		TypeName: "gauge",
	},
	{
		Key:      "write_cache_pc",
		Name:     "spectrum_system_stats_write_cache",
		Desc:     "Information about the system",
		Unit:     "%",
		TypeName: "gauge",
	},
	{
		Key:      "total_cache_pc",
		Name:     "spectrum_system_stats_total_cache",
		Desc:     "Information about the system",
		Unit:     "%",
		TypeName: "gauge",
	},
	{
		Key:      "vdisk_r_mb",
		Name:     "spectrum_system_stats_vdisk_read",
		Desc:     "Information about the system",
		Unit:     "mbps",
		TypeName: "gauge",
	},
	{
		Key:      "vdisk_r_io",
		Name:     "spectrum_system_stats_vdisk_read",
		Desc:     "Information about the system",
		Unit:     "iops",
		TypeName: "gauge",
	},
	{
		Key:      "vdisk_r_ms",
		Name:     "spectrum_system_stats_vdisk_read",
		Desc:     "Information about the system",
		Unit:     "ms",
		TypeName: "gauge",
	},
	{
		Key:      "vdisk_w_mb",
		Name:     "spectrum_system_stats_vdisk_write",
		Desc:     "Information about the system",
		Unit:     "mbps",
		TypeName: "gauge",
	},
	{
		Key:      "vdisk_w_io",
		Name:     "spectrum_system_stats_vdisk_write",
		Desc:     "Information about the system",
		Unit:     "iops",
		TypeName: "gauge",
	},
	{
		Key:      "vdisk_w_ms",
		Name:     "spectrum_system_stats_vdisk_write",
		Desc:     "Information about the system",
		Unit:     "ms",
		TypeName: "gauge",
	},
	{
		Key:      "mdisk_r_mb",
		Name:     "spectrum_system_stats_mdisk_read",
		Desc:     "Information about the system",
		Unit:     "mbps",
		TypeName: "gauge",
	},
	{
		Key:      "mdisk_r_io",
		Name:     "spectrum_system_stats_mdisk_read",
		Desc:     "Information about the system",
		Unit:     "iops",
		TypeName: "gauge",
	},
	{
		Key:      "mdisk_r_ms",
		Name:     "spectrum_system_stats_mdisk_read",
		Desc:     "Information about the system",
		Unit:     "ms",
		TypeName: "gauge",
	},
	{
		Key:      "mdisk_w_mb",
		Name:     "spectrum_system_stats_mdisk_write",
		Desc:     "Information about the system",
		Unit:     "mbps",
		TypeName: "gauge",
	},
	{
		Key:      "mdisk_w_io",
		Name:     "spectrum_system_stats_mdisk_write",
		Desc:     "Information about the system",
		Unit:     "iops",
		TypeName: "gauge",
	},
	{
		Key:      "mdisk_w_ms",
		Name:     "spectrum_system_stats_mdisk_write",
		Desc:     "Information about the system",
		Unit:     "ms",
		TypeName: "gauge",
	},
	{
		Key:      "drive_r_mb",
		Name:     "spectrum_system_stats_drive_read",
		Desc:     "Information about the system",
		Unit:     "mbps",
		TypeName: "gauge",
	},
	{
		Key:      "drive_r_io",
		Name:     "spectrum_system_stats_drive_read",
		Desc:     "Information about the system",
		Unit:     "iops",
		TypeName: "gauge",
	},
	{
		Key:      "drive_r_ms",
		Name:     "spectrum_system_stats_drive_read",
		Desc:     "Information about the system",
		Unit:     "ms",
		TypeName: "gauge",
	},
	{
		Key:      "drive_w_mb",
		Name:     "spectrum_system_stats_drive_write",
		Desc:     "Information about the system",
		Unit:     "mbps",
		TypeName: "gauge",
	},
	{
		Key:      "drive_w_io",
		Name:     "spectrum_system_stats_drive_write",
		Desc:     "Information about the system",
		Unit:     "iops",
		TypeName: "gauge",
	},
	{
		Key:      "drive_w_ms",
		Name:     "spectrum_system_stats_drive_write",
		Desc:     "Information about the system",
		Unit:     "ms",
		TypeName: "gauge",
	},
}

func (pv *systemStatsProvider) Run() {
	logger.Info("Starting provider", "endpoint", pv.clientDesc.endpoint, "provider", pv.moduleName)
	meter := pv.meterProvider.Meter(pv.moduleName)

	// Register Metrics...
	var observableMap map[string]metric.Float64Observable
	observableMap = provider.CreateMapMetricDescriptor(meter, lsSystemStatsMetricDescs, logger)

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
		data := c.PostLsSystemStats()
		if data == nil {
			logger.Warn("data is nil", "provider", pv.moduleName, "endpoint", pv.clientDesc.endpoint)
			return nil
		}
		for _, v := range data {
			f, _ := strconv.ParseFloat(v.StatCurrent, 64)
			if observableMap[v.StatName] != nil {
				observer.ObserveFloat64(observableMap[v.StatName], f, clientAttrs)
			}
		}

		return nil
	}, observableArray...)

}
