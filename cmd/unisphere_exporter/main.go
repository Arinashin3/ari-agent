package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"strconv"

	"github.com/Arinashin3/ari-agent/config/cfgUnisphere"
	"github.com/Arinashin3/ari-agent/utils/provider"
	"github.com/Arinashin3/gounity"
	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/common/promslog"
	promslogflag "github.com/prometheus/common/promslog/flag"
	"go.opentelemetry.io/otel/attribute"
	sdkLog "go.opentelemetry.io/otel/sdk/log"
	sdkMetric "go.opentelemetry.io/otel/sdk/metric"
)

//TIP <p>To run your code, right-click the code and select <b>RunMeter</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>RunMeter</b> menu item from here.</p>

//var (
//	cfgFile = kingpin.Flag("config.file", "Path to config file.").Short('c').Default("config.yml").String()
//	//	listen  = kingpin.Flag("listen", "Address to listen on").Short('l').Default(":9748").String()
//)

const serviceName = "unisphere_exporter"

var (
	configFile      = kingpin.Flag("config.file", "Path to config file.").Short('c').Default("config.file").String()
	logger          *slog.Logger
	UsableProviders = make(map[string]Provider)
	Providers       []Provider
	MetricExporter  *sdkMetric.Exporter
	LogExporter     *sdkLog.Exporter
	cfg             *cfgUnisphere.UnisphereConfig
	isFailed        bool
)

type ClientDesc struct {
	endpoint     string
	customLabels []attribute.KeyValue
	hostLabels   []attribute.KeyValue
	client       *gounity.UnisphereClient
}

type Provider interface {
	NewProvider(moduleName string, desc *ClientDesc) Provider
	Run()
}

func registProvider(moduleName string, pv Provider) error {
	UsableProviders[moduleName] = pv
	return nil
}

func main() {
	// Set Flag & Logger
	promslogConfig := &promslog.Config{}
	promslogflag.AddFlags(kingpin.CommandLine, promslogConfig)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	logger = promslog.New(promslogConfig)

	// Load Configuration Set Configurations...
	logger.Info("Load Configs...")
	cfg = cfgUnisphere.NewUnisphereConfiguration()
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
			logger.Error("Failed to create metric exporter.", "error", err)
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

	RegistryProviders()
	// Check Failed
	if isFailed {
		logger.Error("Failed to load configs...")
		os.Exit(1)
	}

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
		if username != "" || password != "" {
			logger.Error("Cannot found the authentication credentials.", "auth", clientConf.Auth)
		}
		insecure, _ = strconv.ParseBool(clientConf.Insecure)
		cm := gounity.NewClient(endpoint, username, password, insecure)

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
	data, err := cl.client.GetSystemInstances([]string{"name", "serialNumber"}, nil)
	if data == nil {
		return errors.New("cannot to Update Attributes")
	}
	if err != nil {
		return err
	}

	var tmp []attribute.KeyValue
	if cl.customLabels != nil {
		tmp = cl.customLabels
	}
	for _, entry := range data.Entries {
		content := entry.Content
		tmp = append(tmp, attribute.String("host.name", content.Name))
		tmp = append(tmp, attribute.String("instance", content.SerialNumber))
	}
	cl.hostLabels = tmp

	return nil
}
