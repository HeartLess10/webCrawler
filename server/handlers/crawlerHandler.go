package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"webScraping/DTOs"
	"webScraping/webCrawler"
)

type CrawlerHandler struct {
	handlerLogger *log.Logger
}

func NewCrawlerHandler(l *log.Logger) *CrawlerHandler {
	return &CrawlerHandler{handlerLogger: l}
}

func (crawlerHandler *CrawlerHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		crawlerHandler.getCrawlerHandler(rw, r)
	}

	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (crawlerHandler *CrawlerHandler) getCrawlerHandler(rw http.ResponseWriter, r *http.Request) {
	// body, err := io.ReadAll(r.Body)
	// if err != nil {
	// 	(*crawlerHandler).handlerLogger.Fatal("Error: ", err)
	// }
	//good url to test: "https://scrapeme.live/shop/"
	var urlToWebCrawl *DTOs.CrawlUrlDto
	err := json.NewDecoder(r.Body).Decode(&urlToWebCrawl) //decoding the body from json to struct
	if err != nil {
		(*crawlerHandler).handlerLogger.Println("Error: ", err)
		return
	}
	webCrawler := webCrawler.NewCrawler()
	webCrawler.StartWebCrawlOnSiteUrl(urlToWebCrawl.Url)
	err = webCrawler.ToJson(rw)
	if err != nil {
		(*crawlerHandler).handlerLogger.Println("Error: ", err)
		return
	}
}
