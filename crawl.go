package main

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

type WebPage struct {
	Href    string `bson:"href"`
	Title   string `bson:"title"`
	Content string `bson:"content"`
}

func Crawl(url string, node *html.Node, q *Queue, visited *Visited, db *DatabaseConnection) {
	wp := crawlWebPage(node, url, q, visited)

	if visited.count < visited.maxLimit {
		db.insertWebpage(wp)
	}

}

func crawlWebPage(node *html.Node, href string, q *Queue, visited *Visited) WebPage {
	wp := WebPage{
		Href: href,
	}

	var (
		tagCount          int
		titleFound        bool
		pageContentLength int
		inBody            bool
		maxTags           = 500
	)

	var extract func(*html.Node)
	extract = func(node *html.Node) {
		if tagCount >= maxTags {
			return
		}
		tagCount++

		if !titleFound && node.Type == html.ElementNode && node.Data == "title" && node.FirstChild != nil && visited.size() <= visited.maxLimit {
			wp.Title = node.FirstChild.Data
			fmt.Printf("Count: %d | %s -> %s\n", visited.size(), href, wp.Title)
			titleFound = true
		}

		if node.Type == html.ElementNode {

			if node.Data == "body" {
				inBody = true
			}

			if node.Data == "a" {
				for _, att := range node.Attr {
					if att.Key == "href" {
						href := att.Val
						if len(href) > 0 && strings.HasPrefix(href, "http") && !visited.contains(href) {
							q.push(href)
						}
						break
					}
				}
			}
		}

		if inBody && node.Type == html.TextNode && pageContentLength < 500 {
			wp.Content += strings.TrimSpace(node.Data) + " "
			pageContentLength += len(node.Data)
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			extract(child)
			if tagCount >= maxTags {
				return
			}
		}

		if inBody && node.Type == html.ElementNode && node.Data == "body" {
			inBody = false
		}

	}

	extract(node)

	return wp
}
