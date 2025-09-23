package collector

import (
	"github.com/Arinashin3/ari-agent/client/unisphere"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	otlplogger "go.opentelemetry.io/otel/sdk/log"
	otlpmetric "go.opentelemetry.io/otel/sdk/metric"
)

type ClientResourceStruct struct {
	endpoint    string
	customLabel []attribute.KeyValue
	attributes  []attribute.KeyValue
	client      *unisphere.Client
}

type ProviderSet struct {
	name     string
	queryId  string
	interval time.Duration
	mp       *otlpmetric.MeterProvider
	lp       *otlplogger.LoggerProvider
	observer []*metric.Observable
}

type MetricDescriptor struct {
	key      string
	name     string
	desc     string
	unit     string
	typeName string
}
