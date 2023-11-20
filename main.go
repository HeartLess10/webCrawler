package main

import webCrawl "webScraping/crawler"

func main() {

	urlToWebCrawl := "https://scrapeme.live/shop/"
	webCrawl.StartWebCrawlOnSiteUrl(urlToWebCrawl)

}
