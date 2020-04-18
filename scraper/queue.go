package scraper

import (
	"container/heap"
	"go.uber.org/zap"
	"time"
)

// An Item is something we manage in a scrapeTime queue.
type Item struct {
	url                   string // The url of the item; arbitrary.
	scrapeTime            int    // The scrapeTime of the item in the queue.
	scrapeIntervalSeconds int
	index                 int // The index is needed by update and is maintained by the heap.Interface methods, Index in the heap
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].scrapeTime < pq[j].scrapeTime
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// Sleep until the current time is less than the time on the prioritized item
func (pq *PriorityQueue) PeekAndIncrement() string {
	last := pq[len(pq)-1]
	timeNow := int(time.Now().Unix())

	if timeNow < last.scrapeTime {
		t := last.scrapeTime - timeNow
		zap.S().Debugf("Sleeping for %d seconds", t)
		time.Sleep(time.Duration(t) * time.Second)
	}

	last.scrapeTime += last.scrapeIntervalSeconds
	return last.url
}