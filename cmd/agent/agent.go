package main

import (
	"log"
	"time"

	"go.uber.org/zap"

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
	agent := service.NewAgent(src, cfg, strg, sugar)
	sugar.Infof("Running agent on %v\npoll interval = %v\nreport interval = %v\n", agent.Config.Address, agent.Config.PollInterval, agent.Config.ReportInterval)
	go agent.PollMetrics()
	time.Sleep(1 * time.Second) //TODO: Костыль. Потом написать сюда WG
	for {
		if err := agent.SendMetric(); err != nil {
			logger.Error("agent.SendMetric",
				zap.Error(err))
		}
		time.Sleep(time.Duration(agent.Config.ReportInterval) * time.Second)
	}
}
