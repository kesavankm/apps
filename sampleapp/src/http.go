package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type HttpServerConfig struct {
	ctx         context.Context
	ServicePort int
}

type HttpServer struct {
	ctx         context.Context
	cfg         HttpServerConfig
	servicePort int
	mux         *http.ServeMux
	app         interface{}
}

func (s *HttpServer) httpPostHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[httpPostHandler] Noop")
}

func (s *HttpServer) httpGetHandler(w http.ResponseWriter, r *http.Request) {
	server := strings.Split(r.URL.Path, "/")[1]
	fmt.Printf("Server %s\n", server)
}

func (s *HttpServer) httpHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		s.httpPostHandler(w, r)
		log.Print("POST")
		return
	case "GET":
		s.httpGetHandler(w, r)
		return
	default:
		log.Printf("Req %+v\n", r)
	}
}

func (s *HttpServer) httpInstallRoutes() {
	s.mux.HandleFunc("/", s.httpHandler)
}

func (s *HttpServer) httpStartServer() {
	address := fmt.Sprintf(":%d", s.servicePort)
	err := http.ListenAndServe(address, s.mux)
	if err != nil {
		log.Printf("[startServer]: err listening and serving http. Err %+v", err)
		panic(err)
	}
}

func (s *HttpServer) httpStartListener() error {
	s.mux = http.NewServeMux()
	s.httpInstallRoutes()
	s.httpStartServer()
	return nil
}

func (s *HttpServer) HttpStartService() error {
	return nil
}

func NewHttpServer(cfg HttpServerConfig) *HttpServer {
	s := &HttpServer{ctx: cfg.ctx, cfg: cfg, servicePort: cfg.ServicePort}
	return s
}
