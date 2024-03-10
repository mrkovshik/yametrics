package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/mrkovshik/yametrics/internal/metrics"
	"github.com/mrkovshik/yametrics/internal/utl"
	"log"
	"net/http"
	"sync"
	"time"
)

var pollInterval time.Duration
var reportInterval time.Duration
var hostPort *string

func main() {
	if err := parseFlags(); err != nil {
		log.Fatal(err)
	}
	var (
		mu            sync.Mutex
		updateCounter int
		metricsValues = sync.Map{}
		src           = metrics.NewRuntimeMetrics()
	)
	fmt.Println("Running agent on", *hostPort)
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
	url := fmt.Sprintf("http://%v/update/", *hostPort)
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

func parseFlags() error {

	hostPort = flag.String("a", "localhost:8080", "server host and port")
	flag.DurationVar(&pollInterval, "r", 2*time.Second, "metrics polling interval")
	flag.DurationVar(&reportInterval, "p", 10*time.Second, "metrics sending to server interval")
	flag.Parse()
	if !utl.ValidateAddress(*hostPort) {
		return errors.New("need address in a form host:port")
	}
	return nil
}
