package counter

import (
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
		if len(urlParts) < 5 || (urlParts[3] != "PollCount" && urlParts[3] != "testCounter") {
			http.Error(res, "Data is missing", http.StatusNotFound)
			return
		}

		counter, err := storage.NewCounter(urlParts[3], urlParts[4])
		if err != nil {
			http.Error(res, "wrong value format", http.StatusBadRequest)
			return
		}

		if err := counter.Update(s.Storage); err != nil {
			http.Error(res, "Error updating counter", http.StatusBadRequest)
			return
		}
		res.Write([]byte("Counter successfully updated"))
	}
}
