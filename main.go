package main

import (
	"encoding/json"
	"os"

	"github.com/gocolly/colly"
)

type Racket struct {
	Model    string `json:"model"`
	Price    string `json:"price"`
	ImageUrl string `json:"imageUrl"`
}

func main() {
	c := colly.NewCollector(colly.AllowedDomains("www.bullpadel.com"))

	var items []Racket

	c.OnHTML("div.left-column div[class='thumbnail-container']", func(e *colly.HTMLElement) {
		item := Racket{
			Model:    e.ChildText("h3"),
			Price:    e.ChildText("span[itemprop='price']"),
			ImageUrl: e.ChildAttr("img", "src"),
		}
		items = append(items, item)
	})

	c.OnHTML("nav[class='pagination']", func(e *colly.HTMLElement) {
		nextPage := e.ChildAttr("a[rel='next']", "href")
		c.Visit(nextPage)
	})

	c.Visit("https://www.bullpadel.com/gb/39-proline")

	content, err := json.Marshal(items)
	if err != nil {
		panic(err)
	}
	os.WriteFile("BullPadelRackets.json", content, 0644)
}
