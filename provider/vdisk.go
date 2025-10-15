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

type vdiskProvider struct {
	moduleName    string
	interval      time.Duration
	meterProvider *sdkMetric.MeterProvider
	clientDesc    *ClientDesc
}

func init() {
	moduleName := "vdisk"
	registProvider(moduleName, &vdiskProvider{moduleName: moduleName})
}

func (pv *vdiskProvider) IsDefaultEnabled() bool {
	return true
}

func (pv *vdiskProvider) NewProvider(cfg *config.SpectrumConfig, moduleName string, cl *ClientDesc) Provider {
	pvConf := cfg.Providers.Vdisk
	enabled := pvConf.GetEnabled(pv.IsDefaultEnabled())
	interval := pvConf.GetInterval()

	if !enabled {
		return nil
	}
	if MetricExporter == nil {
		return nil
	}
	mp := NewMeterProvider(serviceName, interval, MetricExporter)
	return &vdiskProvider{
		moduleName:    moduleName,
		interval:      interval,
		meterProvider: mp,
		clientDesc:    cl,
	}
}

var VDiskMetricDescs = []*MetricDescriptor{
	{
		Key:      "status",
		Name:     "spectrum_vdisk_status",
		Desc:     "Information about the vdisk",
		Unit:     "",
		TypeName: "gauge",
	},
	{
		Key:      "capacity",
		Name:     "spectrum_vdisk_capacity",
		Desc:     "Information about the vdisk",
		Unit:     "mb",
		TypeName: "gauge",
	},
}

func (pv *vdiskProvider) Run(logger *slog.Logger) {
	logger.Info("Starting provider", "endpoint", pv.clientDesc.endpoint, "provider", pv.moduleName)
	meter := pv.meterProvider.Meter(pv.moduleName)

	// Register Metrics...
	var observableMap map[string]metric.Float64Observable
	observableMap = CreateMapMetricDescriptor(meter, VDiskMetricDescs, logger)

	// Register Metrics for Observables...
	var observableVDisk []metric.Observable
	for _, observable := range observableMap {
		observableVDisk = append(observableVDisk, observable)
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
		data, err := c.GetVDisk()
		if err != nil {
			logger.Error("Failed to post", "err", err, "endpoint", pv.clientDesc.endpoint, "provider", pv.moduleName)
			return nil
		}
		if data == nil {
			logger.Warn("data is nil", "provider", pv.moduleName, "endpoint", pv.clientDesc.endpoint)
			return nil
		}

		for _, v := range data {
			// VDisk Attributes...
			additionalAttrs := metric.WithAttributes(
				attribute.String("vdisk.id", v.Id),
				attribute.String("vdisk.name", v.Name),
				attribute.String("vdisk.type", v.Type),
			)
			observer.ObserveFloat64(observableMap["status"], v.Status.Enum(), clientAttrs, additionalAttrs)
			observer.ObserveFloat64(observableMap["capacity"], v.Capacity.Bytes().ToMiB(), clientAttrs, additionalAttrs)
		}

		return nil
	}, observableVDisk...)

}
