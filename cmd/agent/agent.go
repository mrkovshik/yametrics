package main

import (
	"fmt"
	"github.com/mrkovshik/yametrics/internal/metrics"
	"net/http"
	"sync"
	"time"
)

func main() {
	parseFlags()
	var (
		mu            sync.Mutex
		updateCounter int
		metricsValues = sync.Map{}
		src           = metrics.NewRuntimeMetrics()
	)
	fmt.Println("Running agent on", addr.String())
	go func() {
		for {
			fmt.Println("Starting to update metrics")
			src.StoreMetrics(&metricsValues)
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

func sendMetric(name, value, metricType string) {
	url := fmt.Sprintf("http://%v/update/", addr.String())
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
