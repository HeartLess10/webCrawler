package main

import (
	webCrawl "webScraping/crawler"
	//"github.com/labstack/echo"
)

func main() {
	//e := echo.New()

	urlToWebCrawl := "https://scrapeme.live/shop/"
	webCrawl.NewCrawler().StartWebCrawlOnSiteUrl(urlToWebCrawl)
	//e.Logger.Fatal(e.Start(":1323"))

}
