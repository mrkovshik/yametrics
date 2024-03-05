package gauge

import (
	"github.com/mrkovshik/yametrics/internal/metrics"
	"net/http"
	"strings"

	"github.com/mrkovshik/yametrics/internal/service"
	"github.com/mrkovshik/yametrics/internal/storage"
)

func Handler(s *service.Service) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(res, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
			return
		}
		urlParts := strings.Split(req.URL.Path, "/")
		if len(urlParts) < 5 || !verifyGaugeName(urlParts[3]) {
			http.Error(res, "Data is missing", http.StatusNotFound)
			return
		}

		gauge, err := storage.NewGauge(urlParts[3], urlParts[4])
		if err != nil {
			http.Error(res, "wrong value format", http.StatusBadRequest)
			return
		}
		if err := gauge.Update(s.Storage); err != nil {
			http.Error(res, "Error updating counter", http.StatusBadRequest)
			return
		}
		res.Write([]byte("gauge successfully updated"))

	}
}

func verifyGaugeName(name string) bool {
	_, ok := metrics.MetricNamesMap[name]
	return ok
}
