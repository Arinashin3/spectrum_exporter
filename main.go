package main

import (
	"log/slog"
	"os"
	"spectrum_exporter/config"
	"spectrum_exporter/provider"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/common/promslog"
	promslogflag "github.com/prometheus/common/promslog/flag"
)

var (
	configFile = kingpin.Flag("config.file", "Path to config file.").Short('c').Default("config.yml").String()
	logger     *slog.Logger
	isFailed   bool
)

func main() {
	promslogConfig := &promslog.Config{}
	promslogflag.AddFlags(kingpin.CommandLine, promslogConfig)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	logger = promslog.New(promslogConfig)

	logger.Info("Load Configs...")
	var cfg *config.SpectrumConfig
	cfg = config.NewConfiguration()
	err := cfg.LoadFile(configFile)
	if err != nil {
		isFailed = true
		logger.Error("Failed to load config file.", "error", err)
	}
	if !provider.RegistryProviders(cfg, logger) {
		isFailed = true
	}
	if isFailed {
		logger.Info("Load Configs success.")
		os.Exit(1)
	}
	provider.RunProviders(logger)
}
