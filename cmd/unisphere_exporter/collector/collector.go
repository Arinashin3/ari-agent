package collector

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"strings"

	"github.com/Arinashin3/ari-agent/client/unisphere"
	"github.com/Arinashin3/ari-agent/cmd/unisphere_exporter/config"
	"github.com/Arinashin3/ari-agent/flag"
	"github.com/Arinashin3/ari-agent/utils/provider"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/log"
	sdkMetric "go.opentelemetry.io/otel/sdk/metric"
)

const serviceName = "unisphere_exporter"

var (
	ClientList      []*ClientResourceStruct
	meterProviders  []MeterProvider
	loggerProviders []LoggerProvider
	logger          *slog.Logger
	MeterExporter   *sdkMetric.Exporter
	LogExporter     *log.Exporter
	loadedConfig    bool
)

func InitClientList() {
	if logger == nil {
		logger = flag.Logger
	}
	if !loadedConfig {
		cfg := config.GetConfig()
		ctx := context.Background()
		for _, cl := range cfg.Clients {
			logger.Debug("Create the Otel Client", "endpoint", cl.Endpoint)

			var attrs []attribute.KeyValue
			for k, v := range cl.Labels {
				attrs = append(attrs, attribute.String(k, v))
			}
			insecure, err := strconv.ParseBool(cl.Insecure)
			if err != nil {
				logger.Error("Failed to parse insecure flag", "endpoint", cl.Endpoint, "err", err)
				return
			}
			auth := cfg.SearchAuth(cl.Auth)

			uc := unisphere.NewClient(cl.Endpoint, auth, insecure)
			ClientList = append(ClientList, &ClientResourceStruct{
				endpoint:    cl.Endpoint,
				customLabel: attrs,
				client:      uc,
			})
		}
		// Init MeterExporter
		if cfg.Server.Metrics.Enabled {
			mode := config.GetMetricsMode()
			endpoint := config.GetMetricsEndpoint()
			insecure := config.GetMetricsInsecure()

			MeterExporter, _ = provider.NewMetricExporter(ctx, mode, endpoint, insecure)
		}

		// Init LogExporter
		if cfg.Server.Logs.Enabled {
			mode := config.GetLogsMode()
			endpoint := config.GetLogsEndpoint()
			insecure := config.GetLogsInsecure()

			LogExporter, _ = provider.NewLogExporter(ctx, mode, endpoint, insecure)
		}
		loadedConfig = true
	}
}

type MeterProvider interface {
	RunMeter()
}
type LoggerProvider interface {
	RunLogger()
}

func Run() {
	for _, c := range ClientList {
		c.UpdateClientAttributes()
	}
	for _, mp := range meterProviders {
		go mp.RunMeter()
	}
	for _, lp := range loggerProviders {
		go lp.RunLogger()
	}
}

// // Real Start
func (c *ClientResourceStruct) UpdateClientAttributes() error {
	uc := c.client

	// System (Get Hostname & SerialNumber)
	// Update Client's Labels
	data, err := uc.GetSystem([]string{"name", "serialNumber"})
	if err != nil {
		return err
	}

	var tmpAttrs []attribute.KeyValue
	if c.customLabel != nil {
		tmpAttrs = c.customLabel
	}
	for _, entries := range data.Entries {
		content := entries.Content
		tmpAttrs = append(tmpAttrs, attribute.String("host.name", content.Name))
		tmpAttrs = append(tmpAttrs, attribute.String("instance", content.SerialNumber))
	}
	c.attributes = tmpAttrs

	return nil
}

func registMeterProvider(provider MeterProvider) {
	if provider != nil {
		meterProviders = append(meterProviders, provider)
	}
}

func registLoggerProvider(provider LoggerProvider) {
	if provider != nil {
		loggerProviders = append(loggerProviders, provider)
	}
}

// Map Key will be to Lower
func mappingMetricDescriptor(meter metric.Meter, mds []*MetricDescriptor) map[string]metric.Float64Observable {
	mdmap := make(map[string]metric.Float64Observable)
	var err error
	for _, md := range mds {
		var tmp metric.Float64Observable
		desc := metric.WithDescription(md.desc)
		unit := metric.WithUnit(md.unit)
		switch md.typeName {
		case "counter":
			tmp, err = meter.Float64ObservableCounter(md.name, desc, unit)
		case "gauge":
			tmp, err = meter.Float64ObservableGauge(md.name, desc, unit)
		default:
			err = errors.New("unknown metric type")
		}
		if err != nil {
			logger.Error("cannot create metrics", "err", err)
		}
		mdmap[strings.ToLower(md.key)] = tmp
	}
	return mdmap

}
