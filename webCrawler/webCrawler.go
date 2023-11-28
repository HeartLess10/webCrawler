package webCrawler

import (
	"fmt"
	"sync"
	"time"
	"webScraping/webCrawler/util"
)

type Crawler interface {
	StartWebCrawlOnSiteUrl(string) map[string][]string
	GetStartUrl() string
	GetMapUrls() map[string][]string
}

type webCrawler struct {
	domainUrl     string              `json:"-"`
	functionMutex sync.Mutex          `json:"-"`
	someMapMutex  sync.RWMutex        `json:"-"`
	mapUrls       map[string][]string `json:"urls ,Omitempty"`
	wg            sync.WaitGroup      `json:"-"`
}

func NewCrawler() Crawler {
	return &webCrawler{"", sync.Mutex{}, sync.RWMutex{}, make(map[string][]string), sync.WaitGroup{}}
}

func (wb *webCrawler) GetMapUrls() map[string][]string {
	return wb.mapUrls
}

func (wb *webCrawler) GetStartUrl() string {
	return wb.domainUrl
}

func (wb *webCrawler) StartWebCrawlOnSiteUrl(url string) map[string][]string {
	fmt.Printf("Started to web crawl %s domain\n", url)
	wb.domainUrl = url
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
	util.TestNoDuplicateUrlInMap(wb.domainUrl, wb.mapUrls)
	return wb.mapUrls
}

func (wb *webCrawler) crawlWebSite(crawlUrl string) {
	wb.wg.Add(1)
	defer wb.wg.Done()
	wb.changeMapValue(crawlUrl, util.ExtractSiteURLs(crawlUrl, wb.domainUrl))
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
