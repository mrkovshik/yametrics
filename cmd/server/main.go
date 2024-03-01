package main

import (
	"github.com/mrkovshik/yametrics/internal/storage"
	"net/http"
)

var metricsMap = storage.NewMemStorage()

func updateMetric(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.Write([]byte("Only GET requests are allowed!"))
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	metricName := req.FormValue("name")
	metricValue := req.FormValue("value")
	switch req.FormValue("type") {
	case "gauge":
		if err := metricsMap.UpdateGauge(metricName, metricValue); err != nil {
			res.Write([]byte("Error updating gauge"))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	case "counter":
		if err := metricsMap.UpdateCounter(metricName, metricValue); err != nil {
			res.Write([]byte("Error updating counter"))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		res.Write([]byte("invalid metric type"))
		res.WriteHeader(http.StatusBadRequest)
		return

	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/{type}/{name}/{value}`, updateMetric)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}

}
