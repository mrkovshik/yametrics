package main

import (
	"fmt"
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

	agent := service.NewAgent(src, log.Default(), cfg, strg) //TODO: implement zap logger
	fmt.Printf("Running agent on %v\npoll interval = %v\nreport interval = %v\n", agent.Config.Address, agent.Config.PollInterval, agent.Config.ReportInterval)
	go agent.PollMetrics()
	time.Sleep(1 * time.Second)
	for {
		if err := agent.SendMetric(); err != nil {
			log.Println(err)
		}
		time.Sleep(time.Duration(agent.Config.ReportInterval) * time.Second)
	}
}
