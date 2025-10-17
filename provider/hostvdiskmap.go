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

type hostvdiskmapProvider struct {
	moduleName    string
	interval      time.Duration
	meterProvider *sdkMetric.MeterProvider
	clientDesc    *ClientDesc
}

func init() {
	moduleName := "hostvdiskmap"
	registProvider(moduleName, &hostvdiskmapProvider{moduleName: moduleName})
}

func (pv *hostvdiskmapProvider) IsDefaultEnabled() bool {
	return true
}

func (pv *hostvdiskmapProvider) NewProvider(cfg *config.SpectrumConfig, moduleName string, cl *ClientDesc) Provider {
	pvConf := cfg.Providers.Hostvdiskmap
	enabled := pvConf.GetEnabled(pv.IsDefaultEnabled())
	interval := pvConf.GetInterval()

	if !enabled {
		return nil
	}
	if MetricExporter == nil {
		return nil
	}
	mp := NewMeterProvider(serviceName, interval, MetricExporter)
	return &hostvdiskmapProvider{
		moduleName:    moduleName,
		interval:      interval,
		meterProvider: mp,
		clientDesc:    cl,
	}
}

var HostVdiskMapMetricDescs = []*MetricDescriptor{
	{
		Key:      "info",
		Name:     "spectrum_hostvdiskmap_mapping_info",
		Desc:     "Information about the hostvdiskmap",
		Unit:     "",
		TypeName: "gauge",
	},
}

func (pv *hostvdiskmapProvider) Run(logger *slog.Logger) {
	c := pv.clientDesc.client
	logger.Info("Starting provider", "endpoint", c.Endpoint(), "provider", pv.moduleName)
	meter := pv.meterProvider.Meter(pv.moduleName)

	// Register Metrics...
	var observableMap map[string]metric.Float64Observable
	observableMap = CreateMapMetricDescriptor(meter, HostVdiskMapMetricDescs, logger)

	// Register Metrics for Observables...
	var observableHostVdiskMap []metric.Observable
	for _, observable := range observableMap {
		observableHostVdiskMap = append(observableHostVdiskMap, observable)
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
		data, err := c.GetHostVdiskMap()
		if err != nil {
			logger.Error("Failed to post", "err", err, "endpoint", c.Endpoint(), "provider", pv.moduleName)
			return nil
		}
		if data == nil {
			logger.Warn("data is nil", "provider", pv.moduleName, "endpoint", c.Endpoint())
			return nil
		}

		for _, v := range data {
			// HostVdiskMap Attributes...
			additionalAttrs := metric.WithAttributes(
				attribute.String("target.id", v.Id),
				attribute.String("mapping.type", v.MappingType),
				attribute.String("scsi.id", v.SCSIId),
				attribute.String("vdisk.id", v.VdiskId),
			)
			observer.ObserveFloat64(observableMap["info"], 1, clientAttrs, additionalAttrs)
		}

		return nil
	}, observableHostVdiskMap...)

}
