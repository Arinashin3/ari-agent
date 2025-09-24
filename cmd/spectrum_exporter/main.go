package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/Arinashin3/ari-agent/client/spectrum"
	"github.com/Arinashin3/ari-agent/config"
	"github.com/Arinashin3/ari-agent/utils/provider"
	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/common/promslog"
	promslogflag "github.com/prometheus/common/promslog/flag"
	"go.opentelemetry.io/otel/attribute"
	sdkLog "go.opentelemetry.io/otel/sdk/log"
	sdkMetric "go.opentelemetry.io/otel/sdk/metric"
)

var (
	configFile      = kingpin.Flag("config.file", "Path to config file.").Short('c').Default("config.yml").String()
	logger          *slog.Logger
	UsableProviders = make(map[string]Provider)
	Providers       = make(map[string]Provider)
	MetricExporter  *sdkMetric.Exporter
	LogExporter     *sdkLog.Exporter
	configFailed    bool
)

type ClientDesc struct {
	endpoint     string
	customLabels []attribute.KeyValue
	hostLabels   []attribute.KeyValue
	client       *spectrum.Client
}

type Provider interface {
	NewProvider(moduleName string, interval time.Duration, clientDesc *ClientDesc) Provider
	IsDefaultEnabled() bool
}

func registProvider(moduleName string, pv Provider) error {
	UsableProviders[moduleName] = pv
	return nil
}

func main() {
	// Set Flag &nd Logger
	promslogConfig := &promslog.Config{}
	promslogflag.AddFlags(kingpin.CommandLine, promslogConfig)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	logger = promslog.New(promslogConfig)

	// Load Configuration Set Configuratios...
	logger.Info("Load Configs...")
	cfg := config.NewSpectrumConfiguration()
	err := cfg.LoadFile(configFile)
	if err != nil {
		configFailed = true
		logger.Error("Failed to load config file.", "error", err)
	}

	// Define MetricExporter
	ctx := context.Background()
	endpoint := cfg.GetMetricsEndpoint()
	mode := cfg.GetMetricsMode()
	insecure := cfg.GetMetricsInsecure()
	if endpoint != "" {
		MetricExporter, err = provider.NewMetricExporter(ctx, mode, endpoint, insecure)
		if err != nil {
			configFailed = true
			logger.Error("Failed to create the Metric Exporter...", "error", err)
		}
	}

	// Define LogExporter
	endpoint = cfg.GetLogsEndpoint()
	mode = cfg.GetLogsMode()
	insecure = cfg.GetLogsInsecure()
	if endpoint != "" {
		LogExporter, err = provider.NewLogExporter(ctx, mode, endpoint, insecure)
		if err != nil {
			configFailed = true
			logger.Error("Failed to create the Log Exporter...", "error", err)
		}
	}

	if configFailed {
		logger.Error("Failed to load configs...")
		os.Exit(1)
	}

	for moduleName, pv := range UsableProviders {
		pvConf := cfg.GetProvider(moduleName)
		if pv.IsDefaultEnabled() {
			logger.Info("moduleName", "name", moduleName, "config", pvConf)
		}

	}
	// Run Application
	RunApplication(cfg, logger)

}

func RunApplication(cfg config.Config, logger *slog.Logger) {

}
