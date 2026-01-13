package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

func main() {
	db, err := NewDatabaseConnection()

	if err != nil {
		log.Fatal(err)
	}

	defer CloseConnection(db)

	targets, err := GetScrapTargets(db)

	if err != nil {
		log.Fatal(err)
	}

	for _, target := range targets {
		queryItems(db, target)
	}
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

	err := collector.Visit("https://kicker.de")

	if err != nil {
		panic(err)
	}
}
