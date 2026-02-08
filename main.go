package main

import (
	"encoding/json"
	"fmt"
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
	Material   string `json:"material"`
	Series     string `json:"series"`
}

var regexWeight = regexp.MustCompile(`(?im)\b(?:approx(?:imate)?\.?\s*)?weight\s*(?:[:·-])\s*([0-9]{2,4}(?:\s*[-–]\s*[0-9]{2,4})?\s*(?:g|grs?|grams?))\.?`)
var regexShape = regexp.MustCompile(`(?im)\b(?:shape|form)\s*(?:[:·-])\s*(hybrid|diamond|round|geometric|teardrop|tear\s*drop)\b`)
var regexMaterial = regexp.MustCompile(`(?im)[•\s]*\b(?:comp\.?\s*exterior(?:\s*composition)?|exterior\s*comp\.?(?:osition)?|exterior\s*composition|outer\s*comp\.?(?:osition)?|outer\s*composition|outer\s*shell)\s*(?:[:·-])\s*([^•\r\n]+?)(?:\.\s*-\s*|•|\r?\n|$)`)
var brand = "BullPadel"

func ScrapeRacketPage(url string) (Weight string, Shape string, Material string) {

	weight := ""
	shape := ""
	material := ""

	c := colly.NewCollector(colly.AllowedDomains("www.bullpadel.com"))
	c.OnHTML("div[class='description-short']", func(e *colly.HTMLElement) {
		Description := e.Text

		tmpWeight := regexWeight.FindStringSubmatch(Description)
		if len(tmpWeight) >= 2 {
			weight = strings.ReplaceAll(strings.TrimSpace(tmpWeight[1]), " ", "")
		}

		tmpShape := regexShape.FindStringSubmatch(Description)
		if len(tmpShape) >= 2 {
			shape = tmpShape[1]
		}

		tmpMaterial := regexMaterial.FindStringSubmatch(Description)
		if len(tmpMaterial) >= 2 {
			material = tmpMaterial[1]
		}
	})

	if err := c.Visit(url); err != nil {
		fmt.Println("visit error: ", err)
	}
	return weight, shape, material
}

func ScrapeRacket(url string) {
	var items []Racket
	c := colly.NewCollector(colly.AllowedDomains("www.bullpadel.com"))

	series := strings.Split(url, "/")[4]
	series = strings.ReplaceAll(series, "-", " ")
	series = strings.Title(series)
	partCut := strings.Split(series, " ")[0]
	series = strings.ReplaceAll(series, partCut, "")
	series = strings.ReplaceAll(series, " ", "")

	c.OnHTML("div.left-column div[class='thumbnail-container']", func(e *colly.HTMLElement) {
		item := Racket{
			Brand:      brand,
			Model:      e.ChildText("h3"),
			Price:      e.ChildText("span[itemprop='price']"),
			ImageUrl:   e.ChildAttr("img", "src"),
			RacketPage: e.ChildAttr("a", "href"),
		}
		if strings.Contains(item.Model, "PACK") {
			return
		}
		brand = strings.ToUpper(brand)
		item.Model = strings.Replace(item.Model, brand, "", 1)
		item.Model = strings.Replace(item.Model, "RACKET", "", 1)
		item.Model = strings.Replace(item.Model, " ", "", 1)

		if strings.HasPrefix(item.RacketPage, "/") {
			item.RacketPage = "https://www.bullpadel.com" + item.RacketPage
		}

		weight, shape, material := ScrapeRacketPage(item.RacketPage)

		item.Series = series
		item.Weight = weight
		item.Shape = shape

		material = strings.TrimSpace(material)
		material = strings.TrimRight(material, ". ")
		item.Material = material

		items = append(items, item)
	})

	c.OnHTML("nav[class='pagination']", func(e *colly.HTMLElement) {
		nextPage := e.ChildAttr("a[rel='next']", "href")
		if nextPage != "" {
			if err := c.Visit(nextPage); err != nil {
				fmt.Println("next page visit error: ", err)
			}
		}
	})

	if err := c.Visit(url); err != nil {
		fmt.Println("visit error: ", err)
	}

	content, err := json.Marshal(items)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
	if err := os.WriteFile(brand+"Rackets"+series+".json", content, 0644); err != nil {
		fmt.Println("error: ", err)
	}
}

func main() {
	//ScrapeRacket("https://www.bullpadel.com/gb/39-proline")
	//ScrapeRacket("https://www.bullpadel.com/gb/234-ltd-collection")
	fmt.Println("Done")
}
