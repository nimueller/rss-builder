package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	DatabaseUrl     string        `env:"DATABASE_URL" envDefault:"postgres://postgres:postgres@localhost:5432/postgres"`
	WebServerHost   string        `env:"WEBSERVER_HOST" envDefault:"localhost"`
	WebServerPort   int           `env:"WEBSERVER_PORT" envDefault:"8080"`
	ScraperInterval time.Duration `env:"SCRAPER_INTERVAL" envDefault:"10m"`
}

func main() {
	config := Config{}
	err := env.Parse(&config)

	if err != nil {
		log.Fatal(err)
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
