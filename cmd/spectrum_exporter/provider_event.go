package main

import (
	"time"

	"context"

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
	//
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
	ctime := time.Now().Add(-1 * time.Hour).UTC()
	cl := pv.clientDesc.client
	lp := pv.loggerProvider

	for {
		pvlogger := lp.Logger(pv.moduleName, log.WithInstrumentationAttributes(pv.clientDesc.hostLabels...))

		data := cl.PostLsEventLog(ctime)

		if data == nil {
			time.Sleep(pv.interval)
			continue
		}

		for _, event := range data {
			record := log.Record{}
			eventTime, err := time.ParseInLocation("060102150405", event.LastTimestamp, time.Local)
			if err != nil {
				logger.Error("Error parsing timestamp", "err", err)
			}
			record.SetObservedTimestamp(eventTime)
			if event.ErrorCode != "" {
				record.AddAttributes(
					log.String("level", "ALERT"),
				)
			} else {
				record.AddAttributes(
					log.String("level", "INFO"),
				)
			}
			record.AddAttributes(
				log.String("error.code", event.ErrorCode),
				log.String("message.id", event.EventId),
				log.String("object.name", event.ObjectName),
				log.String("status", event.Status),
			)
			record.SetBody(log.StringValue(event.Description))

			pvlogger.Emit(ctx, record)
		}
		ctime = time.Now().Local()
		time.Sleep(pv.interval)
	}
}
