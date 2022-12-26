package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metrics struct {
	registry         *prometheus.Registry
	numActiveDevices prometheus.Gauge
}

func (m *metrics) PublishMetrics(nsReport NetScanReport) error {
	m.numActiveDevices.Set(float64(nsReport.numActiveDevices))
	return nil
}

func NewMetrics() *metrics {
	// Create a non-global registry.
	reg := prometheus.NewRegistry()

	m := &metrics{
		numActiveDevices: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "ns_num_active_devices",
			Help: "Number of active devices",
		}),
	}
	reg.MustRegister(m.numActiveDevices)

	m.registry = reg
	return m
}

func (ns *NetScan) httpInstallMetricsRoutes() {
	reg := ns.metrics.registry
	ns.mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
}
