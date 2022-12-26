package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	alerts "github.com/kesavankm/alerter"
)

type SampleAppConfig struct {
	EnableHTTPService bool
	ServicePort       int
	EnableAlerts      bool
	EnableMetrics     bool
}

type SampleApp struct {
	ctx            context.Context
	cfg            SampleAppConfig
	servicePort    int
	mux            *http.ServeMux
	s              *HttpServer
	targetNetworks []string
	pollFrequency  time.Duration
	alertClient    alerts.AlertIntf
	metrics        *metrics
}

func (sa *SampleApp) SetPollFrequency(t time.Duration) {
	sa.pollFrequency = t
}

func NewSampleApp(ctx context.Context) *SampleApp {
	var err error
	sa := &SampleApp{ctx: ctx}
	sa.servicePort = 10080
	sa.alertClient, err = alerts.NewAlertClient(ctx, alerts.AlertConfig{ClientType: "slack", ApiToken: os.Getenv("SLACK_BOT_TOKEN")})
	if err != nil {
		panic(err)
	}

	sa.s = NewHttpServer(HttpServerConfig{ctx: ctx, ServicePort: sa.servicePort})
	sa.metrics = NewMetrics()
	sa.metrics.app = sa

	return sa
}

func (sa *SampleApp) HelloWorld() {
	log.Printf("Hello World")
}

func (sa *SampleApp) Run() {
}
