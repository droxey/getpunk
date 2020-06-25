package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/droxey/getpunk/logger"
	"github.com/gocolly/colly"
)

// Shield represents a badge on shields.io.
// See the schema details @ https://shields.io/endpoint.
type Shield struct {
	SchemaVersion int    `json:"schemaVersion"`
	Label         string `json:"label"`
	Message       string `json:"message"`
	IsError       bool   `json:"isError"`
	Logo          string `json:"logo"`
	Style         string `json:"style"`
	CacheSeconds  int    `json:"cacheSeconds"`
}

// scrapeHandler initializes Colly, makes a request to the
// URL we're scraping, and returns the scraped data as JSON.
func scrapeHandler(w http.ResponseWriter, r *http.Request) {
	c := colly.NewCollector()

	c.OnHTML("body > div.main_content.project > div:nth-child(1) > div.row.m-t-50 > div > div:nth-child(4) > div:nth-child(2) > b", func(e *colly.HTMLElement) {
		// Clean up the scraped value, returned in the folowing format: 0.87 Ξ ($209).
		fullPrice := strings.Replace(e.Text, "Ξ", "", -1)
		logger.Log.Info("Price scraped: " + fullPrice)

		// Set the struct's properties.
		priceShield := new(Shield)
		priceShield.SchemaVersion = 1
		priceShield.Label = ""
		priceShield.Message = fullPrice
		priceShield.IsError = false
		priceShield.Style = "flat"
		priceShield.Logo = "image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABgAAAAYCAYAAADgdz34AAAArElEQVR42mNgGAWjgOZAzHzJf3yYYsOlnE/gxWRbQozhJFkAUwyjRXRq4GxkjMVw4n2AbDhU6P/tjQ0oGNkyEJ9kC5AMZ4BqRsFIhqJj4ixANzyTxxWOkQ3EIk4yAGvs4Y2GY2QfoIuTbcFaxQQ4RrYAXZzcrIAR/jBMcvjjsgA9FeGwiHoWYLOQahYcURSAY2pZgNcX1DCcATnl0MJwFB9QK+XgDSJqG06xBQCuo19OCl9XGgAAAABJRU5ErkJggg=="
		priceShield.CacheSeconds = 300

		// Serialize the struct to JSON.
		bf := bytes.NewBuffer([]byte{})
		jsonEncoder := json.NewEncoder(bf)
		jsonEncoder.SetEscapeHTML(false)
		jsonEncoder.Encode(priceShield)

		// Return a JSON response containing our data.
		w.Header().Set("Content-Type", "application/json")
		w.Write(bf.Bytes())
	})

	c.Visit("https://www.larvalabs.com/cryptopunks/accountinfo?account=0xf11dfe0321485d9892d20420ae7cb5b3eb9fbb06")
}
func main() {
	host := "0.0.0.0:8888"
	http.HandleFunc("/", scrapeHandler)
	logger.Log.Info("Server started: http://" + host)

	err := http.ListenAndServe(host, nil)
	if err != nil {
		logger.Log.Error(err)
		return
	}
}
