package main

import (
	"fmt"
	"github.com/mrkovshik/yametrics/internal/metrics"
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

var metricsValues sync.Map

func main() {
	var (
		mu            sync.Mutex
		updateCounter int
	)

	go func() {
		for {
			fmt.Println("Starting to update metrics")
			getRuntimeMetrics()
			mu.Lock()
			updateCounter++
			mu.Unlock()
			time.Sleep(pollInterval)
		}
	}()
	time.Sleep(1 * time.Second)
	for {
		fmt.Println("Starting to send metrics")
		for name, _ := range metrics.MetricNamesMap {
			value, _ := metricsValues.Load(name)
			sendMetric(name, fmt.Sprint(value), metricTypeGauge)
		}
		sendMetric("PollCount ", fmt.Sprint(updateCounter), metricTypeCounter)
		time.Sleep(reportInterval)
	}

}

func getRuntimeMetrics() {
	var (
		m runtime.MemStats
	)
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	runtime.ReadMemStats(&m)

	metricsValues.Store("Alloc", float64(m.Alloc))
	metricsValues.Store("BuckHashSys", float64(m.BuckHashSys))
	metricsValues.Store("Frees", float64(m.Frees))
	metricsValues.Store("GCCPUFraction", m.GCCPUFraction)
	metricsValues.Store("GCSys", float64(m.GCSys))
	metricsValues.Store("HeapAlloc", float64(m.HeapAlloc))
	metricsValues.Store("HeapIdle", float64(m.HeapIdle))
	metricsValues.Store("HeapInuse", float64(m.HeapInuse))
	metricsValues.Store("HeapObjects", float64(m.HeapObjects))
	metricsValues.Store("HeapReleased", float64(m.HeapReleased))
	metricsValues.Store("HeapSys", float64(m.HeapSys))
	metricsValues.Store("LastGC", float64(m.LastGC))
	metricsValues.Store("Lookups", float64(m.Lookups))
	metricsValues.Store("MCacheInuse", float64(m.MCacheInuse))
	metricsValues.Store("MCacheSys", float64(m.MCacheSys))
	metricsValues.Store("MSpanInuse", float64(m.MSpanInuse))
	metricsValues.Store("MSpanSys", float64(m.MSpanSys))
	metricsValues.Store("Mallocs", float64(m.Mallocs))
	metricsValues.Store("NextGC", float64(m.NextGC))
	metricsValues.Store("NumForcedGC", float64(m.NumForcedGC))
	metricsValues.Store("NumGC", float64(m.NumGC))
	metricsValues.Store("OtherSys", float64(m.OtherSys))
	metricsValues.Store("PauseTotalNs", float64(m.PauseTotalNs))
	metricsValues.Store("StackInuse", float64(m.StackInuse))
	metricsValues.Store("StackSys", float64(m.StackSys))
	metricsValues.Store("Sys", float64(m.Sys))
	metricsValues.Store("TotalAlloc", float64(m.TotalAlloc))
	metricsValues.Store("RandomValue", random.Float64())

}

func sendMetric(name, value, metricType string) {
	url := "http://localhost:8080/update/"
	metricUpdateURL := fmt.Sprintf("%v%v/%v/%v", url, metricType, name, value)
	response, err := http.Post(metricUpdateURL, "text/plain", nil)
	if err != nil {
		fmt.Println(err)
	}
	if response.StatusCode != http.StatusOK {
		fmt.Printf("status code is %v\n", response.StatusCode)
	}

}
