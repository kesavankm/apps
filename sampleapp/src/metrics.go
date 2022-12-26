package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metrics struct {
	app      *SampleApp
	registry *prometheus.Registry
	gauge1   prometheus.Gauge
}

func NewMetrics() *metrics {
	// Create a non-global registry.
	reg := prometheus.NewRegistry()

	m := &metrics{
		gauge1: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "app_data_metric_1",
			Help: "Sample app data metric",
		}),
	}
	reg.MustRegister(m.gauge1)

	m.registry = reg
	return m
}

func (m *metrics) httpInstallMetricsRoutes() {
	reg := m.registry
	if m.app != nil {
		m.app.mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	}
}
