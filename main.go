package main

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

type Racket struct {
	Model      string `json:"model"`
	Price      string `json:"price"`
	ImageUrl   string `json:"imageUrl"`
	RacketPage string `json:"racketPage"`
	Weight     string `json:"weight"`
	Shape      string `json:"shape"`
}

var items []Racket

func Scrap(url string) {
	c := colly.NewCollector(colly.AllowedDomains("www.bullpadel.com"))

	series := strings.Split(url, "/")[4]

	c.OnHTML("div.left-column div[class='thumbnail-container']", func(e *colly.HTMLElement) {
		item := Racket{
			Model:      e.ChildText("h3"),
			Price:      e.ChildText("span[itemprop='price']"),
			ImageUrl:   e.ChildAttr("img", "src"),
			RacketPage: e.ChildAttr("a", "href"),
		}
		item.Model = strings.Replace(item.Model, "RACKET", "", 1)
		item.Model = strings.Replace(item.Model, "PACK", "", 1)
		item.Model = strings.Replace(item.Model, " ", "", 1)
		items = append(items, item)
	})

	c.OnHTML("nav[class='pagination']", func(e *colly.HTMLElement) {
		nextPage := e.ChildAttr("a[rel='next']", "href")
		c.Visit(nextPage)
	})

	c.Visit(url)

	content, err := json.Marshal(items)
	if err != nil {
		panic(err)
	}
	os.WriteFile("BullPadelRackets"+series+".json", content, 0644)
}

func main() {
	Scrap("https://www.bullpadel.com/gb/39-proline")
	Scrap("https://www.bullpadel.com/gb/234-ltd-collection")
}
