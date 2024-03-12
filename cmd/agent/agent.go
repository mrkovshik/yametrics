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
		strg  = storage.NewAgentMapStorage()
		src   = metrics.NewRuntimeMetrics()
		agent = service.NewAgent(src, log.Default(), config.AgentConfig{}, strg) //TODO: implement zap logger
	)
	if err := agent.Config.GetConfigs(); err != nil {
		log.Fatal(err)
	}
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
