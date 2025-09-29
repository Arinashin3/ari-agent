package provider

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	sdkMetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

func NewMeterProvider(svName string, interval time.Duration, exp *sdkMetric.Exporter) *sdkMetric.MeterProvider {
	return sdkMetric.NewMeterProvider(
		sdkMetric.WithResource(resource.NewSchemaless(attribute.String("service.name", svName))),
		sdkMetric.WithReader(
			sdkMetric.NewPeriodicReader(*exp,
				sdkMetric.WithInterval(interval),
			),
		),
	)
}

func NewMetricExporter(ctx context.Context, mode string, endpoint string, insecure bool) (*sdkMetric.Exporter, error) {
	var exp sdkMetric.Exporter
	var err error
	switch mode {
	case "http":
		if insecure {
			exp, err = otlpmetrichttp.New(ctx,
				otlpmetrichttp.WithEndpointURL(endpoint),
				otlpmetrichttp.WithInsecure(),
			)
		} else {
			exp, err = otlpmetrichttp.New(ctx,
				otlpmetrichttp.WithEndpointURL(endpoint),
			)
		}
	case "grpc":
		if insecure {
			exp, err = otlpmetricgrpc.New(ctx,
				otlpmetricgrpc.WithEndpointURL(endpoint),
				otlpmetricgrpc.WithInsecure(),
			)
		} else {
			exp, err = otlpmetricgrpc.New(ctx,
				otlpmetricgrpc.WithEndpointURL(endpoint),
			)
		}
	}
	return &exp, err
}
