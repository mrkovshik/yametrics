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
	var (
		strg = storage.NewMapStorage()
		src  = metrics.NewRuntimeMetrics()
	)

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("zap.NewDevelopment",
			zap.Error(err))
	}
	if err != nil {
		logger.Fatal("metrics.NewUtilMetrics",
			zap.Error(err))
	}
	cfg, err := config.GetConfigs()
	if err != nil {
		logger.Fatal("config.GetConfigs",
			zap.Error(err))
	}
	defer logger.Sync() //nolint:all
	sugar := logger.Sugar()
	ctx := context.Background()
	agent := service.NewAgent(src, cfg, strg, sugar)
	sugar.Infof("Running agent on %v\npoll interval = %v\nreport interval = %v\n", cfg.Address, cfg.PollInterval, cfg.ReportInterval)

	pollTicker := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	defer pollTicker.Stop()
	pollUtilTicker := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	defer pollUtilTicker.Stop()
	sendTicker := time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)
	defer sendTicker.Stop()

	go agent.PollMetrics(pollTicker.C)
	go agent.PollUitlMetrics(pollUtilTicker.C)
	go agent.SendMetrics(ctx, sendTicker.C)
	if cfg.LoadRateLimit > 0 {
		loadTicker := time.NewTicker(time.Duration(cfg.LoadRateLimit) * time.Millisecond)
		go agent.LoadServer(loadTicker.C)
	}
	select {}
}
