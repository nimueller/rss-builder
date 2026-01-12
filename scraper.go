package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

type Item struct {
	Title      string
	ArticleURL string
	ImageURL   string
	Content    string
}

type ScrapTarget struct {
	URL                    string
	baseGoquerySelector    string
	itemGoquerySelector    string
	imageGoquerySelector   string
	articleGoquerySelector string
}

func main() {
	target := ScrapTarget{
		URL:                    "https://kicker.de",
		baseGoquerySelector:    "#kick__ressort > section:nth-child(2)",
		itemGoquerySelector:    ".kick__slidelist__item",
		imageGoquerySelector:   ".kick__slidelist__item_content_picture img",
		articleGoquerySelector: "main.kick__article__content",
	}

	queryItems(target)
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
