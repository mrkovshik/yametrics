package main

import (
	"go.uber.org/zap"
	"time"

	config "github.com/mrkovshik/yametrics/internal/config/agent"
	"github.com/mrkovshik/yametrics/internal/metrics"
	service "github.com/mrkovshik/yametrics/internal/service/agent"
	storage "github.com/mrkovshik/yametrics/internal/storage/agent"
)

func main() {
	var (
		strg = storage.NewAgentMapStorage()
		src  = metrics.NewRuntimeMetrics()
	)
	logger, _ := zap.NewProduction()
	cfg, err := config.GetConfigs()
	if err != nil {
		logger.Fatal("config.GetConfigs",
			zap.Error(err))
	}
	defer logger.Sync() //nolint:all
	agent := service.NewAgent(src, cfg, strg, logger)
	logger.Info("Running agent on",
		zap.String("Host", agent.Config.Address),
		zap.Int("Poll interval", agent.Config.PollInterval),
		zap.Int("Report interval", agent.Config.ReportInterval))
	go agent.PollMetrics()
	time.Sleep(1 * time.Second)
	for {
		if err := agent.SendMetric(); err != nil {
			logger.Error("agent.SendMetric",
				zap.Error(err))
		}
		time.Sleep(time.Duration(agent.Config.ReportInterval) * time.Second)
	}
}
