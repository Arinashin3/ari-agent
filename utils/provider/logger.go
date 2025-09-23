package provider

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	otlplog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
)

func NewLoggerProvider(svName string, interval time.Duration, exp *otlplog.Exporter) *otlplog.LoggerProvider {
	return otlplog.NewLoggerProvider(
		otlplog.WithResource(resource.NewSchemaless(attribute.String("service.name", svName))),
		otlplog.WithProcessor(
			otlplog.NewBatchProcessor(*exp,
				otlplog.WithExportInterval(interval),
			),
		),
	)
}

func NewLogExporter(ctx context.Context, mode string, endpoint string, insecure bool) (*otlplog.Exporter, error) {
	var exp otlplog.Exporter
	var err error
	switch mode {
	case "http":
		if insecure {
			exp, err = otlploghttp.New(ctx,
				otlploghttp.WithEndpointURL(endpoint),
				otlploghttp.WithInsecure(),
			)
		} else {
			exp, err = otlploghttp.New(ctx,
				otlploghttp.WithEndpointURL(endpoint),
			)
		}
	case "grpc":
		if insecure {
			exp, err = otlploggrpc.New(ctx,
				otlploggrpc.WithEndpointURL(endpoint),
				otlploggrpc.WithInsecure(),
			)
		} else {
			exp, err = otlploggrpc.New(ctx,
				otlploggrpc.WithEndpointURL(endpoint),
			)
		}
	}
	return &exp, err
}
