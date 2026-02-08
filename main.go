package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly/v2"
)

// ElementConfig defines which CSS selectors to extract and their JSON keys
type ElementConfig struct {
	Key      string `json:"key"`
	Selector string `json:"selector"`
	Attr     string `json:"attr,omitempty"` // optional attribute (if empty, use text)
}

func main() {
	// Example URL and elements config
	url := "https://example.com"
	elements := []ElementConfig{
		{Key: "title", Selector: "h1"},
		{Key: "links", Selector: "a", Attr: "href"},
		{Key: "paragraphs", Selector: "p"},
	}

	// Storage for extracted data
	data := make(map[string][]string)
	for _, el := range elements {
		data[el.Key] = []string{}
	}

	// Initialize Colly
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (compatible; GoScraper/1.0)"),
		colly.AllowURLRevisit(),
	)

	// Flexible scraping logic
	for _, el := range elements {
		// Capture loop variable
		el := el
		c.OnHTML(el.Selector, func(e *colly.HTMLElement) {
			var value string
			if el.Attr != "" {
				value = e.Attr(el.Attr)
			} else {
				value = e.Text
			}
			data[el.Key] = append(data[el.Key], value)
		})
	}

	// Error handling
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Request to %s failed: %v", r.Request.URL, err)
	})

	// Start scraping
	fmt.Printf("Scraping %s...\n", url)
	err := c.Visit(url)
	if err != nil {
		log.Fatal(err)
	}

	// Save JSON
	file, err := os.Create("output.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Scraping complete! Data saved to output.json")
}
