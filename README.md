# BullPadel Racket Scraper (Go)

A simple Go scraper that collects BullPadel padel racket data from bullpadel.com and outputs it to JSON files.

## What it does
- Visits a category/series page (e.g. Proline, LTD Collection)
- Extracts for each racket:
    - `brand`, `model`, `price`, `imageUrl`, `racketPage`
- Opens each racket detail page and parses:
    - `weight`, `shape`, `material` (from the product description)
- Writes results to:
    - `BullPadelRackets<Series>.json` (e.g. `BullPadelRacketsProline.json`)

## Output format
Each item looks like:
```json
{
  "brand": "BullPadel",
  "model": "VERTEX 05",
  "price": "€339.99",
  "imageUrl": "https://...",
  "racketPage": "https://...",
  "weight": "365-375g",
  "shape": "Diamond",
  "material": "X-Tend Carbon 12K",
  "series": "Proline"
}