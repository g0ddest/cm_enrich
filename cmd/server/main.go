package main

import (
	"cm_enrich/internal/config"
	"cm_enrich/internal/handlers"
	"cm_enrich/internal/metrics"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	go metrics.SetupPrometheus()

	log.Println("Starting message processing...")
	for {
		handlers.ProcessMessages(cfg)
	}
}
