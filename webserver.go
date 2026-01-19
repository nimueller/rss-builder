package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Content     string `xml:"http://purl.org/rss/1.0/modules/content/ encoded,omitempty"`
}

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

func RunWebServer(ctx context.Context, config Config, db *sql.DB) {
	log.Println("Starting webserver...")
	mux := http.NewServeMux()
	mux.HandleFunc("/rss", func(writer http.ResponseWriter, request *http.Request) {
		handleRssRoute(writer, request, db)
	})
	srv := &http.Server{Addr: config.WebServerHost + ":" + strconv.Itoa(config.WebServerPort), Handler: mux}

	go func() {
		err := srv.ListenAndServe()

		if err != nil {
			panic(err)
		}
	}()

	log.Println("Webserver started. Serving is listening on", config.WebServerHost, "port", config.WebServerPort)
	<-ctx.Done()

	err := srv.Shutdown(ctx)

	if err != nil {
		log.Println(err)
	}
}

func handleRssRoute(w http.ResponseWriter, request *http.Request, db *sql.DB) {
	targetId, err := strconv.ParseInt(request.URL.Query().Get("targetId"), 10, 64)

	if err != nil {
		http.Error(w, fmt.Sprintf("Expected targetId %v as integer 64", targetId), http.StatusBadRequest)
		log.Println(err)
		return
	}

	result, err := GetLatestScrapResult(db, targetId)

	if err != nil {
		http.Error(w, "Error while reading results", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	items := make([]Item, 0, len(result))
	for _, scrapResult := range result {
		var content string

		if !scrapResult.Content.Valid {
			content = ""
		} else {
			content = scrapResult.Content.String
		}

		items = append(items, Item{
			Title:       scrapResult.Title,
			Link:        scrapResult.ArticleUrl,
			Description: scrapResult.Title,
			PubDate:     time.Now().Format(time.RFC1123),
			Content:     content,
		})
	}

	rss := RSS{
		Version: "2.0",
		Channel: Channel{
			Title:       "Mein RSS Feed",
			Link:        "https://example.com",
			Description: "Dies ist ein Testfeed",
			Items:       items,
		},
	}

	xmlData, err := xml.MarshalIndent(rss, "", "  ")
	if err != nil {
		http.Error(w, "Fehler beim Erzeugen des RSS", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/rss+xml")
	w.Write([]byte(xml.Header))
	w.Write(xmlData)
}
