package main

import (
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

type Item struct {
	Title      string
	ArticleURL string
	ImageURL   string
	Content    string
}

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
		queryItems(target)
	}
}

func queryItems(target ScrapTarget) {
	collector := colly.NewCollector()

	collector.OnHTML(target.baseGoquerySelector, func(e *colly.HTMLElement) {
		e.ForEach(target.itemGoquerySelector, func(_ int, itemElement *colly.HTMLElement) {
			item := &Item{
				Title:      itemElement.ChildText("a[title]"),
				ArticleURL: itemElement.ChildAttr("a[title]", "href"),
				ImageURL:   itemElement.ChildAttr(fmt.Sprintf("%s img", target.imageGoquerySelector), "src"),
			}

			context := colly.NewContext()
			context.Put("item", item)
			itemElement.Request.Ctx.Put("item", item)

			err := itemElement.Request.Visit(item.ArticleURL)

			if err != nil {
				fmt.Println(err)
			}
		})
	})

	collector.OnHTML(target.articleGoquerySelector, func(e *colly.HTMLElement) {
		item := e.Request.Ctx.GetAny("item").(*Item)
		articleHtml, err := e.DOM.Html()

		if err != nil {
			fmt.Println(err)
		}

		item.Content = articleHtml
		fmt.Println(item.Content)
	})

	err := collector.Visit("https://kicker.de")

	if err != nil {
		panic(err)
	}
}
