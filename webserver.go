package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
)

func RunWebServer(ctx context.Context, config Config) {
	log.Println("Starting webserver...")
	srv := &http.Server{Addr: config.WebServerHost + ":" + strconv.Itoa(config.WebServerPort)}

	go func() {
		err := srv.ListenAndServe()

		if err != nil {
			panic(err)
		}
	}()

	log.Println("Webserver started. Serving is listening on", config.WebServerHost, "port", config.WebServerPort)
	<-ctx.Done()

	err := srv.Shutdown(context.Background())

	if err != nil {
		log.Println(err)
	}
}
