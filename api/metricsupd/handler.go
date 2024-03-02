package metricsupd

import (
	"github.com/mrkovshik/yametrics/internal/storage"
	"net/http"
	"strings"
)

var metricsMap = storage.NewMemStorage()

func Handler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	urlParts := strings.Split(req.URL.Path, "/")
	metricName := urlParts[3]
	if len(metricName) == 0 {
		http.Error(res, "metric name is missing", http.StatusNotFound)
	}
	metricValue := urlParts[4]
	switch urlParts[2] {
	case "gauge":
		if err := metricsMap.UpdateGauge(metricName, metricValue); err != nil {
			http.Error(res, "Error updating gauge", http.StatusBadRequest)
			return
		}
		res.Write([]byte("Gauge successfully updated"))
	case "counter":
		if err := metricsMap.UpdateCounter(metricName, metricValue); err != nil {
			http.Error(res, "Error updating counter", http.StatusBadRequest)
			return
		}
		res.Write([]byte("Counter successfully updated"))
	default:
		http.Error(res, "invalid metric type", http.StatusBadRequest)
		return

	}
}
