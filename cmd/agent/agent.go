package main

import (
	"fmt"
	config "github.com/mrkovshik/yametrics/internal/config/agent"
	"github.com/mrkovshik/yametrics/internal/metrics"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	var (
		cfg           config.AgentConfig
		mu            sync.Mutex
		updateCounter int
		metricsValues = sync.Map{}
		src           = metrics.NewRuntimeMetrics()
	)
	if err := cfg.GetConfigs(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Running agent on %v\npoll interval = %v\nreport interval = %v\n", cfg.Address, cfg.PollInterval, cfg.ReportInterval)
	go func() {
		for {
			log.Println("Starting to update metrics")
			src.StoreMetrics(&metricsValues)
			mu.Lock()
			updateCounter++
			mu.Unlock()
			time.Sleep(time.Duration(cfg.PollInterval) * time.Second)
		}
	}()
	time.Sleep(1 * time.Second)
	for {
		log.Println("Starting to send metrics")
		for name := range metrics.MetricNamesMap {
			value, _ := metricsValues.Load(name)
			sendMetric(cfg.Address, name, fmt.Sprint(value), metrics.MetricTypeGauge)
		}
		sendMetric(cfg.Address, "PollCount", fmt.Sprint(updateCounter), metrics.MetricTypeCounter)
		time.Sleep(time.Duration(cfg.ReportInterval) * time.Second)
	}

}

func sendMetric(addr, name, value, metricType string) {
	url := fmt.Sprintf("http://%v/update/", addr)
	metricUpdateURL := fmt.Sprintf("%v%v/%v/%v", url, metricType, name, value)
	response, err := http.Post(metricUpdateURL, "text/plain", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		fmt.Printf("status code is %v, while sending %v:%v:%v\n", response.StatusCode, metricType, name, value)
		return
	}
}
