package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"strconv"

	"github.com/Arinashin3/ari-agent/client/spectrum"
	"github.com/Arinashin3/ari-agent/config/cfgSpectrum"
	"github.com/Arinashin3/ari-agent/utils/provider"
	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/common/promslog"
	promslogflag "github.com/prometheus/common/promslog/flag"
	"go.opentelemetry.io/otel/attribute"
	sdkLog "go.opentelemetry.io/otel/sdk/log"
	sdkMetric "go.opentelemetry.io/otel/sdk/metric"
)

const serviceName = "spectrum_exporter"

var (
	configFile      = kingpin.Flag("config.file", "Path to config file.").Short('c').Default("config.yml").String()
	logger          *slog.Logger
	UsableProviders = make(map[string]Provider)
	Providers       []Provider
	MetricExporter  *sdkMetric.Exporter
	LogExporter     *sdkLog.Exporter
	cfg             *cfgSpectrum.SpectrumConfig
	isFailed        bool
)

type ClientDesc struct {
	endpoint     string
	customLabels []attribute.KeyValue
	hostLabels   []attribute.KeyValue
	client       *spectrum.Client
}

type Provider interface {
	NewProvider(moduleName string, clientDesc *ClientDesc) Provider
	Run()
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
	cfg = cfgSpectrum.NewSpectrumConfiguration()
	err := cfg.LoadFile(configFile)
	if err != nil {
		isFailed = true
		logger.Error("Failed to load config file.", "error", err)
	}

	// Define MetricExporter
	ctx := context.Background()
	endpoint := cfg.Server.Metrics.Endpoint + cfg.Server.Metrics.Api_Path
	mode := cfg.Server.Metrics.Mode
	insecure, _ := strconv.ParseBool(cfg.Server.Metrics.Insecure)
	if endpoint != "" {
		MetricExporter, err = provider.NewMetricExporter(ctx, mode, endpoint, insecure)
		if err != nil {
			isFailed = true
			logger.Error("Failed to create the Metric Exporter...", "error", err)
		}
	}

	// Define LogExporter
	endpoint = cfg.Server.Logs.Endpoint + cfg.Server.Logs.Api_Path
	mode = cfg.Server.Logs.Mode
	insecure, _ = strconv.ParseBool(cfg.Server.Logs.Insecure)
	if endpoint != "" {
		LogExporter, err = provider.NewLogExporter(ctx, mode, endpoint, insecure)
		if err != nil {
			isFailed = true
			logger.Error("Failed to create the Log Exporter...", "error", err)
		}
	}

	if isFailed {
		logger.Error("Failed to load configs...")
		os.Exit(1)
	}

	// Run Application
	RegistryProviders()
	RunProviders()

}

func RegistryProviders() {
	for _, clientConf := range cfg.Clients {
		var endpoint string
		var customLabels []attribute.KeyValue
		var insecure bool

		endpoint = clientConf.Endpoint
		for k, v := range clientConf.Labels {
			customLabels = append(customLabels, attribute.String(k, v))
		}
		username, password := cfg.SearchAuth(clientConf.Auth)
		if username == "" || password == "" {
			logger.Error("Cannot found the authentication credentials.", "auth", clientConf.Auth)
			os.Exit(1)
		}
		insecure, _ = strconv.ParseBool(clientConf.Insecure)
		cm := spectrum.NewClient(endpoint, username, password, insecure)

		cl := &ClientDesc{
			endpoint:     endpoint,
			customLabels: customLabels,
			hostLabels:   nil,
			client:       cm,
		}
		_ = UpdateAttributes(cl)

		for k, pv := range UsableProviders {
			tmp := pv.NewProvider(k, cl)
			if tmp != nil {
				Providers = append(Providers, tmp)
			}
		}
	}
}

func RunProviders() {
	for _, pv := range Providers {
		go pv.Run()
	}
	select {}
}

func UpdateAttributes(cl *ClientDesc) error {
	data := cl.client.PostLsSystem()
	if data == nil {
		return errors.New("Cannot post ls system")
	}

	var tmp []attribute.KeyValue
	tmp = append(cl.customLabels, attribute.String("instance", data.Id), attribute.String("host.name", data.Name))

	cl.hostLabels = tmp
	return nil
}
