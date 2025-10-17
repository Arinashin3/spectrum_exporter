package provider

import (
	"context"
	"log/slog"
	"spectrum_exporter/config"
	"spectrum_exporter/gospectrum"
	"spectrum_exporter/gospectrum/types"
	"time"

	"go.opentelemetry.io/otel/log"
	sdkLog "go.opentelemetry.io/otel/sdk/log"
)

func init() {
	moduleName := "eventlog"
	registProvider(moduleName, &eventlogProvider{moduleName: moduleName})
}

func (pv *eventlogProvider) IsDefaultEnabled() bool {
	return true
}

func (pv *eventlogProvider) NewProvider(cfg *config.SpectrumConfig, moduleName string, cl *ClientDesc) Provider {
	pvConf := cfg.Providers.Eventlog
	enabled := pvConf.GetEnabled(pv.IsDefaultEnabled())
	interval := pvConf.GetInterval()
	//
	if !enabled {
		return nil
	}
	if LogExporter == nil {
		return nil
	}
	lp := NewLoggerProvider(serviceName, interval, LogExporter)
	return &eventlogProvider{
		moduleName:     moduleName,
		interval:       interval,
		loggerProvider: lp,
		clientDesc:     cl,
	}
}

type eventlogProvider struct {
	moduleName     string
	interval       time.Duration
	level          int
	loggerProvider *sdkLog.LoggerProvider
	clientDesc     *ClientDesc
}

func (pv *eventlogProvider) Run(logger *slog.Logger) {
	c := pv.clientDesc.client
	logger.Info("Starting provider", "endpoint", c.Endpoint(), "provider", pv.moduleName)
	ctx := context.Background()
	lp := pv.loggerProvider

	// Create And Set Options...
	opts := gospectrum.NewLsEventLogOptions()
	ctime := time.Now().Add(-1 * time.Hour).UTC()

	for {
		if !c.IsLogin() {
			time.Sleep(pv.interval)
			continue
		}
		pvlogger := lp.Logger(pv.moduleName, log.WithInstrumentationAttributes(pv.clientDesc.hostLabels...))

		opts.AddFilterValue("last_timestamp>=" + ctime.Format(types.TimeLayout))

		data, err := c.GetEventLog(opts)
		if err != nil {
			time.Sleep(pv.interval)
			continue
		}

		if data == nil {
			time.Sleep(pv.interval)
			continue
		}

		for _, eventlog := range data {
			record := log.Record{}
			eventTimestamp := eventlog.LastTimestamp.Time()
			if err != nil {
				logger.Error("Error parsing timestamp", "err", err)
			}
			record.SetObservedTimestamp(eventTimestamp)
			if eventlog.ErrorCode != "" {
				record.AddAttributes(
					log.String("level", "ALERT"),
				)
			} else {
				record.AddAttributes(
					log.String("level", "INFO"),
				)
			}
			record.AddAttributes(
				log.String("error.code", eventlog.ErrorCode),
				log.String("message.id", eventlog.EventId),
				log.String("object.name", eventlog.ObjectName),
				log.String("status", eventlog.Status),
			)
			record.SetBody(log.StringValue(eventlog.Description))

			pvlogger.Emit(ctx, record)
		}
		ctime = time.Now().Local()
		time.Sleep(pv.interval)
	}
}
