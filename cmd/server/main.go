package main

import (
	"github.com/mrkovshik/yametrics/api/metricsupd"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, metricsupd.Handler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}

}
