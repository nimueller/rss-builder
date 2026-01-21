package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	DatabaseUrl     string        `env:"DATABASE_URL" envDefault:"postgres://postgres:postgres@localhost:5432/postgres"`
	WebServerHost   string        `env:"WEBSERVER_HOST" envDefault:"127.0.0.1"`
	WebServerPort   int           `env:"WEBSERVER_PORT" envDefault:"8080"`
	ScraperInterval time.Duration `env:"SCRAPER_INTERVAL" envDefault:"10m"`
}

func main() {
	config := parseConfig()
	db := connectToDatabase(config)
	defer CloseConnectionPool(db)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go RunScraperPeriodically(ctx, config, db)
	go RunWebServer(ctx, config, db)

	log.Println("Started. Press Ctrl+C to stop.")
	waitForSignal()
}

func parseConfig() Config {
	config := Config{}
	err := env.Parse(&config)

	if err != nil {
		log.Fatal(err)
	}

	return config
}

func connectToDatabase(config Config) *sql.DB {
	db, err := NewDatabaseConnectionPool(config)

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func waitForSignal() {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
}
