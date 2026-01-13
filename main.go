package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Config struct {
	WebServerHost   string
	WebServerPort   int
	ScraperInterval time.Duration
}

func main() {
	config := Config{
		WebServerHost:   "localhost",
		WebServerPort:   8080,
		ScraperInterval: 10 * time.Minute,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go RunScraperPeriodically(ctx, config)
	go RunWebServer(ctx, config)

	log.Println("Started. Press Ctrl+C to stop.")
	waitForSignal()
}

func waitForSignal() {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
}
