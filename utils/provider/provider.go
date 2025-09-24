package provider

import (
	"errors"
	"log/slog"

	"go.opentelemetry.io/otel/metric"
)

type MetricDescriptor struct {
	Key      string
	Name     string
	Desc     string
	Unit     string
	TypeName string
}

func CreateMapMetricDescriptor(meter metric.Meter, mds []*MetricDescriptor, logger *slog.Logger) map[string]metric.Float64Observable {
	mdmap := make(map[string]metric.Float64Observable)
	var err error
	for _, md := range mds {
		var tmp metric.Float64Observable
		desc := metric.WithDescription(md.Desc)
		unit := metric.WithUnit(md.Unit)
		switch md.TypeName {
		case "counter":
			tmp, err = meter.Float64ObservableCounter(md.Name, desc, unit)
		case "gauge":
			tmp, err = meter.Float64ObservableGauge(md.Name, desc, unit)
		default:
			err = errors.New("unknown metric type")
		}
		if err != nil {
			logger.Warn("cannot create metric", "error", err, "metric_key", md.Key, "metric_type", md.TypeName)
		}
		mdmap[md.Key] = tmp
	}
	return mdmap

}
