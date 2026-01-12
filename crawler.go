package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

type Item struct {
	Title    string
	ImageURL string
}

func main() {
	collector := colly.NewCollector()

	queryItems(collector, "#kick__ressort > section:nth-child(2)", ".kick__slidelist__item", ".kick__slidelist__item_content_picture")

	err := collector.Visit("https://kicker.de")

	if err != nil {
		panic(err)
	}
}

func queryItems(collector *colly.Collector, baseGoquerySelector string, itemGoquerySelector string, imageSelector string) {
	collector.OnHTML(baseGoquerySelector, func(e *colly.HTMLElement) {
		e.ForEach(itemGoquerySelector, func(_ int, itemElement *colly.HTMLElement) {
			item := Item{}
			getItem(itemElement, &item, imageSelector)

			println()
			println("\t", item.Title)
			println("\t", item.ImageURL)
		})
	})
}

func getItem(itemElement *colly.HTMLElement, item *Item, imageSelector string) {
	item.Title = itemElement.ChildText("a[title]")
	item.ImageURL = itemElement.ChildAttr(fmt.Sprintf("%s img", imageSelector), "src")
}
