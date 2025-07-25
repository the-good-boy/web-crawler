package main

import (
	"hash/fnv"
	"sync"
)

type Visited struct {
	data     map[uint64]struct{}
	count    int
	mu       sync.Mutex
	maxLimit int
}

func (vis *Visited) add(url string) {
	vis.mu.Lock()
	defer vis.mu.Unlock()
	vis.data[hashUrl(url)] = struct{}{}
	vis.count++
}

func (vis *Visited) contains(url string) bool {
	vis.mu.Lock()
	defer vis.mu.Unlock()
	_, exists := vis.data[hashUrl(url)]
	return exists
}

func (vis *Visited) size() int {
	vis.mu.Lock()
	defer vis.mu.Unlock()
	return vis.count
}
func hashUrl(url string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(url))
	return h.Sum64()
}
