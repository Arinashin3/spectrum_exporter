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

type hostProvider struct {
	moduleName    string
	interval      time.Duration
	meterProvider *sdkMetric.MeterProvider
	clientDesc    *ClientDesc
}

func init() {
	moduleName := "host"
	registProvider(moduleName, &hostProvider{moduleName: moduleName})
}

func (pv *hostProvider) IsDefaultEnabled() bool {
	return true
}

func (pv *hostProvider) NewProvider(cfg *config.SpectrumConfig, moduleName string, cl *ClientDesc) Provider {
	pvConf := cfg.Providers.Host
	enabled := pvConf.GetEnabled(pv.IsDefaultEnabled())
	interval := pvConf.GetInterval()

	if !enabled {
		return nil
	}
	if MetricExporter == nil {
		return nil
	}
	mp := NewMeterProvider(serviceName, interval, MetricExporter)
	return &hostProvider{
		moduleName:    moduleName,
		interval:      interval,
		meterProvider: mp,
		clientDesc:    cl,
	}
}

var HostMetricDescs = []*MetricDescriptor{
	{
		Key:      "status",
		Name:     "spectrum_host_status",
		Desc:     "Information about the host",
		Unit:     "",
		TypeName: "gauge",
	},
	{
		Key:      "port_count",
		Name:     "spectrum_host_port_count",
		Desc:     "Information about the host",
		Unit:     "",
		TypeName: "gauge",
	},
	{
		Key:      "iogrp_count",
		Name:     "spectrum_host_iogrp_count",
		Desc:     "Information about the host",
		Unit:     "",
		TypeName: "gauge",
	},
}

func (pv *hostProvider) Run(logger *slog.Logger) {
	c := pv.clientDesc.client
	logger.Info("Starting provider", "endpoint", c.Endpoint(), "provider", pv.moduleName)
	meter := pv.meterProvider.Meter(pv.moduleName)

	// Register Metrics...
	var observableMap map[string]metric.Float64Observable
	observableMap = CreateMapMetricDescriptor(meter, HostMetricDescs, logger)

	// Register Metrics for Observables...
	var observableHost []metric.Observable
	for _, observable := range observableMap {
		observableHost = append(observableHost, observable)
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
		data, err := c.GetHost()
		if err != nil {
			logger.Error("Failed to post", "err", err, "endpoint", c.Endpoint(), "provider", pv.moduleName)
			return nil
		}
		if data == nil {
			logger.Warn("data is nil", "provider", pv.moduleName, "endpoint", c.Endpoint())
			return nil
		}

		for _, v := range data {
			// Host Attributes...
			additionalAttrs := metric.WithAttributes(
				attribute.String("target.id", v.Id),
				attribute.String("target.name", v.Name),
				attribute.String("protocol", v.Protocol),
			)
			observer.ObserveFloat64(observableMap["status"], v.Status.Enum(), clientAttrs, additionalAttrs)
			observer.ObserveFloat64(observableMap["port_count"], v.PortCount.Float(), clientAttrs, additionalAttrs)
			observer.ObserveFloat64(observableMap["iogrp_count"], v.IogrpCount.Float(), clientAttrs, additionalAttrs)
		}

		return nil
	}, observableHost...)

}
