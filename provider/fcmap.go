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

type fcmapProvider struct {
	moduleName    string
	interval      time.Duration
	meterProvider *sdkMetric.MeterProvider
	clientDesc    *ClientDesc
}

func init() {
	moduleName := "flashcopy"
	registProvider(moduleName, &fcmapProvider{moduleName: moduleName})
}

func (pv *fcmapProvider) IsDefaultEnabled() bool {
	return true
}

func (pv *fcmapProvider) NewProvider(cfg *config.SpectrumConfig, moduleName string, cl *ClientDesc) Provider {
	pvConf := cfg.Providers.Flashcopy
	enabled := pvConf.GetEnabled(pv.IsDefaultEnabled())
	interval := pvConf.GetInterval()

	if !enabled {
		return nil
	}
	if MetricExporter == nil {
		return nil
	}
	mp := NewMeterProvider(serviceName, interval, MetricExporter)
	return &fcmapProvider{
		moduleName:    moduleName,
		interval:      interval,
		meterProvider: mp,
		clientDesc:    cl,
	}
}

var FcmapMetricDescs = []*MetricDescriptor{
	{
		Key:      "StartTime",
		Name:     "spectrum_flashcopy_start_timestamp",
		Desc:     "Information about the fcmap",
		Unit:     "",
		TypeName: "gauge",
	},
	{
		Key:      "Status",
		Name:     "spectrum_flashcopy_status",
		Desc:     "Information about the fcmap",
		Unit:     "",
		TypeName: "gauge",
	},
	{
		Key:      "Progress",
		Name:     "spectrum_flashcopy_progress",
		Desc:     "Information about the fcmap",
		Unit:     "%",
		TypeName: "gauge",
	},
	{
		Key:      "CleanProgress",
		Name:     "spectrum_flashcopy_clean_progress",
		Desc:     "Information about the fcmap",
		Unit:     "%",
		TypeName: "gauge",
	},
	{
		Key:      "CopyRate",
		Name:     "spectrum_flashcopy_copy_rate",
		Desc:     "Information about the fcmap",
		Unit:     "",
		TypeName: "gauge",
	},
	{
		Key:      "SpaceAllocatedToVdisks",
		Name:     "spectrum_fcmap_allocated_to_vdisks",
		Desc:     "Information about the fcmap",
		Unit:     "mb",
		TypeName: "gauge",
	},
}

func (pv *fcmapProvider) Run(logger *slog.Logger) {
	c := pv.clientDesc.client
	logger.Info("Starting provider", "endpoint", c.Endpoint(), "provider", pv.moduleName)
	meter := pv.meterProvider.Meter(pv.moduleName)

	// Register Metrics...
	var observableMap map[string]metric.Float64Observable
	observableMap = CreateMapMetricDescriptor(meter, FcmapMetricDescs, logger)

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
		data, err := c.GetFlashCopyMap()
		if err != nil {
			logger.Error("Failed to post", "err", err, "endpoint", c.Endpoint(), "provider", pv.moduleName)
			return nil
		}
		if data == nil {
			logger.Warn("data is nil", "provider", pv.moduleName, "endpoint", c.Endpoint())
			return nil
		}

		for _, v := range data {
			fcmapAttr := metric.WithAttributes(
				[]attribute.KeyValue{
					attribute.String("fc.name", v.Name),
					attribute.String("group.name", v.GroupName),
					attribute.String("source.volume", v.SourceVdiskName),
					attribute.String("target.volume", v.TargetVdiskName),
				}...,
			)
			observer.ObserveFloat64(observableMap["StartTime"], float64(v.StartTime.Time().Unix()), clientAttrs, fcmapAttr)
			observer.ObserveFloat64(observableMap["Status"], v.Status.Enum(), clientAttrs, fcmapAttr)
			observer.ObserveFloat64(observableMap["Progress"], v.Progress.Float(), clientAttrs, fcmapAttr)
			observer.ObserveFloat64(observableMap["CleanProgress"], v.CleanProgress.Float(), clientAttrs, fcmapAttr)
			observer.ObserveFloat64(observableMap["CopyRate"], v.CopyRate.Float(), clientAttrs, fcmapAttr)

		}

		return nil
	}, observableArray...)

}
