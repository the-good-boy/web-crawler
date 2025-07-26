package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/net/html"
)

func main() {
	webArchiveAccess := true
	if godotenv.Load() != nil {
		fmt.Println("Error loading .env file. No access to web archive.")
		webArchiveAccess = false
	}
	db := DatabaseConnection{access: webArchiveAccess, uri: "", client: nil, collection: nil}
	db.connect()

	beginUrl := flag.String("url", "https://commoncrawl.org/", "the url you want to start the crawling with!")
	crawlLimit := flag.Int("crawlSize", 20, "Maximum number of webpages you want to crawl.")
	flag.Parse()
	visited := Visited{
		data:     make(map[uint64]struct{}),
		maxLimit: *crawlLimit,
	}

	queue := Queue{
		totalCnt: 0,
		currCnt:  0,
		hrefs:    make([]string, 0),
	}

	ticker := time.NewTicker(1 * time.Minute)

	done := make(chan bool)

	crawlerStats := Stats{pagesPerMinute: "0 0\n", crawledRatioPerMinute: "0 0\n", startTime: time.Now()}

	// tick every minute
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				crawlerStats.update(&visited, &queue, t)
			}
		}
	}()

	c := make(chan *html.Node)
	visited.add(*beginUrl)
	go getHtmlNode(*beginUrl, c)

	htmlNode := <-c
	fmt.Println("Starting crawling with url: ", *beginUrl)
	Crawl(*beginUrl, htmlNode, &queue, &visited, &db)

	for queue.size() > 0 && visited.count < visited.maxLimit {
		url := queue.pop()
		visited.add(url)

		go getHtmlNode(url, c)
		htmlNode := <-c

		// Check if htmlNode is empty before parsing
		if htmlNode == nil || htmlNode.Type == html.ErrorNode {
			continue
		}

		go Crawl(url, htmlNode, &queue, &visited, &db)
	}

	ticker.Stop()
	done <- true
	db.disconnect()
	fmt.Println("\n------------------CRAWLER STATS------------------")
	fmt.Printf("Total queued: %d\n", queue.totalCnt)
	fmt.Printf("To be crawled (Queue) size: %d\n", queue.size())
	fmt.Printf("Crawled size: %d\n", visited.size())
	crawlerStats.print()
}

func getHtmlNode(urlString string, c chan *html.Node) {
	resp, err := http.Get(urlString)
	var doc *html.Node
	if err != nil {
		doc = &html.Node{}
		c <- doc
		return
	}
	defer resp.Body.Close()

	doc, err = html.Parse(resp.Body)
	if err != nil {
		doc = &html.Node{}
		c <- doc
		return
	}

	c <- doc
}
