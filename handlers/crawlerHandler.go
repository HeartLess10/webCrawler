package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"webScraping/DTOs"
	"webScraping/webCrawler"
)

type GetCrawlerHandler struct {
	handlerLogger *log.Logger
}

func NewGetCrawlerHandler(l *log.Logger) *GetCrawlerHandler {
	return &GetCrawlerHandler{handlerLogger: l}
}

func (crawlerHandler *GetCrawlerHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		(*crawlerHandler).handlerLogger.Fatal("Error: ", err)
	}
	//good url to test: "https://scrapeme.live/shop/"
	var urlToWebCrawl *DTOs.CrawlUrlDto
	json.Unmarshal(body, urlToWebCrawl)

	urls := webCrawler.NewCrawler().StartWebCrawlOnSiteUrl(urlToWebCrawl.Url)
	jsonBytes, err := json.Marshal(urls)
	if err != nil {
		(*crawlerHandler).handlerLogger.Fatal("Error: ", err)
	}
	rw.Write(jsonBytes)
}
