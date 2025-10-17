package provider

import (
	"context"
	"log/slog"
	"spectrum_exporter/config"
	"time"

	"go.opentelemetry.io/otel/metric"
	sdkMetric "go.opentelemetry.io/otel/sdk/metric"
)

type healthCheckProvider struct {
	moduleName    string
	interval      time.Duration
	meterProvider *sdkMetric.MeterProvider
	clientDesc    *ClientDesc
}

func init() {
	moduleName := "healthCheck"
	registProvider(moduleName, &healthCheckProvider{moduleName: moduleName})
}

func (pv *healthCheckProvider) IsDefaultEnabled() bool {
	return true
}

func (pv *healthCheckProvider) NewProvider(cfg *config.SpectrumConfig, moduleName string, cl *ClientDesc) Provider {
	interval := 1 * time.Minute
	if MetricExporter == nil {
		return nil
	}
	mp := NewMeterProvider(serviceName, interval, MetricExporter)
	return &healthCheckProvider{
		moduleName:    moduleName,
		interval:      interval,
		meterProvider: mp,
		clientDesc:    cl,
	}
}

var HealthCheckMetricDescs = []*MetricDescriptor{
	{
		Key:      "up",
		Name:     "spectrum_up",
		Desc:     "Information about the healthCheck",
		Unit:     "",
		TypeName: "gauge",
	},
}

func (pv *healthCheckProvider) Run(logger *slog.Logger) {
	c := pv.clientDesc.client
	logger.Debug("Starting provider", "endpoint", c.Endpoint(), "provider", pv.moduleName)
	meter := pv.meterProvider.Meter(pv.moduleName)

	// Register Metrics...
	var observableMap map[string]metric.Float64Observable
	observableMap = CreateMapMetricDescriptor(meter, HealthCheckMetricDescs, logger)

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
		check, err := c.Login()
		if err != nil {
			logger.Warn("Failed to login", "err", err, "endpoint", c.Endpoint(), "provider", pv.moduleName)
		}
		if check {
		}
		result := c.IsLogin()
		var f float64
		if result {
			_ = UpdateAttributes(pv.clientDesc)
			f = 1
		} else {
			f = 0
		}

		observer.ObserveFloat64(observableMap["up"], f, clientAttrs)

		return nil
	}, observableArray...)
}
