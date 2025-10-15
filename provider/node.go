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

type nodeProvider struct {
	moduleName    string
	interval      time.Duration
	meterProvider *sdkMetric.MeterProvider
	clientDesc    *ClientDesc
}

func init() {
	moduleName := "node"
	registProvider(moduleName, &nodeProvider{moduleName: moduleName})
}

func (pv *nodeProvider) IsDefaultEnabled() bool {
	return true
}

func (pv *nodeProvider) NewProvider(cfg *config.SpectrumConfig, moduleName string, cl *ClientDesc) Provider {
	pvConf := cfg.Providers.Node
	enabled := pvConf.GetEnabled(pv.IsDefaultEnabled())
	interval := pvConf.GetInterval()

	if !enabled {
		return nil
	}
	if MetricExporter == nil {
		return nil
	}
	mp := NewMeterProvider(serviceName, interval, MetricExporter)
	return &nodeProvider{
		moduleName:    moduleName,
		interval:      interval,
		meterProvider: mp,
		clientDesc:    cl,
	}
}

var NodeMetricDescs = []*MetricDescriptor{
	{
		Key:      "status",
		Name:     "spectrum_node_status",
		Desc:     "Information about the node",
		Unit:     "",
		TypeName: "gauge",
	},
	{
		Key:      "config",
		Name:     "spectrum_node_config",
		Desc:     "Information about the node",
		Unit:     "",
		TypeName: "gauge",
	},
}

func (pv *nodeProvider) Run(logger *slog.Logger) {
	logger.Info("Starting provider", "endpoint", pv.clientDesc.endpoint, "provider", pv.moduleName)
	meter := pv.meterProvider.Meter(pv.moduleName)

	// Register Metrics...
	var observableMap map[string]metric.Float64Observable
	observableMap = CreateMapMetricDescriptor(meter, NodeMetricDescs, logger)

	// Register Metrics for Observables...
	var observableNode []metric.Observable
	for _, observable := range observableMap {
		observableNode = append(observableNode, observable)
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
		data, err := c.GetNodeCanister()
		if err != nil {
			logger.Error("Failed to post", "err", err, "endpoint", pv.clientDesc.endpoint, "provider", pv.moduleName)
			return nil
		}
		if data == nil {
			logger.Warn("data is nil", "provider", pv.moduleName, "endpoint", pv.clientDesc.endpoint)
			return nil
		}

		for _, v := range data {
			// Node Attributes...
			additionalAttrs := metric.WithAttributes(
				attribute.String("enclosure.id", v.EnclosureId),
				attribute.String("canister.id", v.CanisterId),
				attribute.String("node.id", v.Id),
				attribute.String("node.name", v.Name),
				attribute.String("node.wwnn", v.WWNN),
			)
			var vconfig float64
			if v.ConfigNode.Bool() {
				vconfig = 1
			} else {
				vconfig = 0
			}
			observer.ObserveFloat64(observableMap["status"], v.Status.Enum(), clientAttrs, additionalAttrs)
			observer.ObserveFloat64(observableMap["config"], vconfig, clientAttrs, additionalAttrs)
		}

		return nil
	}, observableNode...)

}
