package main

import (
	"context"
	"log"
	"time"

	"go.uber.org/zap"

	config "github.com/mrkovshik/yametrics/internal/config/agent"
	"github.com/mrkovshik/yametrics/internal/metrics"
	service "github.com/mrkovshik/yametrics/internal/service/agent"
	"github.com/mrkovshik/yametrics/internal/storage"
)

func main() {
	// Initialize storage and metrics source
	strg := storage.NewMapStorage()
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

	// Create agent instance with dependencies
	agent := service.NewAgent(src, cfg, strg, sugar)

	// Log agent configuration
	sugar.Infof("Running agent on %v\npoll interval = %v\nreport interval = %v\n", cfg.Address, cfg.PollInterval, cfg.ReportInterval)

	// Create tickers for polling and sending metrics
	pollTicker := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	defer pollTicker.Stop()
	pollUtilTicker := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	defer pollUtilTicker.Stop()
	sendTicker := time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)
	defer sendTicker.Stop()

	// Start goroutines for polling and sending metrics
	go agent.PollMetrics(pollTicker.C)
	go agent.PollUitlMetrics(pollUtilTicker.C)
	go agent.SendMetrics(context.Background(), sendTicker.C)

	// Block indefinitely to keep the agent running
	select {}
}
