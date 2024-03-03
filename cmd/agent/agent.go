package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"
)

const (
	metricTypeGauge   = "gauge"
	metricTypeCounter = "counter"
	pollInterval      = 2 * time.Second
	reportInterval    = 10 * time.Second
)

var metricsValues = make(map[string]float64)

func main() {
	var (
		mu            sync.Mutex
		updateCounter int
	)

	go func() {
		for {
			mu.Lock()
			getRuntimeMetrics()
			updateCounter++
			mu.Unlock()
			time.Sleep(pollInterval)
		}
	}()

	for {
		for name, value := range metricsValues {
			go sendMetric(name, fmt.Sprint(value), metricTypeGauge)
			sendMetric("PollCount ", fmt.Sprint(updateCounter), metricTypeCounter)

		}
		getRuntimeMetrics()
		updateCounter++

		time.Sleep(reportInterval)
	}

}

func getRuntimeMetrics() {
	var (
		m runtime.MemStats
	)
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	runtime.ReadMemStats(&m)

	metricsValues["Alloc"] = float64(m.Alloc)
	metricsValues["BuckHashSys"] = float64(m.BuckHashSys)
	metricsValues["Frees"] = float64(m.Frees)
	metricsValues["GCCPUFraction"] = float64(m.GCCPUFraction)
	metricsValues["GCSys"] = float64(m.GCSys)
	metricsValues["HeapAlloc"] = float64(m.HeapAlloc)
	metricsValues["HeapIdle"] = float64(m.HeapIdle)
	metricsValues["HeapInuse"] = float64(m.HeapInuse)
	metricsValues["HeapObjects"] = float64(m.HeapObjects)
	metricsValues["HeapReleased"] = float64(m.HeapReleased)
	metricsValues["HeapSys"] = float64(m.HeapSys)
	metricsValues["LastGC"] = float64(m.LastGC)
	metricsValues["Lookups"] = float64(m.Lookups)
	metricsValues["MCacheInuse"] = float64(m.MCacheInuse)
	metricsValues["MCacheSys"] = float64(m.MCacheSys)
	metricsValues["MSpanInuse"] = float64(m.MSpanInuse)
	metricsValues["MSpanSys"] = float64(m.MSpanSys)
	metricsValues["Mallocs"] = float64(m.Mallocs)
	metricsValues["NextGC"] = float64(m.NextGC)
	metricsValues["NumForcedGC"] = float64(m.NumForcedGC)
	metricsValues["NumGC"] = float64(m.NumGC)
	metricsValues["OtherSys"] = float64(m.OtherSys)
	metricsValues["PauseTotalNs"] = float64(m.PauseTotalNs)
	metricsValues["StackInuse"] = float64(m.StackInuse)
	metricsValues["StackSys"] = float64(m.StackSys)
	metricsValues["Sys"] = float64(m.Sys)
	metricsValues["TotalAlloc"] = float64(m.TotalAlloc)
	metricsValues["RandomValue"] = random.Float64()

}

func sendMetric(name, value, metricType string) {
	url := "http://localhost:8080/update/"
	metricUpdateURL := fmt.Sprintf("%v%v/%v/%v", url, metricType, name, value)
	response, err := http.Post(metricUpdateURL, "text/plain", nil)
	if err != nil {
		fmt.Println(err)
	}
	if response.StatusCode != http.StatusOK {
		fmt.Printf("status code is %v", response.StatusCode)
	}

}
