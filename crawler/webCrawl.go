package webCrawl

import (
	"fmt"
	"sync"
	"time"
	"webScraping/crawler/util"
)

type crawler interface {
	StartWebCrawlOnSiteUrl(string) map[string][]string
	GetStartUrl() string
	GetMapUrls() map[string][]string
}

type webCrawler struct {
	domainName    string
	functionMutex sync.Mutex
	someMapMutex  sync.RWMutex
	mapUrls       map[string][]string
	wg            sync.WaitGroup
}

func NewCrawler() crawler {
	return &webCrawler{"", sync.Mutex{}, sync.RWMutex{}, make(map[string][]string), sync.WaitGroup{}}
}

func (wb *webCrawler) GetMapUrls() map[string][]string {
	return wb.mapUrls
}

func (wb *webCrawler) GetStartUrl() string {
	return wb.domainName
}

func (wb *webCrawler) StartWebCrawlOnSiteUrl(url string) map[string][]string {
	fmt.Printf("Started to web crawl %s domain\n", url)
	wb.domainName = url
	startTime := time.Now()
	wb.crawlWebSite(url)
	wb.wg.Wait()

	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	util.PrintTimeInReadableFormat(elapsedTime)
	fmt.Printf("Finshed web crawling %s domain\n", url)
	// fmt.Println("Finshed crawling all the links:")
	// for url, _ := range wb.mapUrls {
	// 	fmt.Printf("%s   ,", url)

	// 	utils.TestNoDuplicateUrlInMap(url, wb.mapUrls)
	// }

	return wb.mapUrls
}

func (wb *webCrawler) crawlWebSite(crawlUrl string) {
	wb.wg.Add(1)
	defer wb.wg.Done()
	wb.changeMapValue(crawlUrl, util.ExtractSiteURLs(crawlUrl, wb.domainName))
	for _, url := range wb.mapUrls[crawlUrl] {
		wb.someMapMutex.RLock()
		_, ok := wb.mapUrls[url]
		wb.someMapMutex.RUnlock()
		if !ok {
			wb.changeMapValue(url, nil)
			go wb.crawlWebSite(url)
		}
	}
}

func (wb *webCrawler) changeMapValue(key string, value []string) {
	wb.someMapMutex.Lock()
	wb.mapUrls[key] = value
	wb.someMapMutex.Unlock()
}
