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

type arrayProvider struct {
	moduleName    string
	interval      time.Duration
	meterProvider *sdkMetric.MeterProvider
	clientDesc    *ClientDesc
}

func init() {
	moduleName := "array"
	registProvider(moduleName, &arrayProvider{moduleName: moduleName})
}

func (pv *arrayProvider) IsDefaultEnabled() bool {
	return true
}

func (pv *arrayProvider) NewProvider(cfg *config.SpectrumConfig, moduleName string, cl *ClientDesc) Provider {
	pvConf := cfg.Providers.Array
	enabled := pvConf.GetEnabled(pv.IsDefaultEnabled())
	interval := pvConf.GetInterval()

	if !enabled {
		return nil
	}
	if MetricExporter == nil {
		return nil
	}
	mp := NewMeterProvider(serviceName, interval, MetricExporter)
	return &arrayProvider{
		moduleName:    moduleName,
		interval:      interval,
		meterProvider: mp,
		clientDesc:    cl,
	}
}

var ArrayMetricDescs = []*MetricDescriptor{
	{
		Key:      "status",
		Name:     "spectrum_array_status",
		Desc:     "Information about the array",
		Unit:     "",
		TypeName: "gauge",
	},
	{
		Key:      "capacity",
		Name:     "spectrum_array_capacity",
		Desc:     "Information about the array",
		Unit:     "mb",
		TypeName: "gauge",
	},
	{
		Key:      "raid_level",
		Name:     "spectrum_array_raid_level",
		Desc:     "Information about the array",
		Unit:     "",
		TypeName: "gauge",
	},
}

func (pv *arrayProvider) Run(logger *slog.Logger) {
	c := pv.clientDesc.client
	logger.Info("Starting provider", "endpoint", c.Endpoint(), "provider", pv.moduleName)
	meter := pv.meterProvider.Meter(pv.moduleName)

	// Register Metrics...
	var observableMap map[string]metric.Float64Observable
	observableMap = CreateMapMetricDescriptor(meter, ArrayMetricDescs, logger)

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
		if !c.IsLogin() {
			return nil
		}
		data, err := c.GetArraySummary()
		if err != nil {
			logger.Error("Failed to post", "err", err, "endpoint", c.Endpoint(), "provider", pv.moduleName)
			return nil
		}
		if data == nil {
			logger.Warn("data is nil", "provider", pv.moduleName, "endpoint", c.Endpoint())
			return nil
		}

		for _, v := range data {
			// Array Attributes...
			additionalAttrs := metric.WithAttributes(
				attribute.String("mdisk.name", v.MdiskName),
				attribute.String("group.name", v.MdiskGrpName),
			)
			observer.ObserveFloat64(observableMap["status"], v.Status.Enum(), clientAttrs, additionalAttrs)
			observer.ObserveFloat64(observableMap["capacity"], v.Capacity.Bytes().ToMiB(), clientAttrs, additionalAttrs)
			observer.ObserveFloat64(observableMap["raid_level"], v.RaidLevel.Enum(), clientAttrs, additionalAttrs)
		}

		return nil
	}, observableArray...)

}
