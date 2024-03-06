package main

import (
	"fmt"
	"github.com/mrkovshik/yametrics/internal/metrics"
	"net/http"
	"sync"
	"time"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

func main() {
	var (
		mu            sync.Mutex
		updateCounter int
		metricsValues = sync.Map{}
		src           = metrics.NewRuntimeMetrics()
	)

	go func() {
		for {
			fmt.Println("Starting to update metrics")
			getMetrics(src, &metricsValues)
			mu.Lock()
			updateCounter++
			mu.Unlock()
			time.Sleep(pollInterval)
		}
	}()
	time.Sleep(1 * time.Second)
	for {
		fmt.Println("Starting to send metrics")
		for name := range metrics.MetricNamesMap {
			value, _ := metricsValues.Load(name)
			sendMetric(name, fmt.Sprint(value), metrics.MetricTypeGauge)
		}
		sendMetric("PollCount", fmt.Sprint(updateCounter), metrics.MetricTypeCounter)
		time.Sleep(reportInterval)
	}

}

func getMetrics(source metrics.MetricSource, m *sync.Map) {
	source.GetMetrics(m)

}

func sendMetric(name, value, metricType string) {
	url := "http://localhost:8080/update/"
	metricUpdateURL := fmt.Sprintf("%v%v/%v/%v", url, metricType, name, value)
	response, err := http.Post(metricUpdateURL, "text/plain", nil)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		fmt.Printf("status code is %v\n", response.StatusCode)
	}

}
