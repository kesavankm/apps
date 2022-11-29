package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-ping/ping"
)

type connStats struct {
	// PacketsSent
	pktsRcv int
	pass    int
	fail    int
}

type ConnInfo struct {
	ctx               context.Context
	active            bool
	nc                *NetConn
	targetServer      string
	pinger            *ping.Pinger
	pingerStats       *ping.Statistics
	pingCount         int
	metrics           *Metrics
	lossObserved      int
	prevLossCount     int
	pollingInterval   int
	lastPollTimeStamp time.Time
	totalLossCount    int
}

type NetConn struct {
	ctx           context.Context
	mux           *http.ServeMux
	targetServers []string
	conns         map[string]*ConnInfo
	pollFrequency int // in seconds
}

func (nc *NetConn) httpInstallRoutes() {
	nc.mux.HandleFunc("/", nc.httpHandler)
	nc.httpInstallMetricsRoutes()
}

func (nc *NetConn) httpGetHandler(w http.ResponseWriter, r *http.Request) {
	// log.Printf("Get. Req %+v\n", r)
	// log.Printf("host %+v, URLPath %s, URl %+v\n\n", r.Host, r.URL.Path, r.URL)
	server := strings.Split(r.URL.Path, "/")[1]
	if conn, ok := nc.conns[server]; ok {
		// conn.dispStats()
		stats := conn.getStats()
		// log.Printf("Stats: \n%+v\n", stats)
		// log.Printf("Conn: \n%+v\n", conn)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
		json.NewEncoder(w).Encode(conn)
	} else {
		log.Printf("Warn. No open conn for server %+v\n", server)
	}
}

func (nc *NetConn) httpHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		// post.postPostHandler(w, r)
		log.Print("POST")
		return
	case "GET":
		nc.httpGetHandler(w, r)
		return
	default:
		log.Printf("Req %+v\n", r)
	}
}

func (nc *NetConn) httpStartServer() {
	// log.Printf("[startServer] Enter")
	// defer log.Printf("[startServer] Exit")
	address := fmt.Sprintf(":%d", 12192)
	err := http.ListenAndServe(address, nc.mux)
	if err != nil {
		log.Printf("[startServer]: err listening and serving http. Err %+v", err)
		panic(err)
	}
}

func (nc *NetConn) httpStartListener() error {
	nc.mux = http.NewServeMux()
	nc.httpInstallRoutes()
	nc.httpStartServer()
	return nil
}

func (c *ConnInfo) pingSetup() {
	// log.Printf("[pingSetup] c %+v", c)
	var err error
	c.pinger, err = ping.NewPinger(c.targetServer)
	if err != nil {
		panic(err)
	}
	c.pinger.RecordRtts = false
	c.pinger.SetPrivileged(true)
	c.pingCount = 100
}

func (c *ConnInfo) doPing() {
	err := c.pinger.Run() // blocks until finished
	if err != nil {
		panic(err)
	}
}

func (c *ConnInfo) startPing() {
	// log.Printf("Starting Ping service for %s\n", c.targetServer)
	c.pingSetup()
	c.active = true
	go c.collectStats()
	go c.recordMetrics()
	c.doPing()
}

func (c *ConnInfo) stopPing() {
	// log.Printf("Stopping Ping service for %s\n", c.targetServer)
	c.active = true
	c.pinger.Stop()
}

func (c *ConnInfo) dispStats() {
	log.Printf("Stats for %s: %+v\n", c.targetServer, c.pinger.Statistics())
}

func (c *ConnInfo) getStats() *ping.Statistics {
	return c.pinger.Statistics()
}

func (nc *NetConn) stopAllMonitor() {
	for _, v := range nc.conns {
		v.stopPing()
	}
}

func (nc *NetConn) startMonitor() {
	for _, v := range nc.conns {
		v.startPing()
	}
}

func NewConnInfo(ctx context.Context, server string) *ConnInfo {
	c := ConnInfo{ctx: ctx, targetServer: server, pollingInterval: 10}
	c.NewConnMetrics()
	return &c
}

func (c *ConnInfo) SetParentNetConn(nc *NetConn) {
	c.nc = nc
}

func NewNetConn(ctx context.Context) *NetConn {
	return &NetConn{ctx: ctx}
}

func run(nc *NetConn) {
	defaultServer := "google.com"
	nc.targetServers = append(nc.targetServers, defaultServer)
	nc.conns = make(map[string]*ConnInfo)
	nc.conns[defaultServer] = NewConnInfo(nc.ctx, defaultServer)
	nc.conns[defaultServer].SetParentNetConn(nc)
	go nc.doSpeedtest()
	go nc.startMonitor()
	nc.httpStartListener()
}

func main() {
	log.Printf("[main] Enter v0.0.17")
	ctx := context.Background()
	log.Printf("main")
	n := NewNetConn(ctx)
	run(n)
}
