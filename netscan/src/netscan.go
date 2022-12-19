package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Ullaakut/nmap"

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

func (ns *NetScan) scanAndReport() {
	scan := func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		scanner, err := nmap.NewScanner(
			nmap.WithTargets("192.168.86.0/24"),
			nmap.WithPingScan(),
			nmap.WithContext(ctx),
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

		// // Use the results to print an example output
		// for _, host := range result.Hosts {
		// 	if len(host.Ports) == 0 || len(host.Addresses) == 0 {
		// 		continue
		// 	}

		// 	fmt.Printf("Host %q:\n", host.Addresses[0])

		// 	for _, port := range host.Ports {
		// 		fmt.Printf("\tPort %d/%s %s %s\n", port.ID, port.Protocol, port.State, port.Service.Name)
		// 	}
		// }
		fmt.Printf("Nmap done: %d hosts up scanned in %3f seconds\n", len(result.Hosts), result.Stats.Finished.Elapsed)
		return len(result.Hosts)
	}

	go func() {
		for {
			numDevices := scan()
			ns.metrics.numActiveDevices.Set(float64(numDevices))
			time.Sleep(5 * time.Minute)
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
	apiVersion := "0.0.6"
	log.Printf("[main] Enter version %s\n", apiVersion)
	ctx := context.Background()
	log.Printf("main")
	n := NewNetScan(ctx)
	run(n)
}
