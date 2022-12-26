package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Ullaakut/nmap/v2"

	alerts "github.com/kesavankm/alerter"
)

type NetScan struct {
	ctx            context.Context
	servicePort    int
	mux            *http.ServeMux
	targetNetworks []string
	scanFrequency  time.Duration
	alertClient    alerts.AlertIntf
	metrics        *metrics
}

type NetScanReport struct {
	numActiveDevices int
}

func countByOS(result *nmap.Run) {
	var (
		noHostOS   int
		noHostAddr int
	)
	// Count the number of each OS for all hosts.
	for _, host := range result.Hosts {
		if len(host.Addresses) == 0 {
			noHostAddr++
			// continue
		}

		if len(host.OS.Matches) == 0 {
			noHostOS++
			// continue
		}
		log.Printf("Host %s, Total matches %d\n", host.Addresses[0], len(host.OS.Matches))
		var hostName string
		if len(host.Hostnames) > 0 {
			hostName = host.Hostnames[0].Name
		}
		var bestOSMatch, bestClassFamily, bestClassVendor string
		var hostAccuracy int
		if len(host.OS.Matches) > 0 {
			hostAccuracy = host.OS.Matches[0].Accuracy
			bestOSMatch = host.OS.Matches[0].Name
			bestClassVendor = host.OS.Matches[0].Classes[0].Vendor
			bestClassFamily = host.OS.Matches[0].Classes[0].Family
		}
		log.Printf("  Host %s, Acc: %d, OS: %s, Vendor: %s, OSFamily: %s\n\n",
			hostName, hostAccuracy, bestOSMatch,
			bestClassVendor, bestClassFamily)
	}

	fmt.Printf("Discovered hosts : up:%d(total:%d), noHostOSMatch:%d, noHostAddr:%d.\n",
		result.Stats.Hosts.Up, result.Stats.Hosts.Total, noHostOS, noHostAddr)
}

func (ns *NetScan) scanAndReport() {
	scan := func() NetScanReport {
		log.Printf("Initiating scan")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		scanner, err := nmap.NewScanner(
			nmap.WithTargets("192.168.86.0/24"),
			// nmap.WithPingScan(),
			nmap.WithContext(ctx),
			nmap.WithFastMode(),
			nmap.WithOSDetection(),
			nmap.WithOSScanLimit(),
		)

		if err != nil {
			log.Fatalf("unable to create nmap scanner: %v", err)
		}

		result, warnings, err := scanner.Run()
		if err != nil {
			log.Fatalf("unable to run nmap scan: %v", err)
		}

		if warnings != nil {
			log.Printf("Warnings: \n %v", warnings)
		}

		countByOS(result)

		fmt.Printf("Nmap done: %d hosts up scanned in %3f seconds\n", len(result.Hosts), result.Stats.Finished.Elapsed)
		return NetScanReport{numActiveDevices: len(result.Hosts)}
	}

	go func() {
		for {
			nsReport := scan()
			// ns.metrics.numActiveDevices.Set(float64(numDevices))
			ns.metrics.PublishMetrics(nsReport)
			time.Sleep(ns.scanFrequency)
		}
	}()
}

func (ns *NetScan) httpGetHandler(w http.ResponseWriter, r *http.Request) {
	server := strings.Split(r.URL.Path, "/")[1]
	fmt.Printf("Server %s\n", server)
}

func (ns *NetScan) httpHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		// post.postPostHandler(w, r)
		log.Print("POST")
		return
	case "GET":
		ns.httpGetHandler(w, r)
		return
	default:
		log.Printf("Req %+v\n", r)
	}
}

func (ns *NetScan) httpInstallRoutes() {
	ns.mux.HandleFunc("/", ns.httpHandler)
	ns.httpInstallMetricsRoutes()
}

func (ns *NetScan) httpStartServer() {
	// log.Printf("[startServer] Enter")
	// defer log.Printf("[startServer] Exit")
	address := fmt.Sprintf(":%d", ns.servicePort)
	err := http.ListenAndServe(address, ns.mux)
	if err != nil {
		log.Printf("[startServer]: err listening and serving http. Err %+v", err)
		panic(err)
	}
}

func (ns *NetScan) httpStartListener() error {
	ns.mux = http.NewServeMux()
	ns.httpInstallRoutes()
	ns.httpStartServer()
	return nil
}

func NewNetScan(ctx context.Context) *NetScan {
	var err error
	ns := &NetScan{ctx: ctx}
	ns.servicePort = 10080
	ns.alertClient, err = alerts.NewAlertClient(ctx, alerts.AlertConfig{ClientType: "slack", ApiToken: os.Getenv("SLACK_BOT_TOKEN")})
	if err != nil {
		panic(err)
	}
	ns.scanFrequency = 5 * time.Minute

	ns.metrics = NewMetrics()

	return ns
}

func getDefaultNetwork(intf string) string {
	return ""
}

func run(ns *NetScan) {
	defaultNetwork := getDefaultNetwork("")
	ns.targetNetworks = append(ns.targetNetworks, defaultNetwork)
	go ns.scanAndReport()
	ns.httpStartListener()
}

func main() {
	apiVersion := "0.0.7i"
	log.Printf("[main] Enter version %s\n", apiVersion)
	ctx := context.Background()
	log.Printf("main")
	n := NewNetScan(ctx)
	run(n)
}
