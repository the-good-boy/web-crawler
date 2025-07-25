package main

import (
	"fmt"
	"time"
)

type Stats struct {
	pagesPerMinute        string // 0 0 \n 1 100
	crawledRatioPerMinute string
	startTime             time.Time
}

func (c *Stats) update(crawled *Visited, queue *Queue, t time.Time) {
	c.pagesPerMinute += fmt.Sprintf("%f %d\n", t.Sub(c.startTime).Minutes(), crawled.size())
	c.crawledRatioPerMinute += fmt.Sprintf("%f %f\n", t.Sub(c.startTime).Minutes(), float64(crawled.size())/float64(queue.size()))
}
func (c *Stats) print() {
	fmt.Println("Pages crawled per minute:")
	fmt.Println(c.pagesPerMinute)
	fmt.Println("Crawl to Queued Ratio per minute:")
	fmt.Println(c.crawledRatioPerMinute)
}
