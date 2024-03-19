package main

import (
	"fmt"
	"go.uber.org/zap"
	"log"
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

	cfg, err := config.GetConfigs()
	if err != nil {
		log.Fatal(err)
	}
	logger, _ := zap.NewProduction()

	defer logger.Sync() //nolint:all

	agent := service.NewAgent(src, cfg, strg, logger)
	logger.Info(fmt.Sprintf("Running agent on %v\npoll interval = %v\nreport interval = %v\n", agent.Config.Address, agent.Config.PollInterval, agent.Config.ReportInterval))
	go agent.PollMetrics()
	time.Sleep(1 * time.Second)
	for {
		if err := agent.SendMetric(); err != nil {
			logger.Fatal("agent.SendMetric",
				zap.Error(err))
		}
		time.Sleep(time.Duration(agent.Config.ReportInterval) * time.Second)
	}
}
