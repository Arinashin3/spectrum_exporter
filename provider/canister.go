package provider

import (
	"context"
	"log/slog"
	"spectrum_exporter/config"
	"spectrum_exporter/gospectrum/types"
	"time"

	"go.opentelemetry.io/otel/metric"
	sdkMetric "go.opentelemetry.io/otel/sdk/metric"
)

type enclosureCanisterProvider struct {
	moduleName    string
	interval      time.Duration
	meterProvider *sdkMetric.MeterProvider
	clientDesc    *ClientDesc
}

func init() {
	moduleName := "canister"
	registProvider(moduleName, &enclosureCanisterProvider{moduleName: moduleName})
}

func (pv *enclosureCanisterProvider) IsDefaultEnabled() bool {
	return true
}

func (pv *enclosureCanisterProvider) NewProvider(cfg *config.SpectrumConfig, moduleName string, cl *ClientDesc) Provider {
	pvConf := cfg.Providers.Canister
	enabled := pvConf.GetEnabled(pv.IsDefaultEnabled())
	interval := pvConf.GetInterval()

	if !enabled {
		return nil
	}
	if MetricExporter == nil {
		return nil
	}
	mp := NewMeterProvider(serviceName, interval, MetricExporter)
	return &enclosureCanisterProvider{
		moduleName:    moduleName,
		interval:      interval,
		meterProvider: mp,
		clientDesc:    cl,
	}
}

var EnclosureCanisterMetricDescs = []*MetricDescriptor{
	{
		Key:      "status",
		Name:     "spectrum_canister_status",
		Desc:     "Information about the enclosureCanister",
		Unit:     "",
		TypeName: "gauge",
	},
}

func (pv *enclosureCanisterProvider) Run(logger *slog.Logger) {
	logger.Info("Starting provider", "endpoint", pv.clientDesc.endpoint, "provider", pv.moduleName)
	meter := pv.meterProvider.Meter(pv.moduleName)

	// Register Metrics...
	var observableMap map[string]metric.Float64Observable
	observableMap = CreateMapMetricDescriptor(meter, EnclosureCanisterMetricDescs, logger)

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

		// Status Labels... (online, offline, degraded)
		var StatusCount = make(map[string]float64)
		for k, _ := range types.StatusMap {
			StatusCount[string(k)] = 0
		}

		// Request Data
		c := pv.clientDesc.client
		data, err := c.GetEnclosureCanister()
		if err != nil {
			logger.Error("Failed to post enclosureCanister info", "err", err)
			return nil
		}
		if data == nil {
			logger.Warn("data is nil", "provider", pv.moduleName, "endpoint", pv.clientDesc.endpoint)
			return nil
		}

		for _, v := range data {
			observer.ObserveFloat64(observableMap["status"], v.Status.Enum(), clientAttrs)
		}

		return nil
	}, observableArray...)

}
