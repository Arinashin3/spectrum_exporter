package provider

import (
	"context"
	"log/slog"
	"spectrum_exporter/config"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	sdkMetric "go.opentelemetry.io/otel/sdk/metric"
)

type enclosureProvider struct {
	moduleName    string
	interval      time.Duration
	meterProvider *sdkMetric.MeterProvider
	clientDesc    *ClientDesc
}

func init() {
	moduleName := "enclosure"
	registProvider(moduleName, &enclosureProvider{moduleName: moduleName})
}

func (pv *enclosureProvider) IsDefaultEnabled() bool {
	return true
}

func (pv *enclosureProvider) NewProvider(cfg *config.SpectrumConfig, moduleName string, cl *ClientDesc) Provider {
	pvConf := cfg.Providers.Enclosure
	enabled := pvConf.GetEnabled(pv.IsDefaultEnabled())
	interval := pvConf.GetInterval()

	if !enabled {
		return nil
	}
	if MetricExporter == nil {
		return nil
	}
	mp := NewMeterProvider(serviceName, interval, MetricExporter)
	return &enclosureProvider{
		moduleName:    moduleName,
		interval:      interval,
		meterProvider: mp,
		clientDesc:    cl,
	}
}

var EnclosuregMetricDescs = []*MetricDescriptor{
	{
		Key:      "status",
		Name:     "spectrum_enclosure_status",
		Desc:     "Information about the enclosure",
		Unit:     "",
		TypeName: "gauge",
	},
}

func (pv *enclosureProvider) Run(logger *slog.Logger) {
	logger.Info("Starting provider", "endpoint", pv.clientDesc.endpoint, "provider", pv.moduleName)
	meter := pv.meterProvider.Meter(pv.moduleName)

	// Register Metrics...
	var observableMap map[string]metric.Float64Observable
	observableMap = CreateMapMetricDescriptor(meter, EnclosuregMetricDescs, logger)

	// Register Metrics for Observables...
	var observableArray []metric.Observable
	for _, observable := range observableMap {
		observableArray = append(observableArray, observable)
	}

	// ==============================
	// Callback
	// ==============================
	meter.RegisterCallback(func(ctx context.Context, observer metric.Observer) error {
		// Client Attributes
		clientAttrs := metric.WithAttributes(pv.clientDesc.hostLabels...)

		// Request Data
		c := pv.clientDesc.client
		if !c.HealthCheck() {
			return nil
		}
		data, err := c.GetEnclosure()
		if err != nil {
			logger.Error("Failed to post", "err", err, "endpoint", pv.clientDesc.endpoint, "provider", pv.moduleName)
			return nil
		}
		if data == nil {
			logger.Warn("data is nil", "provider", pv.moduleName, "endpoint", pv.clientDesc.endpoint)
			return nil
		}

		for _, v := range data {
			additionalAttrs := metric.WithAttributes(
				attribute.String("enclosure.id", v.Id),
				attribute.String("enclosure.type", v.Type),
				attribute.String("serial.number", v.SerialNumber),
			)
			observer.ObserveFloat64(observableMap["status"], v.Status.Enum(), clientAttrs, additionalAttrs)
		}

		return nil
	}, observableArray...)

}
