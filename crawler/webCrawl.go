package webCrawl

import (
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	startUrl      string
	functionMutex sync.Mutex
	someMapMutex                      = sync.RWMutex{}
	myMap         map[string][]string = make(map[string][]string)
	wg            sync.WaitGroup
)

func StartWebCrawlOnSiteUrl(url string) map[string][]string {
	startUrl = url
	startTime := time.Now()
	crawlWebSite(url)
	wg.Wait()

	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	printTimeInReadableFormat(elapsedTime)
	fmt.Println("Finshed crawling all the links:")
	for url, _ := range myMap {
		fmt.Printf("%s   ,", url)

		testNoDuplicateUrlInMap(url)
	}

	return myMap
}

func crawlWebSite(crawlUrl string) {
	wg.Add(1)
	defer wg.Done()
	changeMapValue(crawlUrl, extractSiteURLs(crawlUrl))
	for _, url := range myMap[crawlUrl] {
		someMapMutex.RLock()
		_, ok := myMap[url]
		someMapMutex.RUnlock()
		if !ok {
			changeMapValue(url, nil)
			go crawlWebSite(url)
		}
	}
}

func changeMapValue(key string, value []string) {
	someMapMutex.Lock()
	myMap[key] = value
	someMapMutex.Unlock()
}

func extractSiteURLs(url string) []string {
	var documentLinks []string = nil
	var response *http.Response = getSiteHTML(url)
	if response != nil {
		defer response.Body.Close()
		if response.StatusCode == 200 {
			document, err := goquery.NewDocumentFromReader(response.Body)
			if err == nil {
				documentLinks = scrapeDocumentForLinks(document)
				fmt.Printf("url:%s,\n urls:%v\n\n", url, documentLinks)
			} else {
				printError("When trying to find links", url)
			}
		} else {
			printError("Status code "+response.Status, url)
		}
	} else {
		printError("Bad response", url)
	}
	return documentLinks
}

func getSiteHTML(siteUrl string) *http.Response {
	response, err := http.Get(siteUrl)
	if err != nil {
		printError("Url was not valid", siteUrl)
		return nil
	}
	return response
}

func scrapeDocumentForLinks(doc *goquery.Document) []string {
	var result []string
	doc.Find("body a").Each(func(index int, item *goquery.Selection) {
		link, _ := item.Attr("href")
		if !slices.Contains(result, link) && strings.Contains(link, startUrl) {
			result = append(result, link)
		}
	})
	return result
}

func printTimeInReadableFormat(elapsedTime time.Duration) {
	// Round the duration to the nearest millisecond
	roundedDuration := elapsedTime.Round(time.Millisecond)

	// Extract hours, minutes, seconds, and milliseconds
	hours := int(roundedDuration.Hours())
	minutes := int(roundedDuration.Minutes()) % 60
	seconds := int(roundedDuration.Seconds()) % 60
	milliseconds := roundedDuration.Milliseconds() % 1000

	fmt.Printf("Elapsed time: %02d:%02d:%02d.%03d\n", hours, minutes, seconds, milliseconds)

}

func printError(errorMessage string, url string) {
	redColor := "\033[31m"
	resetColor := "\033[0m"
	fmt.Printf("%sError: %s ,url: %s %s\n", redColor, errorMessage, url, resetColor)
}

func testNoDuplicateUrlInMap(url string) {
	//THIS IS FOR SEEING IF THE URL IS ONLY SHOWN ONICE IN THE MAP
	counter := 0
	for tempUrl, _ := range myMap {
		if tempUrl == url {
			counter++
		}
		if counter > 1 {
			printError("Error The url is inside the map more then once ", url)
			log.Fatal("Fix this error")
		}
	}
}
