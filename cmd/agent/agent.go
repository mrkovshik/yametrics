package main

import (
	"context"
	"log"

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

	go agent.PollMetrics()
	go agent.SendMetrics(ctx)
	select {}
}
