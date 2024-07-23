package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	config "github.com/mrkovshik/yametrics/internal/config/agent"
	"github.com/mrkovshik/yametrics/internal/metrics"
	service "github.com/mrkovshik/yametrics/internal/service/agent"
	"github.com/mrkovshik/yametrics/internal/storage"
)

var (
	buildVersion, buildDate, buildCommit string
)

func main() {
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)
	// Initialize storage and metrics source
	strg := storage.NewInMemoryStorage()
	src := metrics.NewRuntimeMetrics()

	// Initialize logging with zap
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("zap.NewDevelopment", zap.Error(err))
	}

	// Flushes buffered log entries before program exits
	defer logger.Sync() //nolint:all
	sugar := logger.Sugar()

	// Get configuration settings
	cfg, err := config.GetConfigs()
	if err != nil {
		logger.Fatal("config.GetConfigs", zap.Error(err))
	}

	ctx, stopServices := context.WithCancel(context.Background())
	defer stopServices()

	// Create agent instance with dependencies
	agent := service.NewAgent(src, &cfg, strg, sugar)

	// Log agent configuration
	sugar.Infof(
		"Running agent on %v\n"+
			"key = %v\n"+
			"key is set = %v\n"+
			"poll interval = %v\n"+
			"poll interval is set = %v\n"+
			"report interval = %v\n"+
			"report interval is set = %v\n"+
			"crypto key = %v\n"+
			"crypto key is set = %v\n"+
			"config file path = %v\n"+
			"config file path is set = %v\n"+
			"rate limit = %v\n"+
			"rate limit is set = %v\n",
		&cfg.Address,
		cfg.Key,
		cfg.KeyIsSet,
		cfg.PollInterval,
		cfg.PollIntervalIsSet,
		cfg.ReportInterval,
		cfg.ReportIntervalIsSet,
		cfg.CryptoKey,
		cfg.CryptoKeyIsSet,
		cfg.ConfigFilePath,
		cfg.ConfigFilePathIsSet,
		cfg.RateLimit,
		cfg.RateLimitIsSet,
	)

	// Create tickers for polling and sending metrics
	pollTicker := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	defer pollTicker.Stop()
	pollUtilTicker := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	defer pollUtilTicker.Stop()
	sendTicker := time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)
	defer sendTicker.Stop()

	pollMetricsStopped := make(chan struct{})
	pollUtilMetricsStopped := make(chan struct{})
	sendMetricsStopped := make(chan struct{})

	// Start goroutines for polling and sending metrics
	go agent.PollMetrics(pollTicker.C, pollMetricsStopped)
	go agent.PollUtilMetrics(pollUtilTicker.C, pollUtilMetricsStopped)
	go agent.SendMetrics(ctx, sendTicker.C, sendMetricsStopped)

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	select {
	case <-sigs:
		sugar.Info("Received shutdown signal")
		pollTicker.Stop()
		pollUtilTicker.Stop()
		sendTicker.Stop()
	}
	<-pollMetricsStopped
	<-pollUtilMetricsStopped
	<-sendMetricsStopped
}
