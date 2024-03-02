package main

import (
	"github.com/mrkovshik/yametrics/api/update_metrics"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, update_metrics.Handler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}

}
