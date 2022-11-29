package main

import (
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/showwin/speedtest-go/speedtest"
)

func (nc *NetConn) doSpeedtest() {
	log.Printf("[doSpeedtest] Enter")
	user, _ := speedtest.FetchUserInfo()
	// Get a list of servers near a specified location
	// user.SetLocationByCity("Tokyo")
	// user.SetLocation("Osaka", 34.6952, 135.5006)

	downloadSpeedGg := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nc_download_speed",
		Help: "Download speed in Mbps",
	})

	uploadSpeedGg := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nc_upload_speed",
		Help: "Upload speed in Mbps",
	})

	latencyGg := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nc_latency",
		Help: "Download speed in Mbps",
	})

	prometheus.MustRegister(downloadSpeedGg)
	prometheus.MustRegister(uploadSpeedGg)
	prometheus.MustRegister(latencyGg)

	go func() {
		count := 0
		for {
			serverList, _ := speedtest.FetchServers(user)
			targets, _ := serverList.FindServer([]int{})
			s := targets[0]
			s.PingTest()
			s.DownloadTest(false)
			s.UploadTest(false)

			log.Printf("Count: %d, Host %s, Latency: %s, Download: %f, Upload: %f\n",
				count, s.Host, s.Latency, s.DLSpeed, s.ULSpeed)
			count++
			downloadSpeedGg.Set(float64(s.DLSpeed))
			uploadSpeedGg.Set(float64(s.ULSpeed))
			latencyGg.Set(float64(s.Latency))
			time.Sleep(5 * time.Minute)
		}
	}()
}
