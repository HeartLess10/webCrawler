package webCrawler

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"
	"webScraping/webCrawler/util"
)

type Crawler interface {
	StartWebCrawlOnSiteUrl(string) map[string][]string
	GetDomainUrl() string
	ToJson(io.Writer) error
}

type webCrawler struct {
	domainUrl     string              `json:"-"` //there is no need for json - becuse the field is not exported
	functionMutex sync.Mutex          `json:"-"`
	someMapMutex  sync.RWMutex        `json:"-"`
	MapUrls       map[string][]string `json:"mapUrls ,Omitempty"`
	wg            sync.WaitGroup      `json:"-"`
}

func NewCrawler() Crawler {
	return &webCrawler{"", sync.Mutex{}, sync.RWMutex{}, make(map[string][]string), sync.WaitGroup{}}
}

func (wc *webCrawler) GetDomainUrl() string {
	return wc.domainUrl
}

func (wc *webCrawler) StartWebCrawlOnSiteUrl(url string) map[string][]string {
	fmt.Printf("Started to web crawl %s domain\n", url)
	wc.domainUrl = url
	startTime := time.Now()
	wc.crawlWebSite(url)
	wc.wg.Wait()

	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	util.PrintTimeInReadableFormat(elapsedTime)
	fmt.Printf("Finshed web crawling %s domain\n", url)
	// fmt.Println("Finshed crawling all the links:")
	// for url, _ := range wc.mapUrls {
	// 	fmt.Printf("%s   ,", url)

	// 	utils.TestNoDuplicateUrlInMap(url, wc.mapUrls)
	// }
	util.TestNoDuplicateUrlInMap(wc.domainUrl, wc.MapUrls)
	return wc.MapUrls
}

func (wc *webCrawler) crawlWebSite(crawlUrl string) {
	wc.wg.Add(1)
	defer wc.wg.Done()
	wc.changeMapValue(crawlUrl, util.ExtractSiteURLs(crawlUrl, wc.domainUrl))
	for _, url := range wc.MapUrls[crawlUrl] {
		wc.someMapMutex.RLock()
		_, ok := wc.MapUrls[url]
		wc.someMapMutex.RUnlock()
		if !ok {
			wc.changeMapValue(url, nil)
			go wc.crawlWebSite(url)
		}
	}
}

func (wc *webCrawler) changeMapValue(key string, value []string) {
	wc.someMapMutex.Lock()
	wc.MapUrls[key] = value
	wc.someMapMutex.Unlock()
}

func (wc *webCrawler) ToJson(w io.Writer) error {
	ecoder := json.NewEncoder(w)
	return ecoder.Encode(wc)
}
