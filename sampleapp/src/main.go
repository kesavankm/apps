package main

import (
	"context"
	"log"
)

func main() {
	apiVersion := "0.0.7"
	log.Printf("[main] Enter version %s\n", apiVersion)
	ctx := context.Background()
	log.Printf("main")
	n := NewSampleApp(ctx)
	n.Run()
}
