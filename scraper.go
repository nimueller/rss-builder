package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gocolly/colly"
)

func RunScraperPeriodically(ctx context.Context, config Config, db *sql.DB) {
	log.Println("Starting scraper...")
	ticker := time.NewTicker(config.ScraperInterval)
	defer ticker.Stop()
	scrapeOnce(db)
	log.Println("Scraping timer set. It will run every", config.ScraperInterval)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			scrapeOnce(db)
		}
	}
}

func scrapeOnce(db *sql.DB) {
	log.Println("Scraping...")
	targets, err := GetScrapTargets(db)

	if err != nil {
		log.Fatal(err)
	}

	for _, target := range targets {
		queryItems(db, target)
	}
	log.Println("Scraping all targets finished.")
}

func queryItems(db *sql.DB, target ScrapTarget) {
	collector := colly.NewCollector()

	collector.OnHTML(target.baseGoquerySelector, func(e *colly.HTMLElement) {
		e.ForEach(target.itemGoquerySelector, func(_ int, itemElement *colly.HTMLElement) {
			title := itemElement.ChildText("a[title]")
			articleUrl := itemElement.Request.AbsoluteURL(itemElement.ChildAttr("a[title]", "href"))
			imageUrl := itemElement.ChildAttr(fmt.Sprintf("%s img", target.imageGoquerySelector), "src")

			resultId, err := InsertScrapResult(db, target, title, articleUrl, imageUrl)

			if err != nil {
				log.Println(err)
			}

			itemElement.Request.Ctx.Put("item", resultId)

			err = itemElement.Request.Visit(articleUrl)

			if err != nil {
				log.Println(err)
			}
		})
	})

	collector.OnHTML(target.articleGoquerySelector, func(e *colly.HTMLElement) {
		resultId := e.Request.Ctx.GetAny("item").(int64)
		articleHtml, err := e.DOM.Html()

		if err != nil {
			log.Println(err)
		}

		err = UpdateScrapContent(db, resultId, articleHtml)

		if err != nil {
			log.Println(err)
		}
	})

	collector.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	err := collector.Visit(target.URL)

	if err != nil {
		panic(err)
	}
}
