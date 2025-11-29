package scrapper

import (
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	AllocGauge       prometheus.Gauge
	TotallAllocGauge prometheus.Gauge
	SysGauge         prometheus.Gauge
	NumGCTotall      prometheus.Gauge
	LastGC           prometheus.Gauge
}

func New() *Metrics {
	m := &Metrics{AllocGauge: prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "memstats_alloc",
		Help: "runtime.MemStats.Alloc",
	}),
		TotallAllocGauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "total_allocation",
			Help: "total byte allocation size",
		}),
		SysGauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "system_bytes",
			Help: "Bytes get from os",
		}),
		NumGCTotall: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "total_gc_counts",
			Help: "numbers of gc runs",
		}),
		LastGC: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "last_gc_time_seconds",
			Help: "last gc time in seconds",
		}),
	}

	prometheus.MustRegister(m.AllocGauge)
	prometheus.MustRegister(m.TotallAllocGauge)
	prometheus.MustRegister(m.SysGauge)
	prometheus.MustRegister(m.LastGC)
	prometheus.MustRegister(m.NumGCTotall)

	return m
}

func UpdateMetrics(metr *Metrics) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	metr.AllocGauge.Set(float64(m.Alloc))
	metr.TotallAllocGauge.Set(float64(m.TotalAlloc))
	metr.SysGauge.Set(float64(m.Sys))
	metr.NumGCTotall.Set(float64(m.NumGC))

	if m.NumGC > 0 {
		metr.LastGC.Set(float64(m.LastGC) / 1e9)
	}
}
