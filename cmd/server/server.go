package main

import (
	"log"
	"net/http"

	"github.com/mrkovshik/yametrics/api/counter"
	"github.com/mrkovshik/yametrics/api/gauge"
	"github.com/mrkovshik/yametrics/internal/service"
	"github.com/mrkovshik/yametrics/internal/storage"

)

func main() {
	mapStorage := storage.NewMapStorage()

	service := service.NewServiceWithMapStorage(mapStorage, log.Default())
	run(service)

}

func run(s *service.Service) {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/counter/`, counter.Handler(s))
	mux.HandleFunc(`/update/gauge/`, gauge.Handler(s))
	mux.HandleFunc(`/update/`, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid metric type", http.StatusBadRequest)
		return
	})

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
