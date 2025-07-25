package main

import "sync"

type Queue struct {
	totalCnt int
	currCnt  int
	hrefs    []string
	mu       sync.Mutex
}

func (q *Queue) push(url string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.totalCnt++
	q.currCnt++
	q.hrefs = append(q.hrefs, url)
}

func (q *Queue) pop() string {
	q.mu.Lock()
	defer q.mu.Unlock()
	url := q.hrefs[0]
	q.hrefs = q.hrefs[1:]
	q.currCnt--
	return url
}

func (q *Queue) size() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.currCnt
}
