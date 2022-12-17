package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metrics struct {
	registry         *prometheus.Registry
	numActiveDevices prometheus.Gauge
}

func NewMetrics() *metrics {
	// Create a non-global registry.
	reg := prometheus.NewRegistry()

	m := &metrics{
		numActiveDevices: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "ns_num_active_devices",
			Help: "Number of active devices",
		}),
		// hdFailures: prometheus.NewCounterVec(
		// 	prometheus.CounterOpts{
		// 		Name: "hd_errors_total",
		// 		Help: "Number of hard-disk errors.",
		// 	},
		// 	[]string{"device"},
		// ),
	}
	reg.MustRegister(m.numActiveDevices)

	// m.numActiveDevices.Set(123.4)
	// reg.MustRegister(m.hdFailures)

	m.registry = reg
	return m
}

func (ns *NetScan) httpInstallMetricsRoutes() {
	reg := ns.metrics.registry
	ns.mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
}
