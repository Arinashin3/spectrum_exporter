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

type mdiskProvider struct {
	moduleName    string
	interval      time.Duration
	meterProvider *sdkMetric.MeterProvider
	clientDesc    *ClientDesc
}

func init() {
	moduleName := "mdisk"
	registProvider(moduleName, &mdiskProvider{moduleName: moduleName})
}

func (pv *mdiskProvider) IsDefaultEnabled() bool {
	return true
}

func (pv *mdiskProvider) NewProvider(cfg *config.SpectrumConfig, moduleName string, cl *ClientDesc) Provider {
	pvConf := cfg.Providers.Mdisk
	enabled := pvConf.GetEnabled(pv.IsDefaultEnabled())
	interval := pvConf.GetInterval()

	if !enabled {
		return nil
	}
	if MetricExporter == nil {
		return nil
	}
	mp := NewMeterProvider(serviceName, interval, MetricExporter)
	return &mdiskProvider{
		moduleName:    moduleName,
		interval:      interval,
		meterProvider: mp,
		clientDesc:    cl,
	}
}

var MdiskMetricDescs = []*MetricDescriptor{
	{
		Key:      "status",
		Name:     "spectrum_mdisk_status",
		Desc:     "Information about the mdisk",
		Unit:     "",
		TypeName: "gauge",
	},
	{
		Key:      "capacity",
		Name:     "spectrum_mdisk_capacity",
		Desc:     "Information about the mdisk",
		Unit:     "mb",
		TypeName: "gauge",
	},
}

func (pv *mdiskProvider) Run(logger *slog.Logger) {
	logger.Info("Starting provider", "endpoint", pv.clientDesc.endpoint, "provider", pv.moduleName)
	meter := pv.meterProvider.Meter(pv.moduleName)

	// Register Metrics...
	var observableMap map[string]metric.Float64Observable
	observableMap = CreateMapMetricDescriptor(meter, MdiskMetricDescs, logger)

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
		data, err := c.GetMDisk()
		if err != nil {
			logger.Error("Failed to post", "err", err, "endpoint", pv.clientDesc.endpoint, "provider", pv.moduleName)
			return nil
		}
		if data == nil {
			logger.Warn("data is nil", "provider", pv.moduleName, "endpoint", pv.clientDesc.endpoint)
			return nil
		}

		for _, v := range data {
			// Mdisk Attributes...
			mdiskAttrs := metric.WithAttributes(
				attribute.String("mdisk.name", v.Name),
				attribute.String("group.name", v.MdiskGrpName),
				attribute.String("mode", string(v.Mode)),
			)
			observer.ObserveFloat64(observableMap["status"], v.Status.Enum(), clientAttrs, mdiskAttrs)
			observer.ObserveFloat64(observableMap["capacity"], v.Capacity.Bytes().ToMiB(), clientAttrs, mdiskAttrs)
		}

		return nil
	}, observableArray...)

}
