package main

import (
	"log"
	"net/http"
	"net/http/pprof"
	"stat-utility/scrapper"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	m := scrapper.New()
	go func() {
		for {
			scrapper.UpdateMetrics(m)
			time.Sleep(1 * time.Second)
		}
	}()
	//Добавил нагрузку чтобы посомтреть как будут менять метрики при сборке мусора
	go func() {
		for {
			_ = make([]byte, 1024*1024)
			time.Sleep(50 * time.Millisecond)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/debug/pprof/", pprof.Index)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
