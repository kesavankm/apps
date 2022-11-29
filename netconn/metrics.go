package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	totalPacketLossGg prometheus.Gauge
	connLossGg        prometheus.Gauge
	outageCountCtr    prometheus.Counter
	avgRTTGg          prometheus.Gauge
	minRTTGg          prometheus.Gauge
	maxRTTGg          prometheus.Gauge

	targetServer string
	packetsSent  int
	packetsRcvd  int
	connLoss     int
	avgRTT       float64
}

func (nc *NetConn) httpTestConnLoss(w http.ResponseWriter, r *http.Request) {
	log.Printf("[httpTestConnLoss] Enter v7")
	nc.conns["google.com"].metrics.connLossGg.Set(float64(1))
}

func (nc *NetConn) httpOutageCount(w http.ResponseWriter, r *http.Request) {
	log.Printf("[httpOutageCount] Enter v7")
	nc.conns["google.com"].metrics.outageCountCtr.Inc()
}

func (nc *NetConn) httpNumbersHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[httpNumbersHandler] Enter v7")
	nc.conns["google.com"].recordMetrics()
}

func (nc *NetConn) httpTestCountersHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[httpTestCountersHandler] Enter v7")
	nc.conns["google.com"].recordMetrics()
}

func (c *ConnInfo) NewConnMetrics() {
	// Metrics init
	m := Metrics{}

	m.totalPacketLossGg = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nc_total_loss_packet_count",
		Help: "The total number of packets lost",
	})

	m.connLossGg = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nc_connection_loss",
		Help: "Connection loss observed",
	})

	m.outageCountCtr = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "nc_outage_count",
		Help: "Number of times outage observed",
	})

	m.avgRTTGg = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nc_avg_rtt",
		Help: "Average RTT",
	})

	m.minRTTGg = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nc_min_rtt",
		Help: "Minimum RTT",
	})

	m.maxRTTGg = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nc_max_rtt",
		Help: "Maximum RTT",
	})

	prometheus.MustRegister(m.totalPacketLossGg)
	prometheus.MustRegister(m.connLossGg)
	prometheus.MustRegister(m.outageCountCtr)
	prometheus.MustRegister(m.avgRTTGg)
	prometheus.MustRegister(m.minRTTGg)
	prometheus.MustRegister(m.maxRTTGg)

	c.metrics = &m
}

func (c *ConnInfo) collectStats() {
	go func() {
		for {
			c.pingerStats = c.pinger.Statistics()
			c.lastPollTimeStamp = time.Now()
			time.Sleep(time.Duration(c.pollingInterval) * time.Second)
			c.totalLossCount = c.pingerStats.PacketsSent - c.pingerStats.PacketsRecv
		}
	}()
}

func (c *ConnInfo) recordMetrics() {
	var wg sync.WaitGroup

	m := c.metrics
	wg.Add(1)
	go func() {
		for {
			// totalLoss.Set(float64(c.totalLossCount))
			if c.totalLossCount > 1 {
				m.totalPacketLossGg.Set(float64(c.totalLossCount))
			}
			m.avgRTTGg.Set(float64(c.pingerStats.AvgRtt / time.Millisecond))
			m.minRTTGg.Set(float64(c.pingerStats.MinRtt / time.Millisecond))
			m.maxRTTGg.Set(float64(c.pingerStats.MaxRtt / time.Millisecond))
			time.Sleep(30 * time.Second)
		}
	}()
	log.Printf("wow")
	wg.Add(1)
	go func() {
		for {
			if c.totalLossCount > c.prevLossCount+1 {
				log.Printf("\n\n  LOSS: Diff: Stats %d,%d\n LossCt %d\n Conn %+v\n\n",
					c.pingerStats.PacketsSent, c.pingerStats.PacketsRecv, c.totalLossCount, c)
				if c.lossObserved == 0 {
					// Mark as outage only if No loss is observed until now
					m.outageCountCtr.Inc()
				}
				c.lossObserved = 1
				c.prevLossCount = c.totalLossCount
			} else {
				if c.lossObserved == 1 {
					log.Printf("\n\n  LOSS reset \n\n")
					c.lossObserved = 0
				}
			}
			// connLoss.Set(float64(c.lossObserved))
			m.connLossGg.Set(float64(c.lossObserved))
			time.Sleep(10 * time.Second)
		}
	}()
	log.Printf("wow2")
	wg.Wait()
	log.Printf("wow3")
}

func (nc *NetConn) httpInstallMetricsRoutes() {
	nc.mux.HandleFunc("/numbers", nc.httpNumbersHandler)
	nc.mux.Handle("/metrics", promhttp.Handler())
	nc.mux.HandleFunc("/connLoss", nc.httpTestConnLoss)
	nc.mux.HandleFunc("/outageCount", nc.httpOutageCount)
}
