package main

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

type Racket struct {
	Brand      string `json:"brand"`
	Model      string `json:"model"`
	Price      string `json:"price"`
	ImageUrl   string `json:"imageUrl"`
	RacketPage string `json:"racketPage"`
	Weight     string `json:"weight"`
	Shape      string `json:"shape"`
}

func ScrapRacketPage(url string) (Weight string, Shape string) {
	regexWeight := regexp.MustCompile(`(?i)\bweight\s*:\s*([0-9]{2,4}\s*(?:[-–]\s*[0-9]{2,4})?\s*g)\b`)
	regexShape := regexp.MustCompile(`(?i)\bshape\s*:\s*([^\r\n•]+)`)

	var weight string = ""
	var shape string = ""

	c := colly.NewCollector(colly.AllowedDomains("www.bullpadel.com"))
	c.OnHTML("div[class='description-short']", func(e *colly.HTMLElement) {
		Description := e.ChildText("p")
		tmpWeight := regexWeight.FindStringSubmatch(Description)
		if len(tmpWeight) >= 2 {
			weight = strings.ReplaceAll(strings.TrimSpace(tmpWeight[1]), " ", "")
		}
		tmpShape := regexShape.FindStringSubmatch(Description)
		if len(tmpShape) >= 2 {
			shape = tmpShape[1]
		}
	})

	c.Visit(url)
	return weight, shape
}

func ScrapRacket(url string) {
	var items []Racket
	c := colly.NewCollector(colly.AllowedDomains("www.bullpadel.com"))

	series := strings.Split(url, "/")[4]

	c.OnHTML("div.left-column div[class='thumbnail-container']", func(e *colly.HTMLElement) {
		item := Racket{
			Brand:      "BullPadel",
			Model:      e.ChildText("h3"),
			Price:      e.ChildText("span[itemprop='price']"),
			ImageUrl:   e.ChildAttr("img", "src"),
			RacketPage: e.ChildAttr("a", "href"),
		}
		item.Model = strings.Replace(item.Model, "BULLPADEL", "", 1)
		item.Model = strings.Replace(item.Model, "RACKET", "", 1)
		item.Model = strings.Replace(item.Model, "PACK", "", 1)
		item.Model = strings.Replace(item.Model, " ", "", 1)
		weight, shape := ScrapRacketPage(item.RacketPage)
		item.Weight = weight
		item.Shape = shape
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
	ScrapRacket("https://www.bullpadel.com/gb/39-proline")
	ScrapRacket("https://www.bullpadel.com/gb/234-ltd-collection")
}
