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

type driveProvider struct {
	moduleName    string
	interval      time.Duration
	meterProvider *sdkMetric.MeterProvider
	clientDesc    *ClientDesc
}

func init() {
	moduleName := "drive"
	registProvider(moduleName, &driveProvider{moduleName: moduleName})
}

func (pv *driveProvider) IsDefaultEnabled() bool {
	return true
}

func (pv *driveProvider) NewProvider(cfg *config.SpectrumConfig, moduleName string, cl *ClientDesc) Provider {
	pvConf := cfg.Providers.Drive
	enabled := pvConf.GetEnabled(pv.IsDefaultEnabled())
	interval := pvConf.GetInterval()

	if !enabled {
		return nil
	}
	if MetricExporter == nil {
		return nil
	}
	mp := NewMeterProvider(serviceName, interval, MetricExporter)
	return &driveProvider{
		moduleName:    moduleName,
		interval:      interval,
		meterProvider: mp,
		clientDesc:    cl,
	}
}

var DriveMetricDescs = []*MetricDescriptor{
	{
		Key:      "status",
		Name:     "spectrum_drive_status",
		Desc:     "Information about the drive",
		Unit:     "",
		TypeName: "gauge",
	},
	{
		Key:      "capacity",
		Name:     "spectrum_drive_capacity",
		Desc:     "Information about the drive",
		Unit:     "mb",
		TypeName: "gauge",
	},
}

func (pv *driveProvider) Run(logger *slog.Logger) {
	logger.Info("Starting provider", "endpoint", pv.clientDesc.endpoint, "provider", pv.moduleName)
	meter := pv.meterProvider.Meter(pv.moduleName)

	// Register Metrics...
	var observableMap map[string]metric.Float64Observable
	observableMap = CreateMapMetricDescriptor(meter, DriveMetricDescs, logger)

	// Register Metrics for Observables...
	var observableDrive []metric.Observable
	for _, observable := range observableMap {
		observableDrive = append(observableDrive, observable)
	}

	// ==============================
	// Callback
	// ==============================
	meter.RegisterCallback(func(ctx context.Context, observer metric.Observer) error {
		// Client Attributes
		clientAttrs := metric.WithAttributes(pv.clientDesc.hostLabels...)

		// Request Data
		c := pv.clientDesc.client
		data, err := c.GetDrive()
		if err != nil {
			logger.Error("Failed to post drive info", "err", err)
			return nil
		}
		if data == nil {
			logger.Warn("data is nil", "provider", pv.moduleName, "endpoint", pv.clientDesc.endpoint)
			return nil
		}

		for _, v := range data {
			// Drive Attributes...
			additionalAttrs := metric.WithAttributes(
				attribute.String("slot.id", v.SlotId),
				attribute.String("tech.type", v.TechType),
				attribute.String("mdisk.name", v.MdiskName),
				attribute.String("use", v.Use),
			)
			observer.ObserveFloat64(observableMap["status"], v.Status.Enum(), clientAttrs, additionalAttrs)
			observer.ObserveFloat64(observableMap["capacity"], v.Capacity.Bytes().ToMiB(), clientAttrs, additionalAttrs)
		}

		return nil
	}, observableDrive...)

}
