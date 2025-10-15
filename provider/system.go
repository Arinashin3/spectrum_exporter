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

type systemProvider struct {
	moduleName    string
	interval      time.Duration
	meterProvider *sdkMetric.MeterProvider
	clientDesc    *ClientDesc
}

func init() {
	moduleName := "system"
	registProvider(moduleName, &systemProvider{moduleName: moduleName})
}

func (pv *systemProvider) IsDefaultEnabled() bool {
	return true
}

func (pv *systemProvider) NewProvider(cfg *config.SpectrumConfig, moduleName string, cl *ClientDesc) Provider {
	pvConf := cfg.Providers.System
	enabled := pvConf.GetEnabled(pv.IsDefaultEnabled())
	interval := pvConf.GetInterval()

	if !enabled {
		return nil
	}
	if MetricExporter == nil {
		return nil
	}
	mp := NewMeterProvider(serviceName, interval, MetricExporter)
	return &systemProvider{
		moduleName:    moduleName,
		interval:      interval,
		meterProvider: mp,
		clientDesc:    cl,
	}
}

var SystemMetricDescs = []*MetricDescriptor{
	{
		Key:      "info",
		Name:     "spectrum_system_info",
		Desc:     "Information about the system",
		Unit:     "",
		TypeName: "gauge",
	},
	{
		Key:      "TotalVdiskCapacity",
		Name:     "spectrum_system_total_vdisk_capacity",
		Desc:     "Information about the system",
		Unit:     "mb",
		TypeName: "gauge",
	},
	{
		Key:      "TotalMdiskCapacity",
		Name:     "spectrum_system_total_mdisk_capacity",
		Desc:     "Information about the system",
		Unit:     "mb",
		TypeName: "gauge",
	},
	{
		Key:      "TotalUsedCapacity",
		Name:     "spectrum_system_total_used_capacity",
		Desc:     "Information about the system",
		Unit:     "mb",
		TypeName: "gauge",
	},
	{
		Key:      "TotalFreeSpace",
		Name:     "spectrum_system_total_free_space",
		Desc:     "Information about the system",
		Unit:     "mb",
		TypeName: "gauge",
	},
	{
		Key:      "SpaceAllocatedToVdisks",
		Name:     "spectrum_system_allocated_to_vdisks",
		Desc:     "Information about the system",
		Unit:     "mb",
		TypeName: "gauge",
	},
}

func (pv *systemProvider) Run(logger *slog.Logger) {
	logger.Info("Starting provider", "endpoint", pv.clientDesc.endpoint, "provider", pv.moduleName)
	meter := pv.meterProvider.Meter(pv.moduleName)

	// Register Metrics...
	var observableMap map[string]metric.Float64Observable
	observableMap = CreateMapMetricDescriptor(meter, SystemMetricDescs, logger)

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
		data, err := c.GetSystem()
		if err != nil {
			logger.Error("Failed to post", "err", err, "endpoint", pv.clientDesc.endpoint, "provider", pv.moduleName)
			return nil
		}
		if data == nil {
			logger.Warn("data is nil", "provider", pv.moduleName, "endpoint", pv.clientDesc.endpoint)
			return nil
		}

		// Info Attributes
		infoAttrs := metric.WithAttributes(
			attribute.String("product.name", data.ProductName),
			attribute.String("firmware.version", data.CodeLevel),
			attribute.String("ip.address", data.ConsoleIP),
		)
		observer.ObserveFloat64(observableMap["info"], 1, clientAttrs, infoAttrs)
		observer.ObserveFloat64(observableMap["TotalVdiskCapacity"], data.TotalVdiskCapacity.Bytes().ToMiB(), clientAttrs)
		observer.ObserveFloat64(observableMap["TotalMdiskCapacity"], data.TotalMdiskCapacity.Bytes().ToMiB(), clientAttrs)
		observer.ObserveFloat64(observableMap["TotalUsedCapacity"], data.TotalUsedCapacity.Bytes().ToMiB(), clientAttrs)
		observer.ObserveFloat64(observableMap["TotalFreeSpace"], data.TotalFreeSpace.Bytes().ToMiB(), clientAttrs)
		observer.ObserveFloat64(observableMap["SpaceAllocatedToVdisks"], data.SpaceAllocatedToVdisks.Bytes().ToMiB(), clientAttrs)

		return nil
	}, observableArray...)

}
