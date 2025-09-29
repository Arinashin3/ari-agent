package main

import (
	"context"
	"time"

	"github.com/Arinashin3/ari-agent/utils/provider"
	"go.opentelemetry.io/otel/log"
	sdkLog "go.opentelemetry.io/otel/sdk/log"
)

func init() {
	moduleName := "event"
	registProvider(moduleName, &eventProvider{moduleName: moduleName})
}

func (pv *eventProvider) IsDefaultEnabled() bool {
	return true
}

func (pv *eventProvider) NewProvider(moduleName string, cl *ClientDesc) Provider {
	pvConf := cfg.Providers.Event
	enabled := pvConf.GetEnabled(pv.IsDefaultEnabled())
	interval := pvConf.GetInterval()

	if !enabled {
		return nil
	}
	if LogExporter == nil {
		return nil
	}
	lp := provider.NewLoggerProvider(serviceName, interval, LogExporter)
	return &eventProvider{
		moduleName:     moduleName,
		interval:       interval,
		level:          pvConf.Level,
		loggerProvider: lp,
		clientDesc:     cl,
	}
}

type eventProvider struct {
	moduleName     string
	interval       time.Duration
	level          int
	loggerProvider *sdkLog.LoggerProvider
	clientDesc     *ClientDesc
}

func (pv *eventProvider) Run() {
	logger.Info("Starting provider", "endpoint", pv.clientDesc.endpoint, "provider", pv.moduleName)
	ctx := context.Background()
	ctime := time.Now().Add(-time.Hour).UTC()
	uc := pv.clientDesc.client
	lp := pv.loggerProvider

	for {

		pvlogger := lp.Logger(pv.moduleName, log.WithInstrumentationAttributes(pv.clientDesc.hostLabels...))
		var fields = []string{
			"creationTime",
			"severity",
			"messageId",
			"message",
			"source",
		}
		data, err := uc.GetEvent(fields, ctime)
		if err != nil {
			logger.Error("Error to GET EventLog", "err", err)
			time.Sleep(pv.interval)
			continue
		}
		if data == nil {
			time.Sleep(pv.interval)
			continue
		}

		for _, entry := range data.Entries {
			record := log.Record{}
			content := entry.Content
			var severity string
			if pv.level > int(content.Severity) {
				continue
			}
			switch content.Severity {
			case 8:
				severity = "OK"
			case 7:
				severity = "DEBUG"
			case 6:
				severity = "INFO"
			case 5:
				severity = "NOTICE"
			case 4:
				severity = "WARNING"
			case 3:
				severity = "ERROR"
			case 2:
				severity = "CRITICAL"
			case 1:
				severity = "ALERT"
			case 0:
				severity = "EMERGENCY"

			}

			record.SetTimestamp(content.CreationTime)
			record.SetObservedTimestamp(content.CreationTime)
			record.SetBody(log.StringValue(content.Message))
			record.AddAttributes(
				log.String("level", severity),
				log.String("message.id", content.MessageId),
				log.String("source", content.Source),
			)
			pvlogger.Emit(ctx, record)

		}
		ctime = data.Updated.UTC()

		time.Sleep(pv.interval)
	}

}
